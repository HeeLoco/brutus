from typing import List, Optional

from azure.core.exceptions import ResourceNotFoundError, HttpResponseError
from azure.identity import DefaultAzureCredential, ClientSecretCredential
from azure.mgmt.resource import ResourceManagementClient

from api.core.config import Settings
from api.core.logging import get_logger
from api.models.resource_group import ResourceGroup, ResourceGroupCreate, ResourceGroupUpdate


logger = get_logger(__name__)


class AzureService:
    """Azure service for managing Azure resources."""
    
    def __init__(self, settings: Settings):
        self.settings = settings
        self.resource_client = self._create_resource_client()
    
    def _create_resource_client(self) -> ResourceManagementClient:
        """Create Azure resource management client with proper authentication."""
        try:
            # Try service principal authentication first
            if all([
                self.settings.azure_client_id,
                self.settings.azure_client_secret,
                self.settings.azure_tenant_id
            ]):
                logger.info("Using service principal authentication")
                credential = ClientSecretCredential(
                    tenant_id=self.settings.azure_tenant_id,
                    client_id=self.settings.azure_client_id,
                    client_secret=self.settings.azure_client_secret
                )
            else:
                # Fall back to default credential (managed identity, Azure CLI, etc.)
                logger.info("Using default Azure credential")
                credential = DefaultAzureCredential()
            
            client = ResourceManagementClient(
                credential=credential,
                subscription_id=self.settings.azure_subscription_id
            )
            
            logger.info("Successfully created Azure resource client")
            return client
            
        except Exception as e:
            logger.error(f"Failed to create Azure resource client: {e}")
            raise
    
    async def test_connection(self) -> bool:
        """Test Azure connection."""
        try:
            # Try to list resource groups (limited to 1) to test connectivity
            list(self.resource_client.resource_groups.list(top=1))
            logger.info("Azure connectivity test successful")
            return True
        except Exception as e:
            logger.error(f"Azure connectivity test failed: {e}")
            return False
    
    async def list_resource_groups(self) -> List[ResourceGroup]:
        """List all resource groups in the subscription."""
        try:
            logger.info("Listing resource groups")
            resource_groups = []
            
            for rg in self.resource_client.resource_groups.list():
                resource_groups.append(ResourceGroup(
                    id=rg.id,
                    name=rg.name,
                    location=rg.location,
                    tags=rg.tags or {},
                    type=rg.type,
                    managed_by=rg.managed_by
                ))
            
            logger.info(f"Found {len(resource_groups)} resource groups")
            return resource_groups
            
        except Exception as e:
            logger.error(f"Failed to list resource groups: {e}")
            raise
    
    async def get_resource_group(self, name: str) -> ResourceGroup:
        """Get a specific resource group by name."""
        try:
            logger.info(f"Getting resource group: {name}")
            
            rg = self.resource_client.resource_groups.get(name)
            
            return ResourceGroup(
                id=rg.id,
                name=rg.name,
                location=rg.location,
                tags=rg.tags or {},
                type=rg.type,
                managed_by=rg.managed_by
            )
            
        except ResourceNotFoundError:
            logger.warning(f"Resource group not found: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to get resource group {name}: {e}")
            raise
    
    async def create_resource_group(self, resource_group: ResourceGroupCreate) -> ResourceGroup:
        """Create a new resource group."""
        try:
            logger.info(f"Creating resource group: {resource_group.name}")
            
            parameters = {
                "location": resource_group.location,
                "tags": resource_group.tags or {}
            }
            
            result = self.resource_client.resource_groups.create_or_update(
                resource_group.name,
                parameters
            )
            
            logger.info(f"Successfully created resource group: {resource_group.name}")
            
            return ResourceGroup(
                id=result.id,
                name=result.name,
                location=result.location,
                tags=result.tags or {},
                type=result.type,
                managed_by=result.managed_by
            )
            
        except HttpResponseError as e:
            logger.error(f"Azure API error creating resource group {resource_group.name}: {e}")
            raise
        except Exception as e:
            logger.error(f"Failed to create resource group {resource_group.name}: {e}")
            raise
    
    async def update_resource_group(self, name: str, update: ResourceGroupUpdate) -> ResourceGroup:
        """Update an existing resource group."""
        try:
            logger.info(f"Updating resource group: {name}")
            
            # First check if resource group exists
            await self.get_resource_group(name)
            
            # Update only tags (location cannot be changed)
            parameters = {
                "tags": update.tags or {}
            }
            
            result = self.resource_client.resource_groups.update(name, parameters)
            
            logger.info(f"Successfully updated resource group: {name}")
            
            return ResourceGroup(
                id=result.id,
                name=result.name,
                location=result.location,
                tags=result.tags or {},
                type=result.type,
                managed_by=result.managed_by
            )
            
        except ResourceNotFoundError:
            logger.warning(f"Cannot update non-existent resource group: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to update resource group {name}: {e}")
            raise
    
    async def delete_resource_group(self, name: str) -> None:
        """Delete a resource group."""
        try:
            logger.info(f"Deleting resource group: {name}")
            
            # Check if resource group exists first
            await self.get_resource_group(name)
            
            # Start the deletion (this is an async operation in Azure)
            delete_operation = self.resource_client.resource_groups.begin_delete(name)
            
            logger.info(f"Started deletion of resource group: {name}")
            # Note: We don't wait for completion as it can take a long time
            
        except ResourceNotFoundError:
            logger.warning(f"Cannot delete non-existent resource group: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to delete resource group {name}: {e}")
            raise