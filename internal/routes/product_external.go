package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/config"
	"github.com/wafi04/otomaxv2/internal/handler"
	"github.com/wafi04/otomaxv2/internal/integrations/digiflazz"
	"github.com/wafi04/otomaxv2/internal/services/productexternal"
)

func ProductExternalRoutes(r *gin.RouterGroup, cfg config.Config, db *sql.DB) {
	digiService := digiflazz.NewDigiflazzService(digiflazz.DigiConfig{
		DigiKey:      cfg.Digiflazz.DigiKey,
		DigiUsername: cfg.Digiflazz.DigiUsername,
	})

	productExternalService := productexternal.NewProductExternal(digiService, db)
	productExternalHandler := handler.NewProductExternalHandler(productExternalService)
	syncProduct := r.Group("/sync/product")

	{
		syncProduct.GET("/digiflazz", productExternalHandler.GetByDigiflazz)
	}

}
