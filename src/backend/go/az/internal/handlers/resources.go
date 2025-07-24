package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/internal/models"
	"brutus/azure-api/internal/services"
)

type ResourceHandler struct {
	azureService *services.AzureService
}

func NewResourceHandler(azureService *services.AzureService) *ResourceHandler {
	return &ResourceHandler{
		azureService: azureService,
	}
}

func (h *ResourceHandler) ListResourceGroups(c *gin.Context) {
	resourceGroups, err := h.azureService.ListResourceGroups(c.Request.Context())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to list resource groups", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resource_groups": resourceGroups,
		"count":           len(resourceGroups),
	})
}

func (h *ResourceHandler) GetResourceGroup(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	resourceGroup, err := h.azureService.GetResourceGroup(c.Request.Context(), name)
	if err != nil {
		h.handleError(c, http.StatusNotFound, "Resource group not found", err)
		return
	}

	c.JSON(http.StatusOK, resourceGroup)
}

func (h *ResourceHandler) CreateResourceGroup(c *gin.Context) {
	var req models.CreateResourceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resourceGroup, err := h.azureService.CreateResourceGroup(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to create resource group", err)
		return
	}

	c.JSON(http.StatusCreated, resourceGroup)
}

func (h *ResourceHandler) UpdateResourceGroup(c *gin.Context) {
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

	resourceGroup, err := h.azureService.UpdateResourceGroup(c.Request.Context(), name, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to update resource group", err)
		return
	}

	c.JSON(http.StatusOK, resourceGroup)
}

func (h *ResourceHandler) DeleteResourceGroup(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		h.handleError(c, http.StatusBadRequest, "Resource group name is required", nil)
		return
	}

	err := h.azureService.DeleteResourceGroup(c.Request.Context(), name)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to delete resource group", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resource group deletion initiated",
		"name":    name,
	})
}

func (h *ResourceHandler) handleError(c *gin.Context, statusCode int, message string, err error) {
	response := models.ErrorResponse{
		Error:   message,
		Code:    statusCode,
	}

	if err != nil {
		response.Message = err.Error()
	}

	c.JSON(statusCode, response)
}