from fastapi import APIRouter, Depends

from api.core.config import Settings, get_settings
from api.core.logging import get_logger
from api.models.resource_group import HealthCheck
from api.services.azure_service import AzureService

router = APIRouter(tags=["health"])
logger = get_logger(__name__)


@router.get("/health", response_model=HealthCheck, summary="Health check endpoint")
async def health_check(settings: Settings = Depends(get_settings)) -> HealthCheck:
    """
    Check the health of the API and its dependencies.
    
    Returns:
        HealthCheck: Service health status including Azure connectivity
    """
    logger.info("Performing health check")
    
    try:
        # Test Azure connection
        azure_service = AzureService(settings)
        azure_connection = await azure_service.test_connection()
        
        logger.info(f"Health check completed - Azure connection: {azure_connection}")
        
        return HealthCheck(
            status="healthy" if azure_connection else "degraded",
            version=settings.version,
            azure_connection=azure_connection
        )
        
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        return HealthCheck(
            status="unhealthy",
            version=settings.version,
            azure_connection=False
        )