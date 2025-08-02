# Brutus - Azure CAF Bootstrap Tool

Brutus is a Terminal User Interface (TUI) tool written in Go that bootstraps critical Azure resources following Microsoft's Cloud Adoption Framework (CAF). It creates the foundational infrastructure needed for Terraform-based Infrastructure as Code pipelines.

## Features

- 🖥️ **Interactive TUI** - Built with bubbletea for intuitive terminal experience
- 🏗️ **CAF Compliance** - Creates management groups and landing zones per CAF recommendations  
- 🔄 **Multi-Subscription** - Handles multiple Azure subscriptions for different purposes
- 📊 **State Management** - Tracks created resources for documentation and change management
- ⚙️ **Flexible Configuration** - Supports environment variables, CLI arguments, and config files
- 🐳 **Containerized** - Docker-based for development isolation and pipeline integration

## Prerequisites

- Go 1.23 or later
- Azure CLI (for authentication)
- Docker (optional, for containerized usage)

## Getting Started

### Local Development

1. **Clone and build:**
   ```bash
   git clone <repository-url>
   cd brutus
   go mod tidy
   go build -o brutus cmd/brutus/main.go
   ```

2. **Run the TUI:**
   ```bash
   ./brutus
   ```

### Docker Usage

```bash
# Build the image
docker build -t brutus .

# Run the container
docker run -it brutus
```

## Usage

Brutus provides an interactive menu system for:

1. **Azure Environment Setup** - Configure subscriptions and authentication
2. **CAF Management Groups & Landing Zones** - Create CAF-compliant hierarchy
3. **Bootstrap Resources** - Create storage accounts, key vaults, managed identities
4. **View Current State** - Review created resources and configuration
5. **Configuration** - Manage settings and preferences

## Configuration

Brutus supports multiple configuration methods:

- **Environment Variables**: `AZURE_SUBSCRIPTION_ID`, `AZURE_TENANT_ID`, etc.
- **CLI Arguments**: Command-line flags and subcommands  
- **Config Files**: YAML/JSON configuration files
- **Interactive Prompts**: TUI-based configuration

## Cloud Adoption Framework (CAF)

Brutus implements CAF recommendations for:
- Management group hierarchy
- Landing zone structure
- Naming conventions
- Resource organization
- Security baseline

## Authentication

- **Primary**: Azure CLI (`az login`)
- **Automation**: Service Principal support
- **Container**: Managed Identity integration for pipeline scenarios

## Development

### Project Structure
```
brutus/
├── cmd/brutus/           # Main application entry point
├── internal/
│   ├── tui/             # TUI components and models
│   ├── azure/           # Azure provider implementation
│   ├── caf/             # CAF logic and templates
│   ├── config/          # Configuration management
│   └── state/           # State management
├── pkg/                 # Public packages
├── configs/             # Example configuration files
└── docs/               # Documentation
```

## Future Roadmap

- AWS provider support
- GCP provider support  
- DevOps pipeline integration
- Advanced CAF scenarios
- Terraform module generation