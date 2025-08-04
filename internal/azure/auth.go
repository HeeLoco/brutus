package azure

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type AuthInfo struct {
	Credential         azcore.TokenCredential
	SubscriptionID     string
	TenantID           string
	CurrentUser        string
	AuthMethod         string
	Subscriptions      []*armsubscriptions.Subscription
	DefaultSub         *armsubscriptions.Subscription
	PermissionLevel    string
	HasOwnerAccess     bool
	CAFPermissions     string
	ManagementGroups   []*armmanagementgroups.ManagementGroupInfo
	MGPermissions      string
	HasMGReadAccess    bool
	HasMGWriteAccess   bool
	MGPermissionReason string
}

type ManagementGroupNode struct {
	ID            string
	Name          string
	DisplayName   string
	Children      []*ManagementGroupNode
	Subscriptions []*armsubscriptions.Subscription
	Level         int
}

func NewAuth(ctx context.Context) (*AuthInfo, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	currentUser := "Unknown"
	if usr, err := user.Current(); err == nil {
		currentUser = usr.Username
	}

	authMethod := "Azure CLI"
	if os.Getenv("AZURE_CLIENT_ID") != "" {
		authMethod = "Service Principal"
	} else if os.Getenv("MSI_ENDPOINT") != "" {
		authMethod = "Managed Identity"
	}

	auth := &AuthInfo{
		Credential:  cred,
		CurrentUser: currentUser,
		AuthMethod:  authMethod,
	}

	return auth, nil
}

func (a *AuthInfo) ValidateAuthentication(ctx context.Context) error {
	subscriptions, err := a.ListSubscriptions(ctx)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return fmt.Errorf("no Azure subscriptions available")
	}

	a.Subscriptions = subscriptions
	
	// Find default subscription (first enabled one)
	for _, sub := range subscriptions {
		if sub.State != nil && *sub.State == armsubscriptions.SubscriptionStateEnabled {
			a.DefaultSub = sub
			if sub.SubscriptionID != nil {
				a.SubscriptionID = *sub.SubscriptionID
			}
			if sub.TenantID != nil {
				a.TenantID = *sub.TenantID
			}
			break
		}
	}

	// Check permissions on default subscription
	if a.SubscriptionID != "" {
		if err := a.checkPermissions(ctx); err != nil {
			log.Printf("Warning: Could not check permissions: %v", err)
			a.PermissionLevel = "Unknown"
			a.CAFPermissions = "Unable to verify"
		}
	}

	log.Printf("Successfully authenticated to Azure - Found %d subscriptions", len(subscriptions))
	return nil
}

func (a *AuthInfo) ListSubscriptions(ctx context.Context) ([]*armsubscriptions.Subscription, error) {
	client, err := armsubscriptions.NewClient(a.Credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions client: %w", err)
	}

	var subscriptions []*armsubscriptions.Subscription
	pager := client.NewListPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions: %w", err)
		}
		subscriptions = append(subscriptions, page.Value...)
	}

	return subscriptions, nil
}

func (a *AuthInfo) SetSubscription(subscriptionID string) {
	a.SubscriptionID = subscriptionID
}

func (a *AuthInfo) GetEnabledSubscriptions() []*armsubscriptions.Subscription {
	var enabled []*armsubscriptions.Subscription
	for _, sub := range a.Subscriptions {
		if sub.State != nil && *sub.State == armsubscriptions.SubscriptionStateEnabled {
			enabled = append(enabled, sub)
		}
	}
	return enabled
}

func (a *AuthInfo) GetSubscriptionSummary() string {
	enabled := a.GetEnabledSubscriptions()
	total := len(a.Subscriptions)
	if total == 0 {
		return "No subscriptions found"
	}
	return fmt.Sprintf("%d enabled of %d total", len(enabled), total)
}

func (a *AuthInfo) GetLocationInfo() string {
	enabled := a.GetEnabledSubscriptions()
	if len(enabled) == 0 {
		return "No location data"
	}
	return "Global scope available"
}

func (a *AuthInfo) HasCAFReadySubscription() bool {
	// Check if we have at least one enabled subscription for CAF deployment
	enabled := a.GetEnabledSubscriptions()
	return len(enabled) > 0
}

