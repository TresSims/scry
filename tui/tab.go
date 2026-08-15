package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
)

// Tab is an interface for rendering a tab page for the scry tui
type Tab interface {
	// Render takes (w)idth and (h)eight, the most recent facts, and the default
	// style. returns a string filling those dimensions, and optionally using
	// that style
	Render(w, h int, f facts.Cache, style lipgloss.Style) string

	// Name returns the name of the tab
	Name() string
}
