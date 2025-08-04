package tui

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edgard-ott/brutus/internal/azure"
)

// App represents the main TUI application
type App struct {
	keys            keyMap
	width           int
	height          int
	state           AppState
	selectedItem    int
	menuItems       []string
	azureAuth       *azure.AuthInfo
	subscriptions   []*armsubscriptions.Subscription
	authStatus      string
	isAuthenticating bool
}

// AppState represents the current state of the application
type AppState int

const (
	StateMainMenu AppState = iota
	StateAzureSetup
	StateCAFSetup
	StateResourceCreation
)

// Messages for async operations
type authSuccessMsg struct {
	auth          *azure.AuthInfo
	subscriptions []*armsubscriptions.Subscription
}

type authErrorMsg struct {
	err error
}

// authenticateCmd performs Azure authentication in the background
func authenticateCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		
		auth, err := azure.NewAuth(ctx)
		if err != nil {
			return authErrorMsg{err: err}
		}
		
		if err := auth.ValidateAuthentication(ctx); err != nil {
			return authErrorMsg{err: err}
		}
		
		subscriptions, err := auth.ListSubscriptions(ctx)
		if err != nil {
			return authErrorMsg{err: err}
		}
		
		return authSuccessMsg{
			auth:          auth,
			subscriptions: subscriptions,
		}
	}
}

// keyMap defines key bindings
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
}

// ShortHelp returns key bindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Quit}
}

// FullHelp returns extended key bindings
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Back, k.Quit},
	}
}

// NewApp creates a new TUI application
func NewApp() *App {
	menuItems := []string{
		"🔧 Azure Environment Setup",
		"🏗️  CAF Management Groups & Landing Zones",
		"📦 Create Bootstrap Resources",
		"📊 View Current State",
		"⚙️  Configuration",
	}

	return &App{
		keys: keyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			Back: key.NewBinding(
				key.WithKeys("esc", "backspace"),
				key.WithHelp("esc", "back"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		state:        StateMainMenu,
		selectedItem: 0,
		menuItems:    menuItems,
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	a.isAuthenticating = true
	a.authStatus = "🔄 Initializing Azure authentication..."
	return authenticateCmd()
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case authSuccessMsg:
		a.azureAuth = msg.auth
		a.subscriptions = msg.subscriptions
		if msg.auth.DefaultSub != nil && msg.auth.DefaultSub.DisplayName != nil {
			a.authStatus = fmt.Sprintf("✅ Connected to Azure - %s", *msg.auth.DefaultSub.DisplayName)
		} else {
			a.authStatus = "✅ Connected to Azure"
		}
		a.isAuthenticating = false
		return a, nil

	case authErrorMsg:
		a.authStatus = fmt.Sprintf("❌ Authentication failed: %s", msg.err.Error())
		a.isAuthenticating = false
		return a, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keys.Back):
			if a.state != StateMainMenu {
				a.state = StateMainMenu
				a.selectedItem = 0
			}
			return a, nil
		case key.Matches(msg, a.keys.Up):
			if a.state == StateMainMenu {
				if a.selectedItem > 0 {
					a.selectedItem--
				}
			}
			return a, nil
		case key.Matches(msg, a.keys.Down):
			if a.state == StateMainMenu {
				if a.selectedItem < len(a.menuItems)-1 {
					a.selectedItem++
				}
			}
			return a, nil
		case key.Matches(msg, a.keys.Enter):
			if a.state == StateMainMenu && !a.isAuthenticating {
				switch a.selectedItem {
				case 0:
					a.state = StateAzureSetup
				case 1:
					a.state = StateCAFSetup
				case 2:
					a.state = StateResourceCreation
				case 3:
					// View Current State - stay in main menu for now
				case 4:
					// Configuration - stay in main menu for now
				}
			}
			return a, nil
		}
	}

	return a, nil
}

// View implements tea.Model
func (a *App) View() string {
	var content string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1F2937")).
		Padding(0, 1).
		MarginBottom(1)

	header := headerStyle.Render("🏛️  BRUTUS - Azure CAF Bootstrap Tool")

	// Content based on current state
	switch a.state {
	case StateMainMenu:
		content = a.renderMainMenu()
	case StateAzureSetup:
		content = a.renderAzureSetup()
	case StateCAFSetup:
		content = a.renderCAFSetup()
	case StateResourceCreation:
		content = a.renderResourceCreation()
	}

	// Footer with help
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginTop(2)

	footer := footerStyle.Render("Use ↑/↓ to navigate • Enter to select • ESC to go back • q to quit")

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (a *App) renderMainMenu() string {
	// Azure Environment Status Section
	envSection := a.renderEnvironmentStatus()
	
	// Menu Section
	menuStyle := lipgloss.NewStyle().
		Margin(1, 0)

	itemStyle := lipgloss.NewStyle().
		Padding(0, 4).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Padding(0, 4).
		MarginBottom(1).
		Background(lipgloss.Color("#374151")).
		Foreground(lipgloss.Color("#F3F4F6"))

	var items []string
	for i, item := range a.menuItems {
		if i == a.selectedItem {
			items = append(items, selectedStyle.Render("► "+item))
		} else {
			items = append(items, itemStyle.Render("  "+item))
		}
	}

	menuSection := menuStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)

	return lipgloss.JoinVertical(lipgloss.Left, envSection, menuSection)
}

