# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Brutus is a Terminal User Interface (TUI) tool written in Go that bootstraps critical Azure resources following Microsoft's Cloud Adoption Framework (CAF). It creates the foundational infrastructure needed for Terraform-based IaC pipelines, with future extensibility for DevOps scenarios.

## Key Features

- **TUI Interface**: Built with bubbletea for interactive terminal experience
- **CAF Compliance**: Creates management groups and landing zones per CAF recommendations
- **Multi-Subscription**: Handles multiple Azure subscriptions for different purposes
- **State Management**: Tracks created resources for documentation and change management
- **Configuration**: Supports environment variables, CLI arguments, and config files (YAML/JSON)
- **Containerized**: Docker-based for development isolation and pipeline integration

## Development Environment

- **Language**: Go with bubbletea TUI framework
- **Container**: Ubuntu-based Docker image for development and deployment
- **Authentication**: Azure CLI integration for authentication
- **Development**: Docker containers for isolated development environment

## Architecture

### Core Components
- **TUI Layer**: Interactive terminal interface with bubbletea
- **Azure Provider**: Azure SDK integration for resource management
- **CAF Engine**: Management group and landing zone creation logic
- **State Manager**: Resource tracking and state persistence
- **Config System**: Multi-format configuration support

### Resource Creation Scope
**Phase 1 (Current)**:
- Management Groups (CAF hierarchy)
- Resource Groups
- Storage Accounts (for Terraform state)
- Managed Identities
- Key Vaults

**Future Phases**:
- DevOps-focused extensions
- Additional cloud providers (AWS, GCP)

## Common Commands

```bash
# Initialize Go module
go mod init brutus
go mod tidy

# Build the project
go build -o brutus cmd/brutus/main.go

# Run tests
go test ./...
go test -v ./...

# Run the TUI tool
./brutus

# Build Docker image
docker build -t brutus .

# Run in container
docker run -it brutus
```

## Project Structure

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
├── docs/               # Documentation
├── Dockerfile          # Container definition
└── go.mod              # Go module definition
```

## Configuration

Supports multiple configuration methods:
1. Environment variables (AZURE_SUBSCRIPTION_ID, etc.)
2. CLI arguments and subcommands
3. YAML/JSON configuration files
4. Interactive TUI prompts

## Authentication

- Primary: Azure CLI (`az login`)
- Automation: Service Principal support
- Container: Managed Identity integration for pipeline scenarios

## State Management

Brutus maintains state for:
- Created resource inventory
- CAF structure mapping
- Configuration snapshots
- Change tracking for updates