package tui

import "charm.land/lipgloss/v2"

// Tab is an interface for rendering a tab page for the scry tui
type Tab interface {
	// Render takes [w]idth and [h]eight and the default style.
	// returns a string filling those dimensions, and optionally using that style
	Render(w, h int, style lipgloss.Style) string

	// Title returns the name of the tab
	Name() string
}
