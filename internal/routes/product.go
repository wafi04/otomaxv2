package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/repository"
	"github.com/wafi04/otomaxv2/internal/services"
)

func ProductRoutes(r *gin.RouterGroup, DB *sql.DB) {
	productRepo := repository.NewProductRepository(DB)
	productService := services.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	productGroup := r.Group("/products")
	{
		productGroup.GET("", productHandler.GetAll)
	}
}
