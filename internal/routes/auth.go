package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/handler"
)

func AuthRoutes(r *gin.RouterGroup, DB *sql.DB) {


	categoryGroup := r.Group("/auth")
	{
		categoryGroup.GET("", handler.GoogleLoginHandler)
		categoryGroup.GET("/google/callback", handler.GoogleCallbackHandler)
	}

}
