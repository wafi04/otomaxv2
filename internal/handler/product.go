package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/services"
	"github.com/wafi04/otomaxv2/pkg/response"
)


type ProductHandler struct {
	productHandler *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productHandler: service,
	}
}


func (h *ProductHandler) GetAll(c *gin.Context) {
	productList, err := h.productHandler.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get news"})
		return
	}
	response.SuccessResponse(c, http.StatusOK, "Products Retreived Successfully", productList)
}