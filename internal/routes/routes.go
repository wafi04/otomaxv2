package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(r *gin.RouterGroup, DB *sql.DB) {
	SubCatgeoryRoutes(r, DB)
	CategoryRoutes(r, DB)
	NewsRoutes(r, DB)
	MethodRoutes(r, DB)

}
