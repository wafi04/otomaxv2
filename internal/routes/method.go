package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/repository"
	"github.com/wafi04/otomaxv2/internal/services"
)

func MethodRoutes(r *gin.RouterGroup, DB *sql.DB) {
	methodRepo := repository.NewMethodRepository(DB)
	methodService := services.NewMethodService(methodRepo)
	methodHandler := handler.NewMethodHandler(methodService)

	categoryGroup := r.Group("/method")
	{
		categoryGroup.POST("", methodHandler.Create)
		categoryGroup.GET("", methodHandler.GetAll)
		categoryGroup.GET("/groub", methodHandler.GetByGrub)

		// categoryGroup.GET("/:id", methodHandler.GetSubCategoryByID)
		categoryGroup.PUT("/:id", methodHandler.Update)
		categoryGroup.DELETE("/:id", methodHandler.Delete)
	}
}
