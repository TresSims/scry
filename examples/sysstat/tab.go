package main

import (
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

// metric pairs a fact key with the label shown next to its progress bar
type metric struct {
	key   string
	label string

	// optional metrics are hidden until their facter publishes a value
	optional bool
}

// metrics are rendered in declaration order
var metrics = []metric{
	{key: "cpu", label: "CPU"},
	{key: "memory", label: "Memory"},
	{key: "gpu", label: "GPU", optional: true},
	{key: "battery", label: "Battery", optional: true},
}

const labelWidth = 8

// SysStatTab implements [tui.Tab]
type SysStatTab struct {
	w, h int

	f facts.Cache

	bar progress.Model
}

func NewSysStatTab() *SysStatTab {
	return &SysStatTab{
		bar: progress.New(progress.WithDefaultBlend()),
	}
}

func (t *SysStatTab) Init() tea.Cmd {
	return nil
}

func (t *SysStatTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.FactsMsg:
		t.f = msg.Facts
	}

	return t, nil
}

func (t *SysStatTab) View() tea.View {
	var v tea.View

	labelStyle := lipgloss.NewStyle().Width(labelWidth)

	rows := []string{}
	for _, m := range metrics {
		val, ok := t.f[m.key].(float64)
		if !ok && m.optional {
			continue
		}

		rows = append(rows, lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render(m.label),
			t.bar.ViewAs(val),
		))
	}

	v.SetContent(lipgloss.NewStyle().
		Width(t.w).
		Height(t.h).
		Render(lipgloss.JoinVertical(lipgloss.Top, rows...)))

	return v
}

func (t *SysStatTab) Name() string {
	return "SysStat"
}

func (t *SysStatTab) SetSize(w, h int) {
	t.w = w
	t.h = h

	t.bar.SetWidth(max(0, w-labelWidth))
}
