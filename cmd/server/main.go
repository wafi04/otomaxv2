package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/config"
	"github.com/wafi04/otomaxv2/internal/routes"
	"github.com/wafi04/otomaxv2/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Logger.Warn("Failed To Setup Env")
	}
	log.Logger.Info("Setup Env Successfully")
	db, err := config.NewDatabaseConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://your-frontend.com"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group("/api")

	routes.ProductExternalRoutes(api, *cfg, db.SqlDB)

	routes.SetupAllRoutes(api, db.SqlDB)
	r.Run(cfg.Server.Host + ":" + cfg.Server.Port)

}
