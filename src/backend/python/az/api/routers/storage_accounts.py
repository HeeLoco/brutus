from typing import List

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import JSONResponse

from api.core.logging import get_logger
from api.models.resource_group import (
    StorageAccount, StorageAccountCreate, StorageAccountUpdate, StorageAccountList,
    ErrorResponse
)
from api.services.azure_service import AzureService
from api.core.config import get_settings

logger = get_logger(__name__)
router = APIRouter(prefix="/storage-accounts", tags=["storage-accounts"])


def get_azure_service() -> AzureService:
    """Dependency to get Azure service instance."""
    return AzureService(get_settings())


@router.get(
    "",
    response_model=StorageAccountList,
    summary="List storage accounts",
    description="Retrieve all storage accounts in the subscription"
)
async def list_storage_accounts(
    azure_service: AzureService = Depends(get_azure_service)
) -> StorageAccountList:
    """List all storage accounts in the subscription."""
    try:
        logger.info("API: Listing storage accounts")
        
        storage_accounts = await azure_service.list_storage_accounts()
        
        return StorageAccountList(
            storage_accounts=storage_accounts,
            count=len(storage_accounts)
        )
        
    except Exception as e:
        logger.error(f"API: Failed to list storage accounts: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to list storage accounts: {str(e)}"
        )


@router.get(
    "/{resource_group}/{name}",
    response_model=StorageAccount,
    summary="Get storage account",
    description="Retrieve a specific storage account by resource group and name",
    responses={
        404: {"model": ErrorResponse, "description": "Storage account not found"}
    }
)
async def get_storage_account(
    resource_group: str,
    name: str,
    azure_service: AzureService = Depends(get_azure_service)
) -> StorageAccount:
    """Get a specific storage account by resource group and name."""
    try:
        logger.info(f"API: Getting storage account {name} in resource group {resource_group}")
        
        return await azure_service.get_storage_account(resource_group, name)
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"API: Failed to get storage account {name}: {e}")
        if "not found" in str(e).lower():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Storage account '{name}' not found in resource group '{resource_group}'"
            )
        
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to get storage account: {str(e)}"
        )


@router.post(
    "",
    response_model=StorageAccount,
    status_code=status.HTTP_201_CREATED,
    summary="Create storage account",
    description="Create a new storage account",
    responses={
        400: {"model": ErrorResponse, "description": "Invalid request"},
        409: {"model": ErrorResponse, "description": "Storage account already exists"}
    }
)
async def create_storage_account(
    storage_account: StorageAccountCreate,
    azure_service: AzureService = Depends(get_azure_service)
) -> StorageAccount:
    """Create a new storage account."""
    try:
        logger.info(f"API: Creating storage account {storage_account.name}")
        
        result = await azure_service.create_storage_account(storage_account)
        
        logger.info(f"API: Successfully created storage account {storage_account.name}")
        return result
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"API: Failed to create storage account {storage_account.name}: {e}")
        
        error_msg = str(e).lower()
        if "already exists" in error_msg or "conflict" in error_msg:
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail=f"Storage account '{storage_account.name}' already exists"
            )
        elif "invalid" in error_msg or "bad request" in error_msg:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Invalid storage account configuration: {str(e)}"
            )
        
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to create storage account: {str(e)}"
        )


@router.put(
    "/{resource_group}/{name}",
    response_model=StorageAccount,
    summary="Update storage account",
    description="Update an existing storage account",
    responses={
        404: {"model": ErrorResponse, "description": "Storage account not found"}
    }
)
async def update_storage_account(
    resource_group: str,
    name: str,
    update: StorageAccountUpdate,
    azure_service: AzureService = Depends(get_azure_service)
) -> StorageAccount:
    """Update an existing storage account."""
    try:
        logger.info(f"API: Updating storage account {name} in resource group {resource_group}")
        
        result = await azure_service.update_storage_account(resource_group, name, update)
        
        logger.info(f"API: Successfully updated storage account {name}")
        return result
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"API: Failed to update storage account {name}: {e}")
        
        if "not found" in str(e).lower():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Storage account '{name}' not found in resource group '{resource_group}'"
            )
        
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to update storage account: {str(e)}"
        )


@router.delete(
    "/{resource_group}/{name}",
    status_code=status.HTTP_200_OK,
    summary="Delete storage account",
    description="Delete a storage account",
    responses={
        404: {"model": ErrorResponse, "description": "Storage account not found"}
    }
)
async def delete_storage_account(
    resource_group: str,
    name: str,
    azure_service: AzureService = Depends(get_azure_service)
) -> JSONResponse:
    """Delete a storage account."""
    try:
        logger.info(f"API: Deleting storage account {name} in resource group {resource_group}")
        
        await azure_service.delete_storage_account(resource_group, name)
        
        logger.info(f"API: Successfully deleted storage account {name}")
        
        return JSONResponse(
            status_code=status.HTTP_200_OK,
            content={
                "message": f"Storage account '{name}' deleted successfully",
                "resource_group": resource_group,
                "name": name
            }
        )
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"API: Failed to delete storage account {name}: {e}")
        
        if "not found" in str(e).lower():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Storage account '{name}' not found in resource group '{resource_group}'"
            )
        
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to delete storage account: {str(e)}"
        )