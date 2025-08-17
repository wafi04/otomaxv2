package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/repository"
	"github.com/wafi04/otomaxv2/internal/services"
)

func CategoryRoutes(r *gin.RouterGroup, DB *sql.DB) {
	categoryRepo := repository.NewCategoryRepository(DB)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	categoryGroup := r.Group("/categories")
	{
		categoryGroup.POST("", categoryHandler.CreateCategory)
		categoryGroup.GET("", categoryHandler.GetAllCategories)
		categoryGroup.GET("/:code", categoryHandler.GetCategoryByCode)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

}