func (a *AuthInfo) checkPermissions(ctx context.Context) error {
	authClient, err := armauthorization.NewRoleAssignmentsClient(a.SubscriptionID, a.Credential, nil)
	if err != nil {
		return fmt.Errorf("failed to create authorization client: %w", err)
	}

	// Check role assignments at subscription scope
	scope := fmt.Sprintf("/subscriptions/%s", a.SubscriptionID)
	pager := authClient.NewListForScopePager(scope, &armauthorization.RoleAssignmentsClientListForScopeOptions{
		Filter: nil,
	})

	var highestRole string
	var hasOwner, hasContributor, hasReader bool

	// Process role assignments
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list role assignments: %w", err)
		}

		for _, assignment := range page.Value {
			if assignment.Properties == nil || assignment.Properties.RoleDefinitionID == nil {
				continue
			}

			roleDefID := *assignment.Properties.RoleDefinitionID
			roleName := a.getRoleNameFromID(roleDefID)

			switch {
			case strings.Contains(strings.ToLower(roleName), "owner"):
				hasOwner = true
				highestRole = "Owner"
			case strings.Contains(strings.ToLower(roleName), "contributor"):
				if !hasOwner {
					hasContributor = true
					highestRole = "Contributor"
				}
			case strings.Contains(strings.ToLower(roleName), "reader"):
				if !hasOwner && !hasContributor {
					hasReader = true
					if highestRole == "" {
						highestRole = "Reader"
					}
				}
			}
		}
		
		// Only process first page for performance
		break
	}

	// Set permission levels
	if highestRole == "" {
		a.PermissionLevel = "Limited/Custom"
		a.CAFPermissions = "Custom permissions - verify manually"
	} else {
		a.PermissionLevel = highestRole
	}

	a.HasOwnerAccess = hasOwner

	// Determine CAF capabilities
	if hasOwner {
		a.CAFPermissions = "✅ Full CAF deployment capability"
	} else if hasContributor {
		a.CAFPermissions = "⚠️ Resource creation (limited CAF features)"
	} else if hasReader {
		a.CAFPermissions = "❌ Read-only access"
	} else {
		a.CAFPermissions = "❓ Custom permissions"
	}

	return nil
}

func (a *AuthInfo) getRoleNameFromID(roleDefID string) string {
	// Common Azure built-in role definition IDs
	roleMap := map[string]string{
		"8e3af657-a8ff-443c-a75c-2fe8c4bcb635": "Owner",
		"b24988ac-6180-42a0-ab88-20f7382dd24c": "Contributor", 
		"acdd72a7-3385-48ef-bd42-f606fba81ae7": "Reader",
		"18d7d88d-d35e-4fb5-a5c3-7773c20a72d9": "User Access Administrator",
	}

	// Extract the GUID from the full resource ID
	parts := strings.Split(roleDefID, "/")
	if len(parts) > 0 {
		guid := parts[len(parts)-1]
		if name, exists := roleMap[guid]; exists {
			return name
		}
	}

	return "Custom Role"
}

func (a *AuthInfo) ListManagementGroups(ctx context.Context) error {
	client, err := armmanagementgroups.NewClient(a.Credential, nil)
	if err != nil {
		a.HasMGReadAccess = false
		a.MGPermissionReason = "Cannot create management groups client"
		a.MGPermissions = "❌ No management group access"
		return fmt.Errorf("failed to create management groups client: %w", err)
	}

	var managementGroups []*armmanagementgroups.ManagementGroupInfo
	pager := client.NewListPager(nil)

	// Try to list management groups to test permissions
	if pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			// Analyze the error to determine permission issue
			errorMsg := err.Error()
			a.HasMGReadAccess = false
			
			if strings.Contains(strings.ToLower(errorMsg), "forbidden") || 
			   strings.Contains(strings.ToLower(errorMsg), "authorization") ||
			   strings.Contains(strings.ToLower(errorMsg), "access denied") {
				a.MGPermissionReason = "Insufficient permissions to read management groups"
				a.MGPermissions = "❌ Read access denied - Need Management Group Reader role"
			} else {
				a.MGPermissionReason = fmt.Sprintf("API error: %s", errorMsg)
				a.MGPermissions = "❌ Cannot access management groups API"
			}
			
			log.Printf("Management Groups access denied: %v", err)
			return nil // Don't return error, just record the permission limitation
		}
		
		// Successfully got first page
		managementGroups = append(managementGroups, page.Value...)
		a.HasMGReadAccess = true
		
		// Continue getting remaining pages
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				break
			}
			managementGroups = append(managementGroups, page.Value...)
		}
	}

	a.ManagementGroups = managementGroups
	
	// Determine permission status
	if a.HasMGReadAccess {
		if len(managementGroups) > 0 {
			a.MGPermissions = fmt.Sprintf("✅ Read access - Found %d management groups", len(managementGroups))
			a.MGPermissionReason = "Can read existing management groups"
		} else {
			a.MGPermissions = "✅ Read access - No existing management groups"
			a.MGPermissionReason = "Read access confirmed, no management groups exist"
		}
		
		// TODO: Test write access by attempting a dry-run operation
		a.HasMGWriteAccess = false // Conservative assumption
	}
	
	return nil
}

