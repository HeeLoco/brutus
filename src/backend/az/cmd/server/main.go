package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/configs"
	"brutus/azure-api/internal/handlers"
	"brutus/azure-api/internal/middleware"
	"brutus/azure-api/internal/services"
)

func main() {
	// Load configuration
	cfg := configs.Load()

	// Initialize Azure service
	azureService, err := services.NewAzureService(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Azure service:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	resourceHandler := handlers.NewResourceHandler(azureService)

	// Health check routes
	router.GET("/health", healthHandler.Check)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Azure API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Resource group routes
		rg := v1.Group("/resource-groups")
		{
			rg.GET("", resourceHandler.ListResourceGroups)
			rg.POST("", resourceHandler.CreateResourceGroup)
			rg.GET("/:name", resourceHandler.GetResourceGroup)
			rg.PUT("/:name", resourceHandler.UpdateResourceGroup)
			rg.DELETE("/:name", resourceHandler.DeleteResourceGroup)
		}
	}

	// Start server
	log.Printf("Starting server on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}