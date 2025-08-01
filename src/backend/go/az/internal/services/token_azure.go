package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"brutus/azure-api/internal/models"
)

// TokenCredential implements azcore.TokenCredential using a provided access token
type TokenCredential struct {
	accessToken string
}

func NewTokenCredential(accessToken string) *TokenCredential {
	return &TokenCredential{
		accessToken: accessToken,
	}
}

func (t *TokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	// Return the provided token - in production you might want to validate expiry
	return azcore.AccessToken{
		Token: t.accessToken,
		ExpiresOn: time.Now().Add(time.Hour), // Set expiry to 1 hour from now (placeholder)
	}, nil
}

// TokenBasedAzureService uses user access tokens instead of service credentials
type TokenBasedAzureService struct{}

func NewTokenBasedAzureService() *TokenBasedAzureService {
	return &TokenBasedAzureService{}
}

func (s *TokenBasedAzureService) createClientWithToken(accessToken, subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID is required")
	}

	cred := NewTokenCredential(accessToken)
	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource groups client with user token: %w", err)
	}

	return client, nil
}

func (s *TokenBasedAzureService) createStorageClientWithToken(accessToken, subscriptionID string) (*armstorage.AccountsClient, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID is required")
	}

	cred := NewTokenCredential(accessToken)
	client, err := armstorage.NewAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage accounts client with user token: %w", err)
	}

	return client, nil
}

func (s *TokenBasedAzureService) ListResourceGroups(ctx context.Context, accessToken, subscriptionID string) ([]models.ResourceGroup, error) {
	client, err := s.createClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(nil)
	var resourceGroups []models.ResourceGroup

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, rg := range page.Value {
			resourceGroups = append(resourceGroups, models.ResourceGroup{
				ID:       getStringValue(rg.ID),
				Name:     getStringValue(rg.Name),
				Location: getStringValue(rg.Location),
				Tags:     convertTags(rg.Tags),
			})
		}
	}

	return resourceGroups, nil
}

func (s *TokenBasedAzureService) GetResourceGroup(ctx context.Context, accessToken, subscriptionID, name string) (*models.ResourceGroup, error) {
	client, err := s.createClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group %s: %w", name, err)
	}

	return &models.ResourceGroup{
		ID:       getStringValue(resp.ID),
		Name:     getStringValue(resp.Name),
		Location: getStringValue(resp.Location),
		Tags:     convertTags(resp.Tags),
	}, nil
}

func (s *TokenBasedAzureService) CreateResourceGroup(ctx context.Context, accessToken, subscriptionID string, req models.CreateResourceGroupRequest) (*models.ResourceGroup, error) {
	client, err := s.createClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	parameters := armresources.ResourceGroup{
		Location: &req.Location,
		Tags:     convertToAzureTags(req.Tags),
	}

	resp, err := client.CreateOrUpdate(ctx, req.Name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource group %s: %w", req.Name, err)
	}

	return &models.ResourceGroup{
		ID:       getStringValue(resp.ID),
		Name:     getStringValue(resp.Name),
		Location: getStringValue(resp.Location),
		Tags:     convertTags(resp.Tags),
	}, nil
}

func (s *TokenBasedAzureService) UpdateResourceGroup(ctx context.Context, accessToken, subscriptionID, name string, req models.UpdateResourceGroupRequest) (*models.ResourceGroup, error) {
	client, err := s.createClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	parameters := armresources.ResourceGroupPatchable{
		Tags: convertToAzureTags(req.Tags),
	}

	resp, err := client.Update(ctx, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update resource group %s: %w", name, err)
	}

	return &models.ResourceGroup{
		ID:       getStringValue(resp.ID),
		Name:     getStringValue(resp.Name),
		Location: getStringValue(resp.Location),
		Tags:     convertTags(resp.Tags),
	}, nil
}

func (s *TokenBasedAzureService) DeleteResourceGroup(ctx context.Context, accessToken, subscriptionID, name string) error {
	client, err := s.createClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return err
	}

	poller, err := client.BeginDelete(ctx, name, nil)
	if err != nil {
		return fmt.Errorf("failed to begin delete resource group %s: %w", name, err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete resource group %s: %w", name, err)
	}

	return nil
}

// Helper functions for converting Azure SDK types
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func convertTags(azureTags map[string]*string) map[string]string {
	if azureTags == nil {
		return nil
	}
	
	tags := make(map[string]string)
	for k, v := range azureTags {
		if v != nil {
			tags[k] = *v
		}
	}
	return tags
}

func convertToAzureTags(tags map[string]string) map[string]*string {
	if tags == nil {
		return nil
	}
	
	azureTags := make(map[string]*string)
	for k, v := range tags {
		value := v // Create a copy to get pointer
		azureTags[k] = &value
	}
	return azureTags
}

// Storage Account Methods

func (s *TokenBasedAzureService) ListStorageAccounts(ctx context.Context, accessToken, subscriptionID string) ([]models.StorageAccount, error) {
	client, err := s.createStorageClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(nil)
	var storageAccounts []models.StorageAccount

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		for _, sa := range page.Value {
			storageAccount := convertToStorageAccount(sa)
			storageAccounts = append(storageAccounts, *storageAccount)
		}
	}

	return storageAccounts, nil
}

func (s *TokenBasedAzureService) GetStorageAccount(ctx context.Context, accessToken, subscriptionID, resourceGroup, name string) (*models.StorageAccount, error) {
	client, err := s.createStorageClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetProperties(ctx, resourceGroup, name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage account %s: %w", name, err)
	}

	return convertToStorageAccount(&resp.Account), nil
}

