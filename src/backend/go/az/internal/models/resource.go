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

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// StorageAccount represents an Azure storage account
type StorageAccount struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Location          string            `json:"location"`
	ResourceGroup     string            `json:"resourceGroup"`
	Kind              string            `json:"kind"`
	SkuName           string            `json:"skuName"`
	SkuTier           string            `json:"skuTier"`
	AccessTier        string            `json:"accessTier,omitempty"`
	AllowBlobPublic   bool              `json:"allowBlobPublicAccess"`
	AllowSharedKey    bool              `json:"allowSharedKeyAccess"`
	Tags              map[string]string `json:"tags,omitempty"`
	CreationTime      string            `json:"creationTime,omitempty"`
	PrimaryEndpoints  *StorageEndpoints `json:"primaryEndpoints,omitempty"`
}

// StorageEndpoints represents storage account endpoints
type StorageEndpoints struct {
	Blob  string `json:"blob,omitempty"`
	Queue string `json:"queue,omitempty"`
	Table string `json:"table,omitempty"`
	File  string `json:"file,omitempty"`
}

// CreateStorageAccountRequest represents the request to create a storage account
type CreateStorageAccountRequest struct {
	Name              string            `json:"name" binding:"required"`
	ResourceGroup     string            `json:"resourceGroup" binding:"required"`
	Location          string            `json:"location" binding:"required"`
	Kind              string            `json:"kind,omitempty"` // StorageV2, BlobStorage, etc.
	SkuName           string            `json:"skuName,omitempty"` // Standard_LRS, Premium_LRS, etc.
	AccessTier        string            `json:"accessTier,omitempty"` // Hot, Cool, Archive
	AllowBlobPublic   *bool             `json:"allowBlobPublicAccess,omitempty"`
	AllowSharedKey    *bool             `json:"allowSharedKeyAccess,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// UpdateStorageAccountRequest represents the request to update a storage account
type UpdateStorageAccountRequest struct {
	AccessTier      string            `json:"accessTier,omitempty"`
	AllowBlobPublic *bool             `json:"allowBlobPublicAccess,omitempty"`
	AllowSharedKey  *bool             `json:"allowSharedKeyAccess,omitempty"`
	Tags            map[string]string `json:"tags,omitempty"`
}