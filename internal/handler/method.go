package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/model"
	"github.com/wafi04/otomaxv2/internal/services"
	"github.com/wafi04/otomaxv2/pkg/response"
)

type MethodHandler struct {
	methodService *services.MethodService
}

func NewMethodHandler(service *services.MethodService) *MethodHandler {
	return &MethodHandler{
		methodService: service,
	}
}

func (handler *MethodHandler) Create(c *gin.Context) {
	var input model.CreateMethodData
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	data, err := handler.methodService.Create(c.Request.Context(), input)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create method", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Method created successfully", data)
}

func (h *MethodHandler) GetAll(c *gin.Context) {
	// Menggunakan pagination utility
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")
	filterType := c.DefaultQuery("type", "")
	active := c.DefaultQuery("status", "")

	// Hitung pagination
	paginationResult := response.CalculatePagination(&page, &limit)

	// Panggil service dengan response parameters
	data, totalCount, err := h.methodService.GetAll(
		c.Request.Context(),
		paginationResult.Skip,
		paginationResult.Take,
		search,
		filterType,
		active,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}

	// Buat paginated response
	responses := response.CreatePaginatedResponse(
		data,
		paginationResult.CurrentPage,
		paginationResult.ItemsPerPage,
		totalCount,
	)

	response.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", responses)
}

func (h *MethodHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID parameter", err.Error())
		return
	}

	var input model.UpdateMethodData
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	update, err := h.methodService.Update(c.Request.Context(), id, input)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update Method", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Method updated successfully", update)
}

func (h *MethodHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid ID parameter", err.Error())
		return
	}

	err = h.methodService.Delete(c.Request.Context(), id)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete Method", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Method deleted successfully", nil)
}

func (h *MethodHandler) GetByGrub(c *gin.Context) {

	data, err := h.methodService.GetAllGroupedByType(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete Method", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Method deleted successfully", data)
}
