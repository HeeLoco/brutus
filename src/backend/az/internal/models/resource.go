package models

// ResourceGroup represents an Azure resource group
type ResourceGroup struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// CreateResourceGroupRequest represents the request to create a resource group
type CreateResourceGroupRequest struct {
	Name     string            `json:"name" binding:"required"`
	Location string            `json:"location" binding:"required"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// UpdateResourceGroupRequest represents the request to update a resource group
type UpdateResourceGroupRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}