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

	// Initialize token-based Azure service (uses user tokens instead of service credentials)
	tokenAzureService := services.NewTokenBasedAzureService(cfg.SubscriptionID)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	tokenResourceHandler := handlers.NewTokenResourceHandler(tokenAzureService)

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
		// Resource group routes (using user tokens)
		rg := v1.Group("/resource-groups")
		{
			rg.GET("", tokenResourceHandler.ListResourceGroups)
			rg.POST("", tokenResourceHandler.CreateResourceGroup)
			rg.GET("/:name", tokenResourceHandler.GetResourceGroup)
			rg.PUT("/:name", tokenResourceHandler.UpdateResourceGroup)
			rg.DELETE("/:name", tokenResourceHandler.DeleteResourceGroup)
		}
	}

	// Start server
	log.Printf("Starting server on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}