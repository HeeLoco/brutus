from typing import List

from azure.core.exceptions import ResourceNotFoundError as AzureResourceNotFoundError
from fastapi import APIRouter, Depends, HTTPException, status

from api.core.config import Settings, get_settings
from api.core.exceptions import ResourceNotFoundError, AzureConnectionError
from api.core.logging import get_logger, correlation_id
from api.models.resource_group import (
    ResourceGroup,
    ResourceGroupCreate,
    ResourceGroupUpdate,
    ResourceGroupList
)
from api.services.azure_service import AzureService

router = APIRouter(prefix="/resource-groups", tags=["resource groups"])
logger = get_logger(__name__)


def get_azure_service(settings: Settings = Depends(get_settings)) -> AzureService:
    """Dependency to get Azure service instance."""
    return AzureService(settings)


@router.get("/", response_model=ResourceGroupList, summary="List all resource groups")
async def list_resource_groups(
    azure_service: AzureService = Depends(get_azure_service)
) -> ResourceGroupList:
    """
    List all resource groups in the subscription.
    
    Returns:
        ResourceGroupList: List of resource groups with count
    """
    try:
        logger.info("Listing resource groups")
        resource_groups = await azure_service.list_resource_groups()
        
        return ResourceGroupList(
            resource_groups=resource_groups,
            count=len(resource_groups)
        )
        
    except Exception as e:
        logger.error(f"Failed to list resource groups: {e}")
        raise AzureConnectionError(
            detail="Failed to retrieve resource groups from Azure",
            correlation_id=correlation_id.get()
        )


@router.get("/{name}", response_model=ResourceGroup, summary="Get a specific resource group")
async def get_resource_group(
    name: str,
    azure_service: AzureService = Depends(get_azure_service)
) -> ResourceGroup:
    """
    Get a specific resource group by name.
    
    Args:
        name: The resource group name
        
    Returns:
        ResourceGroup: The resource group details
    """
    try:
        logger.info(f"Getting resource group: {name}")
        return await azure_service.get_resource_group(name)
        
    except AzureResourceNotFoundError:
        raise ResourceNotFoundError(
            resource_name=name,
            resource_type="resource group",
            correlation_id=correlation_id.get()
        )
    except Exception as e:
        logger.error(f"Failed to get resource group {name}: {e}")
        raise AzureConnectionError(
            detail=f"Failed to retrieve resource group '{name}' from Azure",
            correlation_id=correlation_id.get()
        )


@router.post("/", response_model=ResourceGroup, status_code=status.HTTP_201_CREATED, 
            summary="Create a new resource group")
async def create_resource_group(
    resource_group: ResourceGroupCreate,
    azure_service: AzureService = Depends(get_azure_service)
) -> ResourceGroup:
    """
    Create a new resource group.
    
    Args:
        resource_group: Resource group creation data
        
    Returns:
        ResourceGroup: The created resource group
    """
    try:
        logger.info(f"Creating resource group: {resource_group.name}")
        return await azure_service.create_resource_group(resource_group)
        
    except Exception as e:
        logger.error(f"Failed to create resource group {resource_group.name}: {e}")
        
        # Check if it's a conflict error (already exists)
        if "already exists" in str(e).lower() or "conflict" in str(e).lower():
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail=f"Resource group '{resource_group.name}' already exists"
            )
        
        raise AzureConnectionError(
            detail=f"Failed to create resource group '{resource_group.name}'",
            correlation_id=correlation_id.get()
        )


@router.put("/{name}", response_model=ResourceGroup, summary="Update a resource group")
async def update_resource_group(
    name: str,
    update_data: ResourceGroupUpdate,
    azure_service: AzureService = Depends(get_azure_service)
) -> ResourceGroup:
    """
    Update an existing resource group.
    
    Args:
        name: The resource group name
        update_data: Resource group update data
        
    Returns:
        ResourceGroup: The updated resource group
    """
    try:
        logger.info(f"Updating resource group: {name}")
        return await azure_service.update_resource_group(name, update_data)
        
    except AzureResourceNotFoundError:
        raise ResourceNotFoundError(
            resource_name=name,
            resource_type="resource group",
            correlation_id=correlation_id.get()
        )
    except Exception as e:
        logger.error(f"Failed to update resource group {name}: {e}")
        raise AzureConnectionError(
            detail=f"Failed to update resource group '{name}'",
            correlation_id=correlation_id.get()
        )


@router.delete("/{name}", status_code=status.HTTP_202_ACCEPTED, 
              summary="Delete a resource group")
async def delete_resource_group(
    name: str,
    azure_service: AzureService = Depends(get_azure_service)
):
    """
    Delete a resource group.
    
    Note: This operation is asynchronous in Azure and may take several minutes.
    
    Args:
        name: The resource group name
        
    Returns:
        dict: Deletion status message
    """
    try:
        logger.info(f"Deleting resource group: {name}")
        await azure_service.delete_resource_group(name)
        
        return {
            "message": f"Resource group '{name}' deletion initiated",
            "note": "Deletion is asynchronous and may take several minutes to complete"
        }
        
    except AzureResourceNotFoundError:
        raise ResourceNotFoundError(
            resource_name=name,
            resource_type="resource group",
            correlation_id=correlation_id.get()
        )
    except Exception as e:
        logger.error(f"Failed to delete resource group {name}: {e}")
        raise AzureConnectionError(
            detail=f"Failed to delete resource group '{name}'",
            correlation_id=correlation_id.get()
        )