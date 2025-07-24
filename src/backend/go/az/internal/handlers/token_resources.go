package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/internal/models"
	"brutus/azure-api/internal/services"
)

type TokenResourceHandler struct {
	tokenAzureService *services.TokenBasedAzureService
}

func NewTokenResourceHandler(tokenAzureService *services.TokenBasedAzureService) *TokenResourceHandler {
	return &TokenResourceHandler{
		tokenAzureService: tokenAzureService,
	}
}

// extractBearerToken extracts the bearer token from the Authorization header
func (h *TokenResourceHandler) extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", &models.AuthError{Message: "Authorization header is required"}
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", &models.AuthError{Message: "Authorization header must be Bearer token"}
	}

	// Extract token (remove "Bearer " prefix)
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", &models.AuthError{Message: "Bearer token is empty"}
	}

	return token, nil
}

// extractSubscriptionID extracts the subscription ID from the X-Azure-Subscription-ID header
func (h *TokenResourceHandler) extractSubscriptionID(c *gin.Context) (string, error) {
	subscriptionID := c.GetHeader("X-Azure-Subscription-ID")
	if subscriptionID == "" {
		return "", &models.AuthError{Message: "X-Azure-Subscription-ID header is required"}
	}
	return subscriptionID, nil
}

func (h *TokenResourceHandler) ListResourceGroups(c *gin.Context) {
	token, err := h.extractBearerToken(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	subscriptionID, err := h.extractSubscriptionID(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	resourceGroups, err := h.tokenAzureService.ListResourceGroups(c.Request.Context(), token, subscriptionID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to list resource groups", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resource_groups": resourceGroups,
		"count":           len(resourceGroups),
	})
}

func (h *TokenResourceHandler) GetResourceGroup(c *gin.Context) {
	token, err := h.extractBearerToken(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	subscriptionID, err := h.extractSubscriptionID(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	resourceGroup, err := h.tokenAzureService.GetResourceGroup(c.Request.Context(), token, subscriptionID, name)
	if err != nil {
		h.handleError(c, http.StatusNotFound, "Resource group not found", err)
		return
	}

	c.JSON(http.StatusOK, resourceGroup)
}

func (h *TokenResourceHandler) CreateResourceGroup(c *gin.Context) {
	token, err := h.extractBearerToken(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	subscriptionID, err := h.extractSubscriptionID(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	var req models.CreateResourceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resourceGroup, err := h.tokenAzureService.CreateResourceGroup(c.Request.Context(), token, subscriptionID, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to create resource group", err)
		return
	}

	c.JSON(http.StatusCreated, resourceGroup)
}

func (h *TokenResourceHandler) UpdateResourceGroup(c *gin.Context) {
	token, err := h.extractBearerToken(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	subscriptionID, err := h.extractSubscriptionID(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	var req models.UpdateResourceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resourceGroup, err := h.tokenAzureService.UpdateResourceGroup(c.Request.Context(), token, subscriptionID, name, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to update resource group", err)
		return
	}

	c.JSON(http.StatusOK, resourceGroup)
}

func (h *TokenResourceHandler) DeleteResourceGroup(c *gin.Context) {
	token, err := h.extractBearerToken(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	subscriptionID, err := h.extractSubscriptionID(c)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	err = h.tokenAzureService.DeleteResourceGroup(c.Request.Context(), token, subscriptionID, name)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to delete resource group", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resource group deletion initiated",
		"name":    name,
	})
}

func (h *TokenResourceHandler) handleAuthError(c *gin.Context, err error) {
	response := models.ErrorResponse{
		Error:   "Authentication failed",
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	}
	c.JSON(http.StatusUnauthorized, response)
}

func (h *TokenResourceHandler) handleError(c *gin.Context, statusCode int, message string, err error) {
	response := models.ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	if err != nil {
		response.Message = err.Error()
	}

	c.JSON(statusCode, response)
}