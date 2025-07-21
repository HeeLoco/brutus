package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/internal/models"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}