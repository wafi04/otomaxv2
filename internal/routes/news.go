package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/repository"
	"github.com/wafi04/otomaxv2/internal/services"
)

func NewsRoutes(r *gin.RouterGroup, DB *sql.DB) {
	newsRepo := repository.NewNewsRepository(DB)
	newsService := services.NewNewsService(newsRepo)
	newsHandler := handler.NewNewsHandler(newsService)

	categoryGroup := r.Group("/news")
	{
		categoryGroup.POST("", newsHandler.Create)
		categoryGroup.GET("", newsHandler.GetAll)
		// categoryGroup.GET("/:id", newsHandler.GetSubCategoryByID)
		categoryGroup.PUT("/:id", newsHandler.Update)
		categoryGroup.DELETE("/:id", newsHandler.Delete)
	}
}
