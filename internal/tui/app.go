package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App represents the main TUI application
type App struct {
	keys         keyMap
	width        int
	height       int
	state        AppState
	selectedItem int
	menuItems    []string
}

// AppState represents the current state of the application
type AppState int

const (
	StateMainMenu AppState = iota
	StateAzureSetup
	StateCAFSetup
	StateResourceCreation
)

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
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
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
			if a.state == StateMainMenu {
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

	return menuStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}

func (a *App) renderAzureSetup() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("Azure Setup - Coming Soon")
}

func (a *App) renderCAFSetup() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("CAF Setup - Coming Soon")
}

func (a *App) renderResourceCreation() string {
	return lipgloss.NewStyle().Margin(1, 0).Render("Resource Creation - Coming Soon")
}