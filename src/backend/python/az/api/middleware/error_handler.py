from fastapi import Request
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException as StarletteHTTPException

from api.core.exceptions import AzureAPIException
from api.core.logging import get_logger, correlation_id
from api.models.resource_group import ErrorResponse, ErrorDetail

logger = get_logger(__name__)


async def azure_exception_handler(request: Request, exc: AzureAPIException) -> JSONResponse:
    """Handle Azure API specific exceptions."""
    
    logger.warning(f"Azure API exception: {exc.detail}")
    
    error_response = ErrorResponse(
        error=exc.__class__.__name__,
        message=exc.detail,
        correlation_id=exc.correlation_id or correlation_id.get() or "unknown"
    )
    
    return JSONResponse(
        status_code=exc.status_code,
        content=error_response.model_dump()
    )


async def validation_exception_handler(request: Request, exc: RequestValidationError) -> JSONResponse:
    """Handle request validation errors."""
    
    logger.warning(f"Validation error: {exc.errors()}")
    
    # Convert validation errors to our format
    error_details = [
        ErrorDetail(
            field=".".join(str(loc) for loc in error["loc"]),
            message=error["msg"],
            code=error["type"]
        )
        for error in exc.errors()
    ]
    
    error_response = ErrorResponse(
        error="ValidationError",
        message="Request validation failed",
        details=error_details,
        correlation_id=correlation_id.get() or "unknown"
    )
    
    return JSONResponse(
        status_code=422,
        content=error_response.model_dump()
    )


async def http_exception_handler(request: Request, exc: StarletteHTTPException) -> JSONResponse:
    """Handle general HTTP exceptions."""
    
    logger.warning(f"HTTP exception: {exc.detail}")
    
    error_response = ErrorResponse(
        error="HTTPException",
        message=exc.detail,
        correlation_id=correlation_id.get() or "unknown"
    )
    
    return JSONResponse(
        status_code=exc.status_code,
        content=error_response.model_dump()
    )


async def general_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """Handle unexpected exceptions."""
    
    logger.error(f"Unexpected error: {str(exc)}", exc_info=True)
    
    error_response = ErrorResponse(
        error="InternalServerError",
        message="An unexpected error occurred",
        correlation_id=correlation_id.get() or "unknown"
    )
    
    return JSONResponse(
        status_code=500,
        content=error_response.model_dump()
    )