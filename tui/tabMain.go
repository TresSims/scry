package tui

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
)

// [MainTab] implements [Tab]
type MainTab struct {
	Engine *facts.Engine
}

func (t *MainTab) Load(m *Model) {
	t.Engine = m.Engine
	m.tabs = append(m.tabs, t)
}

func (t *MainTab) Render(w, h int, style lipgloss.Style) string {
	render := style.Width(w).Height(h).Render(fmt.Sprintf("hostname: %s \n\n\nRandom: %d", t.Engine.Data["hostname"].Cache, t.Engine.Data["count"].Cache))

	return render
}

func (t *MainTab) Name() string {
	return "Main"
}
