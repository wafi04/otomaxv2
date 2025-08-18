package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/config"
	"github.com/wafi04/otomaxv2/internal/integrations/duitku"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/services"
	"github.com/wafi04/otomaxv2/pkg/response"
	"github.com/wafi04/otomaxv2/pkg/utils"
)

type DepositHandler struct {
	cfg         config.Config
	depoService *services.DepositService
}

func NewDepositHandler(depoService *services.DepositService) *DepositHandler {
	return &DepositHandler{
		depoService: depoService,
	}
}

func (h *DepositHandler) Create(c *gin.Context) {
	var input model.RequestFormClient
	duitkuMethod := duitku.NewDuitkuService(&h.cfg)

	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}
	depStr := "DEP"
	callbackUrl := "http://localhost:8080/api/callback/duitku"

	duitkuCall, err := duitkuMethod.CreateTransaction(c, &duitku.DuitkuCreateTransactionParams{
		PaymentAmount:   input.Amount,
		MerchantOrderId: utils.GenerateUniqeID(&depStr),
		ProductDetails:  "Deposit",
		PaymentCode:     input.Method,
		CallbackUrl:     &callbackUrl,
		ReturnUrl:       &callbackUrl,
	})
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	_, err = h.depoService.Create(c.Request.Context(), model.CreateDeposit{
		Amount:            input.Amount,
		Method:            input.Method,
		Username:          "",
		DestinationNumber: "",
		PaymentReferee:    &duitkuCall.Reference,
	})

	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Category created successfully", nil)

}

func (h *DepositHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	paginationResult := response.CalculatePagination(&page, &limit)

	data, totalCount, err := h.depoService.GetAll(
		c.Request.Context(),
		model.FilterDeposit{
			Search: &search,
			Status: &status,
			Limit:  paginationResult.Take,
			Offset: paginationResult.Skip,
		},
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch deposits", err.Error())
		return
	}

	responses := response.CreatePaginatedResponse(
		data,
		paginationResult.CurrentPage,
		paginationResult.ItemsPerPage,
		totalCount,
	)

	response.SuccessResponse(c, http.StatusOK, "Deposits retrieved successfully", responses)
}
