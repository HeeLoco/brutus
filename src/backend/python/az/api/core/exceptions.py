from typing import List, Optional, Any

from fastapi import HTTPException, status


class AzureAPIException(HTTPException):
    """Base exception for Azure API errors."""
    
    def __init__(
        self,
        status_code: int,
        detail: str,
        headers: Optional[dict] = None,
        correlation_id: Optional[str] = None
    ):
        super().__init__(status_code=status_code, detail=detail, headers=headers)
        self.correlation_id = correlation_id


class ResourceNotFoundError(AzureAPIException):
    """Exception raised when a resource is not found."""
    
    def __init__(self, resource_name: str, resource_type: str = "resource", correlation_id: Optional[str] = None):
        detail = f"{resource_type.title()} '{resource_name}' not found"
        super().__init__(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=detail,
            correlation_id=correlation_id
        )


class ResourceAlreadyExistsError(AzureAPIException):
    """Exception raised when trying to create a resource that already exists."""
    
    def __init__(self, resource_name: str, resource_type: str = "resource", correlation_id: Optional[str] = None):
        detail = f"{resource_type.title()} '{resource_name}' already exists"
        super().__init__(
            status_code=status.HTTP_409_CONFLICT,
            detail=detail,
            correlation_id=correlation_id
        )


class AzureConnectionError(AzureAPIException):
    """Exception raised when Azure connection fails."""
    
    def __init__(self, detail: str = "Unable to connect to Azure services", correlation_id: Optional[str] = None):
        super().__init__(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail=detail,
            correlation_id=correlation_id
        )


class ValidationError(AzureAPIException):
    """Exception raised for validation errors."""
    
    def __init__(
        self,
        detail: str,
        errors: Optional[List[dict]] = None,
        correlation_id: Optional[str] = None
    ):
        super().__init__(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            detail=detail,
            correlation_id=correlation_id
        )
        self.errors = errors or []


class AzurePermissionError(AzureAPIException):
    """Exception raised when Azure operation is not permitted."""
    
    def __init__(self, detail: str = "Operation not permitted", correlation_id: Optional[str] = None):
        super().__init__(
            status_code=status.HTTP_403_FORBIDDEN,
            detail=detail,
            correlation_id=correlation_id
        )