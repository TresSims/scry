package tui

import tea "charm.land/bubbletea/v2"

// Tab is an interface for rendering a tab page for the scry tui
type Tab interface {
	// Init implements [tea.Model]
	Init() tea.Cmd

	// Update implements [tea.Model]
	Update(tea.Msg) (tea.Model, tea.Cmd)

	// View implements [tea.View]
	View() tea.View

	// Name returns the name of the tab
	Name() string

	// SetSize passes thet usable bubble tea size to it. The [Tab] is expected to
	// respect this
	SetSize(w, h int)
}