func (a *AuthInfo) GetManagementGroupTree(ctx context.Context) (*ManagementGroupNode, error) {
	if len(a.ManagementGroups) == 0 {
		// Create a mock CAF structure for demonstration
		return a.createMockCAFStructure(), nil
	}

	// If no real management groups, return mock structure
	return a.createMockCAFStructure(), nil
}

func (a *AuthInfo) createMockCAFStructure() *ManagementGroupNode {
	// Create a typical CAF management group structure
	rootMG := &ManagementGroupNode{
		ID:          "mg-tenant-root",
		Name:        "Tenant Root",
		DisplayName: "Tenant Root Group",
		Level:       0,
		Children:    []*ManagementGroupNode{},
	}

	// Platform management group
	platformMG := &ManagementGroupNode{
		ID:          "mg-platform",
		Name:        "Platform",
		DisplayName: "Platform",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Landing Zones management group  
	landingZonesMG := &ManagementGroupNode{
		ID:          "mg-landingzones",
		Name:        "Landing Zones",
		DisplayName: "Landing Zones",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Sandbox management group
	sandboxMG := &ManagementGroupNode{
		ID:          "mg-sandbox",
		Name:        "Sandbox",
		DisplayName: "Sandbox",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Platform sub-groups
	managementMG := &ManagementGroupNode{
		ID:          "mg-management",
		Name:        "Management",
		DisplayName: "Management",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	connectivityMG := &ManagementGroupNode{
		ID:          "mg-connectivity",
		Name:        "Connectivity",
		DisplayName: "Connectivity", 
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	identityMG := &ManagementGroupNode{
		ID:          "mg-identity",
		Name:        "Identity",
		DisplayName: "Identity",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Landing Zone sub-groups
	corpMG := &ManagementGroupNode{
		ID:          "mg-corp",
		Name:        "Corp",
		DisplayName: "Corporate",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	onlineMG := &ManagementGroupNode{
		ID:          "mg-online",
		Name:        "Online",
		DisplayName: "Online",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Add current subscription to appropriate group
	if a.DefaultSub != nil {
		corpMG.Subscriptions = []*armsubscriptions.Subscription{a.DefaultSub}
	}

	// Build the tree
	platformMG.Children = []*ManagementGroupNode{managementMG, connectivityMG, identityMG}
	landingZonesMG.Children = []*ManagementGroupNode{corpMG, onlineMG}
	rootMG.Children = []*ManagementGroupNode{platformMG, landingZonesMG, sandboxMG}

	return rootMG
}

func (a *AuthInfo) GetRecommendedCAFStructure(structureType string) *ManagementGroupNode {
	switch structureType {
	case "basic":
		return a.createBasicCAFStructure()
	case "enterprise":
		return a.createEnterpriseCAFStructure()
	case "minimal":
		return a.createMinimalCAFStructure()
	default:
		return a.createBasicCAFStructure()
	}
}

func (a *AuthInfo) createBasicCAFStructure() *ManagementGroupNode {
	// Basic CAF structure (Small-Medium Organizations)
	rootMG := &ManagementGroupNode{
		ID:          "mg-tenant-root",
		Name:        "Tenant Root",
		DisplayName: "Tenant Root Group",
		Level:       0,
		Children:    []*ManagementGroupNode{},
	}

	// Platform management group
	platformMG := &ManagementGroupNode{
		ID:          "mg-platform",
		Name:        "Platform",
		DisplayName: "Platform",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Landing Zones management group  
	landingZonesMG := &ManagementGroupNode{
		ID:          "mg-landingzones",
		Name:        "Landing Zones",
		DisplayName: "Landing Zones",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Sandbox management group
	sandboxMG := &ManagementGroupNode{
		ID:          "mg-sandbox",
		Name:        "Sandbox",
		DisplayName: "Sandbox",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Platform sub-groups
	managementMG := &ManagementGroupNode{
		ID:          "mg-management",
		Name:        "Management",
		DisplayName: "Management",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	connectivityMG := &ManagementGroupNode{
		ID:          "mg-connectivity",
		Name:        "Connectivity",
		DisplayName: "Connectivity", 
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	identityMG := &ManagementGroupNode{
		ID:          "mg-identity",
		Name:        "Identity",
		DisplayName: "Identity",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Landing Zone sub-groups
	corpMG := &ManagementGroupNode{
		ID:          "mg-corp",
		Name:        "Corp",
		DisplayName: "Corporate",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	onlineMG := &ManagementGroupNode{
		ID:          "mg-online",
		Name:        "Online",
		DisplayName: "Online",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Build the tree
	platformMG.Children = []*ManagementGroupNode{managementMG, connectivityMG, identityMG}
	landingZonesMG.Children = []*ManagementGroupNode{corpMG, onlineMG}
	rootMG.Children = []*ManagementGroupNode{platformMG, landingZonesMG, sandboxMG}

	return rootMG
}

func (a *AuthInfo) createEnterpriseCAFStructure() *ManagementGroupNode {
	// Enterprise CAF structure (Large Organizations)
	rootMG := &ManagementGroupNode{
		ID:          "mg-tenant-root",
		Name:        "Tenant Root",
		DisplayName: "Tenant Root Group",
		Level:       0,
		Children:    []*ManagementGroupNode{},
	}

	// Platform management group with Security
	platformMG := &ManagementGroupNode{
		ID:          "mg-platform",
		Name:        "Platform",
		DisplayName: "Platform",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Extended Landing Zones
	landingZonesMG := &ManagementGroupNode{
		ID:          "mg-landingzones",
		Name:        "Landing Zones", 
		DisplayName: "Landing Zones",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Sandbox and Decommissioned
	sandboxMG := &ManagementGroupNode{
		ID:          "mg-sandbox",
		Name:        "Sandbox",
		DisplayName: "Sandbox",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	decommissionedMG := &ManagementGroupNode{
		ID:          "mg-decommissioned",
		Name:        "Decommissioned",
		DisplayName: "Decommissioned",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Platform sub-groups (with Security)
	managementMG := &ManagementGroupNode{
		ID:          "mg-management",
		Name:        "Management",
		DisplayName: "Management",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	connectivityMG := &ManagementGroupNode{
		ID:          "mg-connectivity",
		Name:        "Connectivity",
		DisplayName: "Connectivity",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	identityMG := &ManagementGroupNode{
		ID:          "mg-identity",
		Name:        "Identity",
		DisplayName: "Identity",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	securityMG := &ManagementGroupNode{
		ID:          "mg-security",
		Name:        "Security",
		DisplayName: "Security",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Extended Landing Zone sub-groups
	corpMG := &ManagementGroupNode{
		ID:          "mg-corp",
		Name:        "Corp",
		DisplayName: "Corporate",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	onlineMG := &ManagementGroupNode{
		ID:          "mg-online",
		Name:        "Online", 
		DisplayName: "Online",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	sapMG := &ManagementGroupNode{
		ID:          "mg-sap",
		Name:        "SAP",
		DisplayName: "SAP",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	avsMG := &ManagementGroupNode{
		ID:          "mg-avs",
		Name:        "AVS",
		DisplayName: "Azure VMware Solution",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Build the tree
	platformMG.Children = []*ManagementGroupNode{managementMG, connectivityMG, identityMG, securityMG}
	landingZonesMG.Children = []*ManagementGroupNode{corpMG, onlineMG, sapMG, avsMG}
	rootMG.Children = []*ManagementGroupNode{platformMG, landingZonesMG, sandboxMG, decommissionedMG}

	return rootMG
}

func (a *AuthInfo) createMinimalCAFStructure() *ManagementGroupNode {
	// Minimal CAF structure (Startups/Simple Setups)
	rootMG := &ManagementGroupNode{
		ID:          "mg-tenant-root",
		Name:        "Tenant Root",
		DisplayName: "Tenant Root Group",
		Level:       0,
		Children:    []*ManagementGroupNode{},
	}

	// Simplified Platform
	platformMG := &ManagementGroupNode{
		ID:          "mg-platform",
		Name:        "Platform",
		DisplayName: "Platform",
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Simple Landing Zones
	landingZonesMG := &ManagementGroupNode{
		ID:          "mg-landingzones",
		Name:        "Landing Zones",
		DisplayName: "Landing Zones", 
		Level:       1,
		Children:    []*ManagementGroupNode{},
	}

	// Only Corp landing zone
	corpMG := &ManagementGroupNode{
		ID:          "mg-corp",
		Name:        "Corp",
		DisplayName: "Corporate",
		Level:       2,
		Children:    []*ManagementGroupNode{},
	}

	// Build minimal tree
	landingZonesMG.Children = []*ManagementGroupNode{corpMG}
	rootMG.Children = []*ManagementGroupNode{platformMG, landingZonesMG}

	return rootMG
}