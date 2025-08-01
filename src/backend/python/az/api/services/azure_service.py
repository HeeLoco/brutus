from typing import List, Optional

from azure.core.exceptions import ResourceNotFoundError, HttpResponseError
from azure.identity import DefaultAzureCredential, ClientSecretCredential
from azure.mgmt.resource import ResourceManagementClient
from azure.mgmt.storage import StorageManagementClient

from api.core.config import Settings
from api.core.logging import get_logger
from api.models.resource_group import (
    ResourceGroup, ResourceGroupCreate, ResourceGroupUpdate,
    StorageAccount, StorageAccountCreate, StorageAccountUpdate, StorageEndpoints
)


logger = get_logger(__name__)


class AzureService:
    """Azure service for managing Azure resources."""
    
    def __init__(self, settings: Settings):
        self.settings = settings
        self.resource_client = self._create_resource_client()
        self.storage_client = self._create_storage_client()
    
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
    
    def _create_storage_client(self) -> StorageManagementClient:
        """Create Azure storage management client with proper authentication."""
        try:
            # Use the same credential as resource client
            if all([
                self.settings.azure_client_id,
                self.settings.azure_client_secret,
                self.settings.azure_tenant_id
            ]):
                logger.info("Using service principal authentication for storage client")
                credential = ClientSecretCredential(
                    tenant_id=self.settings.azure_tenant_id,
                    client_id=self.settings.azure_client_id,
                    client_secret=self.settings.azure_client_secret
                )
            else:
                logger.info("Using default Azure credential for storage client")
                credential = DefaultAzureCredential()
            
            client = StorageManagementClient(
                credential=credential,
                subscription_id=self.settings.azure_subscription_id
            )
            
            logger.info("Successfully created Azure storage client")
            return client
            
        except Exception as e:
            logger.error(f"Failed to create Azure storage client: {e}")
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
    
    # Storage Account Methods
    
    async def list_storage_accounts(self) -> List[StorageAccount]:
        """List all storage accounts in the subscription."""
        try:
            logger.info("Listing storage accounts")
            storage_accounts = []
            
            for sa in self.storage_client.storage_accounts.list():
                # Extract resource group from ID
                resource_group = self._extract_resource_group_from_id(sa.id)
                
                endpoints = None
                if sa.primary_endpoints:
                    endpoints = StorageEndpoints(
                        blob=sa.primary_endpoints.blob,
                        queue=sa.primary_endpoints.queue,
                        table=sa.primary_endpoints.table,
                        file=sa.primary_endpoints.file
                    )
                
                storage_accounts.append(StorageAccount(
                    id=sa.id,
                    name=sa.name,
                    location=sa.location,
                    resource_group=resource_group,
                    kind=sa.kind,
                    sku_name=sa.sku.name if sa.sku else None,
                    sku_tier=sa.sku.tier if sa.sku else None,
                    access_tier=sa.access_tier,
                    allow_blob_public=sa.allow_blob_public_access,
                    allow_shared_key=sa.allow_shared_key_access,
                    tags=sa.tags or {},
                    creation_time=sa.creation_time,
                    primary_endpoints=endpoints
                ))
            
            logger.info(f"Found {len(storage_accounts)} storage accounts")
            return storage_accounts
            
        except Exception as e:
            logger.error(f"Failed to list storage accounts: {e}")
            raise
    
    async def get_storage_account(self, resource_group: str, name: str) -> StorageAccount:
        """Get a specific storage account by name."""
        try:
            logger.info(f"Getting storage account: {name} in resource group: {resource_group}")
            
            sa = self.storage_client.storage_accounts.get_properties(resource_group, name)
            
            endpoints = None
            if sa.primary_endpoints:
                endpoints = StorageEndpoints(
                    blob=sa.primary_endpoints.blob,
                    queue=sa.primary_endpoints.queue,
                    table=sa.primary_endpoints.table,
                    file=sa.primary_endpoints.file
                )
            
            return StorageAccount(
                id=sa.id,
                name=sa.name,
                location=sa.location,
                resource_group=resource_group,
                kind=sa.kind,
                sku_name=sa.sku.name if sa.sku else None,
                sku_tier=sa.sku.tier if sa.sku else None,
                access_tier=sa.access_tier,
                allow_blob_public=sa.allow_blob_public_access,
                allow_shared_key=sa.allow_shared_key_access,
                tags=sa.tags or {},
                creation_time=sa.creation_time,
                primary_endpoints=endpoints
            )
            
        except ResourceNotFoundError:
            logger.warning(f"Storage account not found: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to get storage account {name}: {e}")
            raise
    
    async def create_storage_account(self, storage_account: StorageAccountCreate) -> StorageAccount:
        """Create a new storage account."""
        try:
            logger.info(f"Creating storage account: {storage_account.name}")
            
            from azure.mgmt.storage.models import (
                StorageAccountCreateParameters, Sku, Kind, AccessTier
            )
            
            # Set defaults and validate
            kind = Kind.STORAGE_V2
            if storage_account.kind:
                kind = getattr(Kind, storage_account.kind.upper().replace('V', '_V'), Kind.STORAGE_V2)
            
            sku = Sku(name=storage_account.sku_name or "Standard_LRS")
            
            parameters = StorageAccountCreateParameters(
                location=storage_account.location,
                kind=kind,
                sku=sku,
                tags=storage_account.tags or {}
            )
            
            # Set optional parameters
            if storage_account.access_tier:
                parameters.access_tier = getattr(AccessTier, storage_account.access_tier.upper(), None)
            
            if storage_account.allow_blob_public is not None:
                parameters.allow_blob_public_access = storage_account.allow_blob_public
                
            if storage_account.allow_shared_key is not None:
                parameters.allow_shared_key_access = storage_account.allow_shared_key
            
            # Start creation (this is an async operation)
            create_operation = self.storage_client.storage_accounts.begin_create(
                storage_account.resource_group,
                storage_account.name,
                parameters
            )
            
            # Wait for completion
            result = create_operation.result()
            
            logger.info(f"Successfully created storage account: {storage_account.name}")
            
            endpoints = None
            if result.primary_endpoints:
                endpoints = StorageEndpoints(
                    blob=result.primary_endpoints.blob,
                    queue=result.primary_endpoints.queue,
                    table=result.primary_endpoints.table,
                    file=result.primary_endpoints.file
                )
            
            return StorageAccount(
                id=result.id,
                name=result.name,
                location=result.location,
                resource_group=storage_account.resource_group,
                kind=result.kind,
                sku_name=result.sku.name if result.sku else None,
                sku_tier=result.sku.tier if result.sku else None,
                access_tier=result.access_tier,
                allow_blob_public=result.allow_blob_public_access,
                allow_shared_key=result.allow_shared_key_access,
                tags=result.tags or {},
                creation_time=result.creation_time,
                primary_endpoints=endpoints
            )
            
        except HttpResponseError as e:
            logger.error(f"Azure API error creating storage account {storage_account.name}: {e}")
            raise
        except Exception as e:
            logger.error(f"Failed to create storage account {storage_account.name}: {e}")
            raise
    
    async def update_storage_account(self, resource_group: str, name: str, update: StorageAccountUpdate) -> StorageAccount:
        """Update an existing storage account."""
        try:
            logger.info(f"Updating storage account: {name}")
            
            from azure.mgmt.storage.models import StorageAccountUpdateParameters, AccessTier
            
            # First check if storage account exists
            await self.get_storage_account(resource_group, name)
            
            parameters = StorageAccountUpdateParameters()
            
            if update.tags is not None:
                parameters.tags = update.tags
                
            if update.access_tier:
                parameters.access_tier = getattr(AccessTier, update.access_tier.upper(), None)
                
            if update.allow_blob_public is not None:
                parameters.allow_blob_public_access = update.allow_blob_public
                
            if update.allow_shared_key is not None:
                parameters.allow_shared_key_access = update.allow_shared_key
            
            result = self.storage_client.storage_accounts.update(resource_group, name, parameters)
            
            logger.info(f"Successfully updated storage account: {name}")
            
            endpoints = None
            if result.primary_endpoints:
                endpoints = StorageEndpoints(
                    blob=result.primary_endpoints.blob,
                    queue=result.primary_endpoints.queue,
                    table=result.primary_endpoints.table,
                    file=result.primary_endpoints.file
                )
            
            return StorageAccount(
                id=result.id,
                name=result.name,
                location=result.location,
                resource_group=resource_group,
                kind=result.kind,
                sku_name=result.sku.name if result.sku else None,
                sku_tier=result.sku.tier if result.sku else None,
                access_tier=result.access_tier,
                allow_blob_public=result.allow_blob_public_access,
                allow_shared_key=result.allow_shared_key_access,
                tags=result.tags or {},
                creation_time=result.creation_time,
                primary_endpoints=endpoints
            )
            
        except ResourceNotFoundError:
            logger.warning(f"Cannot update non-existent storage account: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to update storage account {name}: {e}")
            raise
    
    async def delete_storage_account(self, resource_group: str, name: str) -> None:
        """Delete a storage account."""
        try:
            logger.info(f"Deleting storage account: {name}")
            
            # Check if storage account exists first
            await self.get_storage_account(resource_group, name)
            
            # Delete the storage account
            self.storage_client.storage_accounts.delete(resource_group, name)
            
            logger.info(f"Successfully deleted storage account: {name}")
            
        except ResourceNotFoundError:
            logger.warning(f"Cannot delete non-existent storage account: {name}")
            raise
        except Exception as e:
            logger.error(f"Failed to delete storage account {name}: {e}")
            raise
    
    def _extract_resource_group_from_id(self, resource_id: str) -> str:
        """Extract resource group name from Azure resource ID."""
        if not resource_id:
            return ""
        
        parts = resource_id.split('/')
        try:
            # Azure resource ID format: /subscriptions/{sub}/resourceGroups/{rg}/...
            rg_index = parts.index('resourceGroups')
            if rg_index + 1 < len(parts):
                return parts[rg_index + 1]
        except (ValueError, IndexError):
            pass
        
        return ""