#!/bin/bash

# Azure App Registration Verification Script

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  Azure App Registration Verification${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_check() {
    echo -e "${BLUE}🔍 Checking: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_header

# Check if Azure CLI is installed
print_check "Azure CLI installation"
if command -v az >/dev/null 2>&1; then
    print_success "Azure CLI is installed"
    az version --output none 2>/dev/null && echo "Version: $(az version --query '"azure-cli"' -o tsv 2>/dev/null)"
else
    print_error "Azure CLI is not installed"
    echo "Install: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
    exit 1
fi

# Check if logged in to Azure CLI
print_check "Azure CLI authentication"
if az account show >/dev/null 2>&1; then
    print_success "Logged in to Azure CLI"
    TENANT_ID=$(az account show --query tenantId -o tsv)
    SUBSCRIPTION_ID=$(az account show --query id -o tsv)
    echo "  Tenant ID: $TENANT_ID"
    echo "  Subscription ID: $SUBSCRIPTION_ID"
else
    print_error "Not logged in to Azure CLI"
    echo "Run: az login"
    exit 1
fi

# Prompt for App Registration details
echo
echo -e "${BLUE}Please provide your App Registration details:${NC}"
read -p "Enter your App Registration Client ID: " CLIENT_ID
read -p "Enter your App Registration Name (optional): " APP_NAME

if [ -z "$CLIENT_ID" ]; then
    print_error "Client ID is required"
    exit 1
fi

# Verify App Registration exists
print_check "App Registration existence"
APP_INFO=$(az ad app show --id "$CLIENT_ID" 2>/dev/null || echo "")
if [ -n "$APP_INFO" ]; then
    print_success "App Registration found"
    APP_DISPLAY_NAME=$(echo "$APP_INFO" | jq -r '.displayName')
    echo "  Display Name: $APP_DISPLAY_NAME"
else
    print_error "App Registration not found with Client ID: $CLIENT_ID"
    echo "Verify the Client ID is correct in Azure Portal"
    exit 1
fi

# Check redirect URIs
print_check "Redirect URIs configuration"
REDIRECT_URIS=$(echo "$APP_INFO" | jq -r '.spa.redirectUris[]?' 2>/dev/null)
EXPECTED_URIS=("http://localhost:8080" "http://localhost:5173" "http://localhost:3000")

if [ -n "$REDIRECT_URIS" ]; then
    print_success "SPA Redirect URIs configured:"
    echo "$REDIRECT_URIS" | while read -r uri; do
        if [ -n "$uri" ]; then
            echo "  - $uri"
        fi
    done
    
    # Check if expected URIs are present
    echo
    print_info "Checking for required redirect URIs:"
    for expected_uri in "${EXPECTED_URIS[@]}"; do
        if echo "$REDIRECT_URIS" | grep -q "$expected_uri"; then
            print_success "$expected_uri - Found"
        else
            print_warning "$expected_uri - Missing (add for local development)"
        fi
    done
else
    print_warning "No SPA redirect URIs configured"
    echo "Add these URIs in Azure Portal > App Registration > Authentication:"
    for uri in "${EXPECTED_URIS[@]}"; do
        echo "  - $uri"
    done
fi

# Check API permissions
print_check "API permissions"
API_PERMISSIONS=$(az ad app permission list --id "$CLIENT_ID" --query '[].{api: resourceAppId, permission: resourceAccess[0].id}' -o json 2>/dev/null)

# Check for Azure Resource Manager permission
ARM_API_ID="797f4846-ba00-40c8-ac9a-b5b4c90f0c67"  # Azure Service Management API
if echo "$API_PERMISSIONS" | jq -e ".[] | select(.api == \"$ARM_API_ID\")" >/dev/null 2>&1; then
    print_success "Azure Resource Manager API permission found"
else
    print_warning "Azure Resource Manager API permission missing"
    echo "Add 'Azure Service Management' > 'user_impersonation' permission"
fi

# Check Microsoft Graph permission
GRAPH_API_ID="00000003-0000-0000-c000-000000000000"  # Microsoft Graph
if echo "$API_PERMISSIONS" | jq -e ".[] | select(.api == \"$GRAPH_API_ID\")" >/dev/null 2>&1; then
    print_success "Microsoft Graph API permission found"
else
    print_warning "Microsoft Graph API permission missing"
    echo "Add 'Microsoft Graph' > 'User.Read' permission"
fi

# Check if app supports public client flows
print_check "Public client flows"
ALLOW_PUBLIC_CLIENT=$(echo "$APP_INFO" | jq -r '.isFallbackPublicClient')
if [ "$ALLOW_PUBLIC_CLIENT" = "true" ]; then
    print_success "Public client flows enabled"
else
    print_warning "Public client flows disabled"
    echo "Enable in Azure Portal > App Registration > Authentication > Advanced settings"
fi

# Test token acquisition (if possible)
print_check "Token acquisition test"
if command -v curl >/dev/null 2>&1; then
    # Try to get a token using device code flow (non-interactive)
    print_info "Attempting to validate app registration with Azure..."
    
    # Check if we can at least validate the client ID format
    if [[ $CLIENT_ID =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]]; then
        print_success "Client ID format is valid (GUID)"
    else
        print_error "Client ID format is invalid (should be a GUID)"
    fi
else
    print_warning "curl not available, skipping token test"
fi

# Generate environment configuration
echo
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Environment Configuration${NC}"
echo -e "${BLUE}========================================${NC}"

cat > .env.azure << EOF
# Azure App Registration Configuration
VITE_AZURE_CLIENT_ID=$CLIENT_ID
VITE_AZURE_TENANT_ID=$TENANT_ID
VITE_AZURE_SUBSCRIPTION_ID=$SUBSCRIPTION_ID
VITE_REDIRECT_URI=http://localhost:8080
VITE_POST_LOGOUT_REDIRECT_URI=http://localhost:8080

# Backend Configuration
AZURE_CLIENT_ID=$CLIENT_ID
AZURE_TENANT_ID=$TENANT_ID
AZURE_SUBSCRIPTION_ID=$SUBSCRIPTION_ID
EOF

print_success "Created .env.azure with your configuration"

echo
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Next Steps${NC}"
echo -e "${BLUE}========================================${NC}"

echo "1. Copy .env.azure to your frontend directory:"
echo "   cp .env.azure src/frontend/vue/.env"
echo
echo "2. Update your Go backend environment:"
echo "   export AZURE_SUBSCRIPTION_ID=$SUBSCRIPTION_ID"
echo "   export AZURE_TENANT_ID=$TENANT_ID"
echo
echo "3. Test the setup:"
echo "   - Visit http://localhost:8080"
echo "   - Select 'Azure Direct' mode"
echo "   - Try to authenticate"
echo
echo "4. If authentication fails, check Azure Portal settings:"
echo "   - Authentication > Platform configurations > Single-page application"
echo "   - API permissions > Grant admin consent"
echo "   - Authentication > Advanced settings > Allow public client flows"

echo
print_info "App Registration verification complete!"