func (s *TokenBasedAzureService) CreateStorageAccount(ctx context.Context, accessToken, subscriptionID string, req models.CreateStorageAccountRequest) (*models.StorageAccount, error) {
	client, err := s.createStorageClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Set defaults
	kind := armstorage.KindStorageV2
	if req.Kind != "" {
		switch req.Kind {
		case "StorageV2":
			kind = armstorage.KindStorageV2
		case "BlobStorage":
			kind = armstorage.KindBlobStorage
		case "Storage":
			kind = armstorage.KindStorage
		case "FileStorage":
			kind = armstorage.KindFileStorage
		case "BlockBlobStorage":
			kind = armstorage.KindBlockBlobStorage
		}
	}

	skuName := armstorage.SKUNameStandardLRS
	if req.SkuName != "" {
		switch req.SkuName {
		case "Standard_LRS":
			skuName = armstorage.SKUNameStandardLRS
		case "Standard_GRS":
			skuName = armstorage.SKUNameStandardGRS
		case "Standard_RAGRS":
			skuName = armstorage.SKUNameStandardRAGRS
		case "Premium_LRS":
			skuName = armstorage.SKUNamePremiumLRS
		}
	}

	parameters := armstorage.AccountCreateParameters{
		Location: &req.Location,
		Kind:     &kind,
		SKU: &armstorage.SKU{
			Name: &skuName,
		},
		Properties: &armstorage.AccountPropertiesCreateParameters{
			AllowBlobPublicAccess:  req.AllowBlobPublic,
			AllowSharedKeyAccess:   req.AllowSharedKey,
		},
		Tags: convertToAzureTags(req.Tags),
	}

	if req.AccessTier != "" {
		switch req.AccessTier {
		case "Hot":
			tier := armstorage.AccessTierHot
			parameters.Properties.AccessTier = &tier
		case "Cool":
			tier := armstorage.AccessTierCool
			parameters.Properties.AccessTier = &tier
		}
	}

	poller, err := client.BeginCreate(ctx, req.ResourceGroup, req.Name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin create storage account %s: %w", req.Name, err)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage account %s: %w", req.Name, err)
	}

	return convertToStorageAccount(&resp.Account), nil
}

func (s *TokenBasedAzureService) UpdateStorageAccount(ctx context.Context, accessToken, subscriptionID, resourceGroup, name string, req models.UpdateStorageAccountRequest) (*models.StorageAccount, error) {
	client, err := s.createStorageClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return nil, err
	}

	parameters := armstorage.AccountUpdateParameters{
		Properties: &armstorage.AccountPropertiesUpdateParameters{
			AllowBlobPublicAccess: req.AllowBlobPublic,
			AllowSharedKeyAccess:  req.AllowSharedKey,
		},
		Tags: convertToAzureTags(req.Tags),
	}

	if req.AccessTier != "" {
		switch req.AccessTier {
		case "Hot":
			tier := armstorage.AccessTierHot
			parameters.Properties.AccessTier = &tier
		case "Cool":
			tier := armstorage.AccessTierCool
			parameters.Properties.AccessTier = &tier
		}
	}

	resp, err := client.Update(ctx, resourceGroup, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update storage account %s: %w", name, err)
	}

	return convertToStorageAccount(&resp.Account), nil
}

func (s *TokenBasedAzureService) DeleteStorageAccount(ctx context.Context, accessToken, subscriptionID, resourceGroup, name string) error {
	client, err := s.createStorageClientWithToken(accessToken, subscriptionID)
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("failed to delete storage account %s: %w", name, err)
	}

	return nil
}

// Helper functions for storage accounts
func convertToStorageAccount(sa *armstorage.Account) *models.StorageAccount {
	account := &models.StorageAccount{
		ID:       getStringValue(sa.ID),
		Name:     getStringValue(sa.Name),
		Location: getStringValue(sa.Location),
		Tags:     convertTags(sa.Tags),
	}

	if sa.Properties != nil {
		props := sa.Properties
		if props.AccessTier != nil {
			account.AccessTier = string(*props.AccessTier)
		}
		if props.AllowBlobPublicAccess != nil {
			account.AllowBlobPublic = *props.AllowBlobPublicAccess
		}
		if props.AllowSharedKeyAccess != nil {
			account.AllowSharedKey = *props.AllowSharedKeyAccess
		}
		if props.CreationTime != nil {
			account.CreationTime = props.CreationTime.Format(time.RFC3339)
		}
		if props.PrimaryEndpoints != nil {
			account.PrimaryEndpoints = &models.StorageEndpoints{
				Blob:  getStringValue(props.PrimaryEndpoints.Blob),
				Queue: getStringValue(props.PrimaryEndpoints.Queue),
				Table: getStringValue(props.PrimaryEndpoints.Table),
				File:  getStringValue(props.PrimaryEndpoints.File),
			}
		}
	}

	if sa.Kind != nil {
		account.Kind = string(*sa.Kind)
	}

	if sa.SKU != nil {
		if sa.SKU.Name != nil {
			account.SkuName = string(*sa.SKU.Name)
		}
		if sa.SKU.Tier != nil {
			account.SkuTier = string(*sa.SKU.Tier)
		}
	}

	// Extract resource group from ID
	if account.ID != "" {
		parts := parseResourceID(account.ID)
		if rg, exists := parts["resourceGroups"]; exists {
			account.ResourceGroup = rg
		}
	}

	return account
}

// parseResourceID parses an Azure resource ID and returns key-value pairs
func parseResourceID(resourceID string) map[string]string {
	parts := make(map[string]string)
	segments := strings.Split(resourceID, "/")
	
	for i := 0; i < len(segments)-1; i += 2 {
		if i+1 < len(segments) {
			parts[segments[i]] = segments[i+1]
		}
	}
	
	return parts
}