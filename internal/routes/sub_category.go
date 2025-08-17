package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/repository"
	"github.com/wafi04/otomaxv2/internal/services"
)

func SubCatgeoryRoutes(r *gin.RouterGroup, DB *sql.DB) {
	subCategoryRepo := repository.NewSubCategory(DB)
	subCategoryService := services.NewSubCategoryService(subCategoryRepo)
	subCategoryHandler := handler.NewSubCategoryHandler(subCategoryService)

	categoryGroup := r.Group("/subcategories")
	{
		categoryGroup.POST("", subCategoryHandler.CreateSubCategory)
		categoryGroup.GET("", subCategoryHandler.GetAllSubCategories)
		categoryGroup.GET("/:id", subCategoryHandler.GetSubCategoryByID)
		categoryGroup.PUT("/:id", subCategoryHandler.UpdateSubCategory)
		categoryGroup.DELETE("/:id", subCategoryHandler.DeleteSubCategory)
	}
}
