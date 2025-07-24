package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

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