func (a *App) renderEnvironmentStatus() string {
	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1, 2).
		MarginBottom(2)

	if a.isAuthenticating {
		return statusStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				"🔄 Azure Environment Status",
				"",
				a.authStatus,
			),
		)
	}

	if a.azureAuth == nil {
		return statusStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				"❌ Azure Environment Status",
				"",
				a.authStatus,
				"Please ensure you're authenticated with 'az login'",
			),
		)
	}

	// Build comprehensive status information
	var content []string
	content = append(content, "🌐 Azure Environment Status")
	content = append(content, "")
	
	// Authentication info
	content = append(content, fmt.Sprintf("👤 User: %s", a.azureAuth.CurrentUser))
	content = append(content, fmt.Sprintf("🔐 Auth Method: %s", a.azureAuth.AuthMethod))
	content = append(content, fmt.Sprintf("🏢 Tenant ID: %s", a.azureAuth.TenantID))
	
	// Permission info
	if a.azureAuth.PermissionLevel != "" {
		content = append(content, fmt.Sprintf("🔑 Permission Level: %s", a.azureAuth.PermissionLevel))
	}
	content = append(content, "")
	
	// Subscription summary
	content = append(content, "📋 Subscriptions:")
	if len(a.subscriptions) > 0 {
		content = append(content, fmt.Sprintf("   %s", a.azureAuth.GetSubscriptionSummary()))
		
		// Show enabled subscriptions
		enabled := a.azureAuth.GetEnabledSubscriptions()
		for i, sub := range enabled {
			if i >= 3 { // Show max 3 subscriptions
				content = append(content, fmt.Sprintf("   ... and %d more", len(enabled)-3))
				break
			}
			name := "Unknown"
			if sub.DisplayName != nil {
				name = *sub.DisplayName
			}
			id := ""
			if sub.SubscriptionID != nil {
				id = (*sub.SubscriptionID)[:8] + "..."  // Show first 8 chars of ID
			}
			marker := "  "
			if sub == a.azureAuth.DefaultSub {
				marker = "► " // Mark default subscription
			}
			content = append(content, fmt.Sprintf("   %s%s (%s)", marker, name, id))
		}
	} else {
		content = append(content, "   No subscriptions found")
	}
	
	content = append(content, "")
	
	// CAF Readiness with permission awareness
	if a.azureAuth.CAFPermissions != "" {
		content = append(content, fmt.Sprintf("🏗️ CAF Capability: %s", a.azureAuth.CAFPermissions))
	} else {
		cafIcon := "❌"
		cafStatus := "No enabled subscriptions"
		if a.azureAuth.HasCAFReadySubscription() {
			cafIcon = "✅"
			cafStatus = "Ready for CAF deployment"
		}
		content = append(content, fmt.Sprintf("🏗️ CAF Ready: %s %s", cafIcon, cafStatus))
	}
	
	// Location info
	content = append(content, fmt.Sprintf("🌍 Scope: %s", a.azureAuth.GetLocationInfo()))

	return statusStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (a *App) renderAzureSetup() string {
	style := lipgloss.NewStyle().Margin(1, 0)
	
	var content []string
	content = append(content, "🔧 Azure Environment Setup")
	content = append(content, "")
	
	if a.isAuthenticating {
		content = append(content, a.authStatus)
		content = append(content, "")
		content = append(content, "This will use your existing 'az login' session...")
	} else if a.azureAuth != nil {
		content = append(content, a.authStatus)
		content = append(content, "")
		content = append(content, fmt.Sprintf("Tenant ID: %s", a.azureAuth.TenantID))
		content = append(content, "")
		
		if len(a.subscriptions) > 0 {
			content = append(content, "Available Subscriptions:")
			for i, sub := range a.subscriptions {
				if i < 5 { // Show first 5 subscriptions
					name := "Unknown"
					if sub.DisplayName != nil {
						name = *sub.DisplayName
					}
					id := "Unknown"
					if sub.SubscriptionID != nil {
						id = *sub.SubscriptionID
					}
					content = append(content, fmt.Sprintf("  • %s (%s)", name, id))
				}
			}
			if len(a.subscriptions) > 5 {
				content = append(content, fmt.Sprintf("  ... and %d more", len(a.subscriptions)-5))
			}
		}
	} else if a.authStatus != "" {
		content = append(content, a.authStatus)
		content = append(content, "")
		content = append(content, "Please ensure you are logged in with 'az login' and try again.")
	} else {
		content = append(content, "Authentication will be performed using your Azure CLI session.")
		content = append(content, "")
		content = append(content, "Make sure you're logged in with: az login")
	}
	
	return style.Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (a *App) renderCAFSetup() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("CAF Setup - Coming Soon")
}

func (a *App) renderResourceCreation() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("Resource Creation - Coming Soon")
}