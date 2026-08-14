package tui

import "charm.land/lipgloss/v2"

// [MainTab] implements [Tab]
type MainTab struct {
}

func WithMainTab(m *Model) {
	m.tabs = append(m.tabs, MainTab{})
}

func (t MainTab) Render(w, h int, style lipgloss.Style) string {
	render := style.Width(w).Height(h).Render()

	return render
}

func (t MainTab) Name() string {
	return "Main"
}
