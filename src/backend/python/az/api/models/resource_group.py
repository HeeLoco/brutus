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