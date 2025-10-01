package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea/v2"

	"dotfiles-manager/core"
)

// Model represents the TUI state
type Model struct {
	manager     *core.DotfilesManager
	choices     []string
	cursor      int
	selected    map[int]struct{}
	currentView string
	files       []string
	ready       bool
}

// NewModel creates a new TUI model
func NewModel(manager *core.DotfilesManager) Model {
	return Model{
		manager:  manager,
		choices:  []string{"Status", "Apply", "Add File", "View Diff", "Update", "Quit"},
		selected: make(map[int]struct{}),
		currentView: "menu", // menu, status, apply, etc.
		files:    []string{},
		ready:    false,
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		tea.EnterAltScreen,
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyInput(msg)
	case tea.WindowSizeMsg:
		m.ready = true
		return m, nil
	}

	return m, nil
}

// handleKeyInput handles keyboard input
func (m Model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter", " ":
		// Handle menu selection
		switch m.cursor {
		case 0: // Status
			m.currentView = "status"
			status, err := m.manager.Status()
			if err != nil {
				// For now, we'll just return to menu on error
				m.currentView = "menu"
			} else {
				// In a real implementation, we'd parse the status and show it
				// For now, we'll just go back to menu
				m.currentView = "menu"
			}
		case 1: // Apply
			m.currentView = "applying"
			err := m.manager.Apply()
			if err != nil {
				// Handle error
			}
			m.currentView = "menu"
		case 2: // Add File
			m.currentView = "add"
			// Would show add file view in real implementation
			m.currentView = "menu"
		case 3: // View Diff
			m.currentView = "diff"
			diff, err := m.manager.Diff()
			if err != nil {
				// Handle error
			}
			// Would show diff in real implementation
			m.currentView = "menu"
		case 4: // Update
			m.currentView = "updating"
			err := m.manager.Update()
			if err != nil {
				// Handle error
			}
			m.currentView = "menu"
		case 5: // Quit
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing...\n"
	}

	switch m.currentView {
	case "menu":
		return m.renderMenu()
	default:
		return m.renderMenu()
	}
}

// renderMenu renders the main menu
func (m Model) renderMenu() string {
	s := "\n  Dotfiles Manager - Manage your configuration files\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\n  Press q to quit\n"

	return s
}

// Run starts the TUI
func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}