from datetime import datetime
from typing import Dict, Optional

from pydantic import BaseModel, Field


class ResourceGroupBase(BaseModel):
    """Base resource group model."""
    name: str = Field(..., min_length=1, max_length=90, description="Resource group name")
    location: str = Field(..., description="Azure location/region")
    tags: Optional[Dict[str, str]] = Field(default=None, description="Resource tags")


class ResourceGroupCreate(ResourceGroupBase):
    """Model for creating a resource group."""
    pass


class ResourceGroupUpdate(BaseModel):
    """Model for updating a resource group."""
    tags: Optional[Dict[str, str]] = Field(default=None, description="Resource tags to update")


class ResourceGroup(ResourceGroupBase):
    """Complete resource group model with all fields."""
    id: str = Field(..., description="Azure resource ID")
    type: str = Field(default="Microsoft.Resources/resourceGroups", description="Resource type")
    managed_by: Optional[str] = Field(default=None, description="Resource manager")
    
    class Config:
        from_attributes = True


class ResourceGroupList(BaseModel):
    """Model for listing resource groups."""
    resource_groups: list[ResourceGroup] = Field(..., description="List of resource groups")
    count: int = Field(..., description="Total count")


class HealthCheck(BaseModel):
    """Health check response model."""
    status: str = Field(default="healthy", description="Service status")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Check timestamp")
    version: str = Field(..., description="API version")
    azure_connection: bool = Field(..., description="Azure connectivity status")


class ErrorDetail(BaseModel):
    """Error detail model."""
    field: Optional[str] = Field(default=None, description="Field that caused the error")
    message: str = Field(..., description="Error message")
    code: Optional[str] = Field(default=None, description="Error code")


class ErrorResponse(BaseModel):
    """Standard error response model."""
    error: str = Field(..., description="Error type")
    message: str = Field(..., description="Human-readable error message")
    details: Optional[list[ErrorDetail]] = Field(default=None, description="Error details")
    correlation_id: str = Field(..., description="Request correlation ID")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Error timestamp")


# Storage Account Models

class StorageEndpoints(BaseModel):
    """Storage account endpoints model."""
    blob: Optional[str] = Field(default=None, description="Blob endpoint")
    queue: Optional[str] = Field(default=None, description="Queue endpoint")
    table: Optional[str] = Field(default=None, description="Table endpoint")
    file: Optional[str] = Field(default=None, description="File endpoint")


class StorageAccountBase(BaseModel):
    """Base storage account model."""
    name: str = Field(..., min_length=3, max_length=24, description="Storage account name")
    resource_group: str = Field(..., description="Resource group name")
    location: str = Field(..., description="Azure location/region")
    kind: Optional[str] = Field(default="StorageV2", description="Storage account kind")
    sku_name: Optional[str] = Field(default="Standard_LRS", description="SKU name")
    access_tier: Optional[str] = Field(default=None, description="Access tier")
    allow_blob_public: Optional[bool] = Field(default=None, description="Allow blob public access")
    allow_shared_key: Optional[bool] = Field(default=None, description="Allow shared key access")
    tags: Optional[Dict[str, str]] = Field(default=None, description="Storage account tags")


class StorageAccountCreate(StorageAccountBase):
    """Model for creating a storage account."""
    pass


class StorageAccountUpdate(BaseModel):
    """Model for updating a storage account."""
    access_tier: Optional[str] = Field(default=None, description="Access tier")
    allow_blob_public: Optional[bool] = Field(default=None, description="Allow blob public access")
    allow_shared_key: Optional[bool] = Field(default=None, description="Allow shared key access")
    tags: Optional[Dict[str, str]] = Field(default=None, description="Storage account tags")


class StorageAccount(StorageAccountBase):
    """Complete storage account model with all fields."""
    id: str = Field(..., description="Azure resource ID")
    sku_tier: Optional[str] = Field(default=None, description="SKU tier")
    creation_time: Optional[datetime] = Field(default=None, description="Creation timestamp")
    primary_endpoints: Optional[StorageEndpoints] = Field(default=None, description="Primary endpoints")
    
    class Config:
        from_attributes = True


class StorageAccountList(BaseModel):
    """Model for listing storage accounts."""
    storage_accounts: list[StorageAccount] = Field(..., description="List of storage accounts")
    count: int = Field(..., description="Total count")