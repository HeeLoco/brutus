package tui

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edgard-ott/brutus/internal/azure"
)

// App represents the main TUI application
type App struct {
	keys                 keyMap
	width                int
	height               int
	state                AppState
	selectedItem         int
	menuItems            []string
	azureAuth            *azure.AuthInfo
	subscriptions        []*armsubscriptions.Subscription
	authStatus           string
	isAuthenticating     bool
	managementGroupTree  *azure.ManagementGroupNode
	mgMenuItems          []string
	isLoadingMGs         bool
	mgActionFeedback     string
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

type mgLoadSuccessMsg struct {
	tree *azure.ManagementGroupNode
}

type mgLoadErrorMsg struct {
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

// loadMGsCmd loads management groups in the background
func loadMGsCmd(auth *azure.AuthInfo) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		
		// Try to load real management groups first
		if err := auth.ListManagementGroups(ctx); err != nil {
			// If failed, still return mock structure
		}
		
		tree, err := auth.GetManagementGroupTree(ctx)
		if err != nil {
			return mgLoadErrorMsg{err: err}
		}
		
		return mgLoadSuccessMsg{tree: tree}
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

	mgMenuItems := []string{
		"📝 Edit Management Group Structure",
		"✏️  Rename Specific Management Group",
		"🏷️  Rename All Management Groups",
		"🏗️  Apply Basic CAF Management Structure",
		"🔄 Refresh Management Groups",
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
		mgMenuItems:  mgMenuItems,
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

	case mgLoadSuccessMsg:
		a.managementGroupTree = msg.tree
		a.isLoadingMGs = false
		return a, nil

	case mgLoadErrorMsg:
		log.Printf("Failed to load management groups: %v", msg.err)
		a.isLoadingMGs = false
		return a, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keys.Back):
			if a.state != StateMainMenu {
				a.state = StateMainMenu
				a.selectedItem = 0
				a.mgActionFeedback = "" // Clear any feedback when going back
			}
			return a, nil
		case key.Matches(msg, a.keys.Up):
			switch a.state {
			case StateMainMenu:
				if a.selectedItem > 0 {
					a.selectedItem--
				}
			case StateCAFSetup:
				if a.selectedItem > 0 {
					a.selectedItem--
				}
			}
			return a, nil
		case key.Matches(msg, a.keys.Down):
			switch a.state {
			case StateMainMenu:
				if a.selectedItem < len(a.menuItems)-1 {
					a.selectedItem++
				}
			case StateCAFSetup:
				if a.selectedItem < len(a.mgMenuItems)-1 {
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
					a.selectedItem = 0 // Reset selection for CAF menu
					if a.azureAuth != nil && a.managementGroupTree == nil {
						a.isLoadingMGs = true
						return a, loadMGsCmd(a.azureAuth)
					}
				case 2:
					a.state = StateResourceCreation
				case 3:
					// View Current State - stay in main menu for now
				case 4:
					// Configuration - stay in main menu for now
				}
			} else if a.state == StateCAFSetup {
				// Handle management group actions (mock for now)
				switch a.selectedItem {
				case 0:
					// Edit Management Group Structure - mock
					a.mgActionFeedback = "📝 Mock: Opening management group structure editor..."
				case 1:
					// Rename Specific Management Group - mock
					a.mgActionFeedback = "✏️ Mock: Select a specific management group to rename..."
				case 2:
					// Rename All Management Groups - mock
					a.mgActionFeedback = "🏷️ Mock: Renaming all management groups with CAF conventions..."
				case 3:
					// Apply Basic CAF Management Structure - mock
					a.mgActionFeedback = "🏗️ Mock: Applying basic CAF management structure to tenant..."
				case 4:
					// Refresh Management Groups
					a.mgActionFeedback = ""
					if a.azureAuth != nil {
						a.isLoadingMGs = true
						return a, loadMGsCmd(a.azureAuth)
					}
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
	// Permission Information Section
	permissionSection := a.renderMGPermissionInfo()
	
	// Structure Display Section
	structureSection := a.renderManagementGroupStructure()
	
	// Action Menu Section
	actionSection := a.renderCAFActionMenu()
	
	return lipgloss.JoinVertical(lipgloss.Left, permissionSection, structureSection, actionSection)
}

func (a *App) renderMGPermissionInfo() string {
	if a.azureAuth == nil || a.isLoadingMGs {
		return "" // Don't show permission info while loading
	}

	// Determine border color based on permission status
	borderColor := "#DC2626" // Red for no access
	if a.azureAuth.HasMGReadAccess {
		borderColor = "#059669" // Green for access
	} else if a.azureAuth.MGPermissions == "" {
		return "" // Don't show if permissions haven't been checked yet
	}

	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 2).
		MarginBottom(1)

	var content []string
	
	if a.azureAuth.HasMGReadAccess {
		// User has access - show brief status
		content = append(content, "ℹ️  Management Group Permissions")
		content = append(content, "")
		content = append(content, a.azureAuth.MGPermissions)
	} else {
		// User lacks access - show detailed guidance
		content = append(content, "⚠️  Management Group Access Required")
		content = append(content, "")
		content = append(content, a.azureAuth.MGPermissions)
		content = append(content, "")
		content = append(content, "📋 Required Roles:")
		content = append(content, "   • Management Group Reader (read access)")
		content = append(content, "   • Management Group Contributor (full access)")
		content = append(content, "")
		content = append(content, "🔧 How to get access:")
		content = append(content, "   1. Ask your Azure AD Global Administrator")
		content = append(content, "   2. Request assignment at Tenant Root Group level")
		content = append(content, "   3. Or enable 'Access management for Azure resources'")
		content = append(content, "      in Azure AD Properties (if you're Global Admin)")
		content = append(content, "")
		if !a.azureAuth.HasMGReadAccess {
			content = append(content, "💡 Currently showing mock CAF structure for demonstration")
		}
	}

	return infoStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (a *App) renderManagementGroupStructure() string {
	structureStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981")).
		Padding(1, 2).
		MarginBottom(2)

	var content []string
	content = append(content, "🏗️ CAF Management Group Structure")
	content = append(content, "")

	if a.isLoadingMGs {
		content = append(content, "🔄 Loading management groups...")
	} else if a.managementGroupTree != nil {
		// Add note if showing mock structure
		if a.azureAuth != nil && !a.azureAuth.HasMGReadAccess {
			mockStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F59E0B")).
				Italic(true)
			content = append(content, mockStyle.Render("📋 Mock CAF Structure (Demo)"))
			content = append(content, "")
		}
		
		treeLines := a.renderMGTree(a.managementGroupTree, "")
		content = append(content, treeLines...)
	} else {
		content = append(content, "❌ No management group structure available")
		content = append(content, "")
		content = append(content, "Use the 'Refresh Management Groups' option to load the structure.")
	}

	return structureStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (a *App) renderMGTree(node *azure.ManagementGroupNode, prefix string) []string {
	var lines []string
	
	// Determine icon and styling based on level
	var icon string
	switch node.Level {
	case 0:
		icon = "🏛️" // Root
	case 1:
		icon = "🏢" // Top level (Platform, Landing Zones, etc.)
	case 2:
		icon = "📁" // Sub-groups
	default:
		icon = "📂"
	}
	
	// Add the current node
	lines = append(lines, fmt.Sprintf("%s%s %s (%s)", prefix, icon, node.DisplayName, node.ID))
	
	// Add subscriptions if any
	for _, sub := range node.Subscriptions {
		subName := "Unknown Subscription"
		if sub.DisplayName != nil {
			subName = *sub.DisplayName
		}
		subID := ""
		if sub.SubscriptionID != nil {
			subID = (*sub.SubscriptionID)[:8] + "..."
		}
		lines = append(lines, fmt.Sprintf("%s  💳 %s (%s)", prefix, subName, subID))
	}
	
	// Add children recursively
	for i, child := range node.Children {
		childPrefix := prefix
		if i == len(node.Children)-1 {
			childPrefix += "└── "
		} else {
			childPrefix += "├── "
		}
		
		childLines := a.renderMGTree(child, prefix+"    ")
		lines = append(lines, childLines...)
	}
	
	return lines
}

func (a *App) renderCAFActionMenu() string {
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(1, 2)

	var content []string
	content = append(content, "⚙️ Management Group Actions")
	content = append(content, "")

	itemStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginBottom(1).
		Background(lipgloss.Color("#374151")).
		Foreground(lipgloss.Color("#F3F4F6"))

	for i, item := range a.mgMenuItems {
		if i == a.selectedItem {
			content = append(content, selectedStyle.Render("► "+item))
		} else {
			content = append(content, itemStyle.Render("  "+item))
		}
	}

	// Add feedback if any action was performed
	if a.mgActionFeedback != "" {
		content = append(content, "")
		feedbackStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Italic(true)
		content = append(content, feedbackStyle.Render(a.mgActionFeedback))
	}

	return menuStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (a *App) renderResourceCreation() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("Resource Creation - Coming Soon")
}