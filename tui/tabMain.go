package tui

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
)

// [MainTab] implements [Tab]
type MainTab struct{}

func (t *MainTab) Render(w, h int, f facts.Cache, style lipgloss.Style) string {
	render := style.Width(w).Height(h).Render(fmt.Sprintf("hostname: %s \n\n\nRandom: %d", f["hostname"], f["count"]))

	return render
}

func (t *MainTab) Name() string {
	return "Main"
}
