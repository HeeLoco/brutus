package services

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

	"brutus/azure-api/configs"
	"brutus/azure-api/internal/models"
)

type AzureService struct {
	client         *armresources.ResourceGroupsClient
	subscriptionID string
}

func NewAzureService(cfg *configs.Config) (*AzureService, error) {
	var cred azcore.TokenCredential
	var err error

	// Try different authentication methods
	if cfg.ClientID != "" && cfg.ClientSecret != "" && cfg.TenantID != "" {
		// Use client credentials
		cred, err = azidentity.NewClientSecretCredential(cfg.TenantID, cfg.ClientID, cfg.ClientSecret, nil)
	} else {
		// Use default Azure credential (managed identity, Azure CLI, etc.)
		cred, err = azidentity.NewDefaultAzureCredential(nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to obtain Azure credential: %w", err)
	}

	client, err := armresources.NewResourceGroupsClient(cfg.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource groups client: %w", err)
	}

	return &AzureService{
		client:         client,
		subscriptionID: cfg.SubscriptionID,
	}, nil
}

func (s *AzureService) ListResourceGroups(ctx context.Context) ([]models.ResourceGroup, error) {
	pager := s.client.NewListPager(nil)
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

func (s *AzureService) GetResourceGroup(ctx context.Context, name string) (*models.ResourceGroup, error) {
	resp, err := s.client.Get(ctx, name, nil)
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

func (s *AzureService) CreateResourceGroup(ctx context.Context, req models.CreateResourceGroupRequest) (*models.ResourceGroup, error) {
	parameters := armresources.ResourceGroup{
		Location: &req.Location,
		Tags:     convertToAzureTags(req.Tags),
	}

	resp, err := s.client.CreateOrUpdate(ctx, req.Name, parameters, nil)
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

func (s *AzureService) UpdateResourceGroup(ctx context.Context, name string, req models.UpdateResourceGroupRequest) (*models.ResourceGroup, error) {
	parameters := armresources.ResourceGroupPatchable{
		Tags: convertToAzureTags(req.Tags),
	}

	resp, err := s.client.Update(ctx, name, parameters, nil)
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

func (s *AzureService) DeleteResourceGroup(ctx context.Context, name string) error {
	poller, err := s.client.BeginDelete(ctx, name, nil)
	if err != nil {
		return fmt.Errorf("failed to begin delete resource group %s: %w", name, err)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete resource group %s: %w", name, err)
	}

	return nil
}

// Helper functions
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
	for key, value := range azureTags {
		if value != nil {
			tags[key] = *value
		}
	}
	return tags
}

func convertToAzureTags(tags map[string]string) map[string]*string {
	if tags == nil {
		return nil
	}

	azureTags := make(map[string]*string)
	for key, value := range tags {
		azureTags[key] = &value
	}
	return azureTags
}