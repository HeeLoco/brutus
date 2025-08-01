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
	tokenAzureService := services.NewTokenBasedAzureService()

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.SecurityHeaders())   // Comprehensive security headers
	router.Use(middleware.CORS())              // Secure CORS configuration
	router.Use(middleware.SetCSRFToken())      // Generate CSRF tokens
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	tokenResourceHandler := handlers.NewTokenResourceHandler(tokenAzureService)
	storageAccountHandler := handlers.NewStorageAccountHandler(tokenAzureService)
	logHandler := handlers.NewLogHandler()

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
		// Logging routes for frontend logs (no CSRF for logging)
		v1.POST("/logs", logHandler.LogEntry)
		v1.POST("/logs/batch", logHandler.LogBatch)

		// Resource group routes (using user tokens) with CSRF protection
		rg := v1.Group("/resource-groups")
		rg.Use(middleware.CSRFProtection()) // Protect state-changing operations
		{
			rg.GET("", tokenResourceHandler.ListResourceGroups)     // No CSRF for GET
			rg.POST("", tokenResourceHandler.CreateResourceGroup)   // CSRF protected
			rg.GET("/:name", tokenResourceHandler.GetResourceGroup) // No CSRF for GET
			rg.PUT("/:name", tokenResourceHandler.UpdateResourceGroup) // CSRF protected
			rg.DELETE("/:name", tokenResourceHandler.DeleteResourceGroup) // CSRF protected
		}

		// Storage account routes (using user tokens) with CSRF protection
		sa := v1.Group("/storage-accounts")
		sa.Use(middleware.CSRFProtection()) // Protect state-changing operations
		{
			sa.GET("", storageAccountHandler.ListStorageAccounts)                                      // No CSRF for GET
			sa.POST("", storageAccountHandler.CreateStorageAccount)                                    // CSRF protected
			sa.GET("/:resourceGroup/:name", storageAccountHandler.GetStorageAccount)                   // No CSRF for GET
			sa.PUT("/:resourceGroup/:name", storageAccountHandler.UpdateStorageAccount)                // CSRF protected
			sa.DELETE("/:resourceGroup/:name", storageAccountHandler.DeleteStorageAccount)             // CSRF protected
		}
	}

	// Start server
	log.Printf("Starting server on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}