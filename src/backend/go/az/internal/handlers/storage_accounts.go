package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/internal/models"
	"brutus/azure-api/internal/services"
)

type StorageAccountHandler struct {
	tokenAzureService *services.TokenBasedAzureService
}

func NewStorageAccountHandler(tokenAzureService *services.TokenBasedAzureService) *StorageAccountHandler {
	return &StorageAccountHandler{
		tokenAzureService: tokenAzureService,
	}
}

// extractBearerToken extracts the bearer token from the Authorization header
func (h *StorageAccountHandler) extractBearerToken(c *gin.Context) (string, error) {
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
func (h *StorageAccountHandler) extractSubscriptionID(c *gin.Context) (string, error) {
	subscriptionID := c.GetHeader("X-Azure-Subscription-ID")
	if subscriptionID == "" {
		return "", &models.AuthError{Message: "X-Azure-Subscription-ID header is required"}
	}
	return subscriptionID, nil
}

func (h *StorageAccountHandler) ListStorageAccounts(c *gin.Context) {
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

	storageAccounts, err := h.tokenAzureService.ListStorageAccounts(c.Request.Context(), token, subscriptionID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to list storage accounts", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"storage_accounts": storageAccounts,
		"count":           len(storageAccounts),
	})
}

func (h *StorageAccountHandler) GetStorageAccount(c *gin.Context) {
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

	resourceGroup := c.Param("resourceGroup")
	if resourceGroup == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Storage account name is required", nil)
		return
	}

	storageAccount, err := h.tokenAzureService.GetStorageAccount(c.Request.Context(), token, subscriptionID, resourceGroup, name)
	if err != nil {
		h.handleError(c, http.StatusNotFound, "Storage account not found", err)
		return
	}

	c.JSON(http.StatusOK, storageAccount)
}

func (h *StorageAccountHandler) CreateStorageAccount(c *gin.Context) {
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

	var req models.CreateStorageAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	storageAccount, err := h.tokenAzureService.CreateStorageAccount(c.Request.Context(), token, subscriptionID, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to create storage account", err)
		return
	}

	c.JSON(http.StatusCreated, storageAccount)
}

func (h *StorageAccountHandler) UpdateStorageAccount(c *gin.Context) {
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

	resourceGroup := c.Param("resourceGroup")
	if resourceGroup == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Storage account name is required", nil)
		return
	}

	var req models.UpdateStorageAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	storageAccount, err := h.tokenAzureService.UpdateStorageAccount(c.Request.Context(), token, subscriptionID, resourceGroup, name, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to update storage account", err)
		return
	}

	c.JSON(http.StatusOK, storageAccount)
}

func (h *StorageAccountHandler) DeleteStorageAccount(c *gin.Context) {
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

	resourceGroup := c.Param("resourceGroup")
	if resourceGroup == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Storage account name is required", nil)
		return
	}

	err = h.tokenAzureService.DeleteStorageAccount(c.Request.Context(), token, subscriptionID, resourceGroup, name)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to delete storage account", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Storage account deleted successfully",
		"name":          name,
		"resourceGroup": resourceGroup,
	})
}

func (h *StorageAccountHandler) handleAuthError(c *gin.Context, err error) {
	response := models.ErrorResponse{
		Error:   "Authentication failed",
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	}
	c.JSON(http.StatusUnauthorized, response)
}

func (h *StorageAccountHandler) handleError(c *gin.Context, statusCode int, message string, err error) {
	response := models.ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	if err != nil {
		response.Message = err.Error()
	}

	c.JSON(statusCode, response)
}