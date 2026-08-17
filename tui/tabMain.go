package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/TresSims/scry/facts"
)

// [MainTab] implements [Tab]
type MainTab struct {
	w, h int

	f facts.Cache
}

func (_ *MainTab) Init() tea.Cmd {
	return nil
}

func (t *MainTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FactsMsg:
		t.f = msg.Facts
	}
	return t, nil
}

func (t *MainTab) View() tea.View {
	var v tea.View

	table := table.New().
		Row(fmt.Sprintf("hostname: %s", t.f["hostname"]), fmt.Sprintf("connectivity: %t", t.f["connectivity"])).
		Width(t.w)

	renderedTable := table.Render()

	v.SetContent(lipgloss.JoinVertical(lipgloss.Top, renderedTable))

	return v
}

func (t *MainTab) Name() string {
	return "Main"
}

func (t *MainTab) SetSize(w, h int) {
	t.w = w
	t.h = h
}
