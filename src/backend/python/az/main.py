import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException

from api.core.config import get_settings
from api.core.logging import setup_logging, get_logger
from api.core.exceptions import AzureAPIException
from api.middleware.correlation import CorrelationMiddleware
from api.middleware.error_handler import (
    azure_exception_handler,
    validation_exception_handler,
    http_exception_handler,
    general_exception_handler
)
from api.routers import health, resource_groups, storage_accounts

# Initialize logging first
setup_logging()
logger = get_logger(__name__)

# Get settings
settings = get_settings()

# Create FastAPI app
app = FastAPI(
    title=settings.app_name,
    version=settings.version,
    description="Production-ready Azure Resource Management API",
    docs_url="/docs",
    redoc_url="/redoc"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE"],
    allow_headers=["*"],
)

# Add custom middleware
app.add_middleware(CorrelationMiddleware)

# Add exception handlers
app.add_exception_handler(AzureAPIException, azure_exception_handler)
app.add_exception_handler(RequestValidationError, validation_exception_handler)
app.add_exception_handler(HTTPException, http_exception_handler)
app.add_exception_handler(Exception, general_exception_handler)

# Include routers
app.include_router(health.router)
app.include_router(resource_groups.router, prefix=settings.api_v1_prefix)
app.include_router(storage_accounts.router, prefix=settings.api_v1_prefix)


@app.get("/", tags=["root"])
async def root():
    """Root endpoint with API information."""
    return {
        "name": settings.app_name,
        "version": settings.version,
        "status": "running",
        "docs": "/docs",
        "health": "/health"
    }


@app.on_event("startup")
async def startup_event():
    """Application startup tasks."""
    logger.info(f"Starting {settings.app_name} v{settings.version}")
    logger.info(f"Debug mode: {settings.debug}")
    logger.info(f"Docs available at: /docs")


@app.on_event("shutdown")
async def shutdown_event():
    """Application shutdown tasks."""
    logger.info(f"Shutting down {settings.app_name}")


if __name__ == "__main__":
    logger.info(f"Starting server on {settings.host}:{settings.port}")
    
    uvicorn.run(
        "main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.debug,
        log_level=settings.log_level.lower()
    )