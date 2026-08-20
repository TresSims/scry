package tui

import (
	"fmt"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/TresSims/scry/facts"
)

// [MainTab] implements [Tab]
type MainTab struct {
	w, h int

	f facts.Cache

	viewport viewport.Model
	ready    bool
}

func (m *MainTab) Init() tea.Cmd {
	return nil
}

func (t *MainTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case FactsMsg:
		t.f = msg.Facts
	}

	t.viewport, cmd = t.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return t, tea.Batch(cmds...)
}

func (t *MainTab) View() tea.View {
	var v tea.View

	table := table.New().
		Row(fmt.Sprintf("hostname: %s", t.f["hostname"]), fmt.Sprintf("connectivity: %t", t.f["connectivity"])).
		Width(t.w)

	renderedTable := table.Render()

	var (
		renderedList string
		logList      []facts.SyslogLine
		ok           bool
	)

	logListUnknown := t.f["journal"]
	if logList, ok = logListUnknown.([]facts.SyslogLine); !ok {
		renderedList = "Error reading logList"
	}

	if renderedList != "Error reading logList" {
		for _, line := range logList {
			renderedList += line.String() + "\n"
		}
	}

	t.renderViewport(t.w, t.h-lipgloss.Height(renderedTable), renderedList)

	v.SetContent(lipgloss.JoinVertical(
		lipgloss.Top,
		renderedTable,
		t.viewport.View(),
	))

	return v
}

func (t *MainTab) Name() string {
	return "Main"
}

func (t *MainTab) SetSize(w, h int) {
	t.w = w
	t.h = h
}

func (t *MainTab) renderViewport(w, h int, content string) {
	if !t.ready {
		t.viewport = viewport.New(viewport.WithHeight(h), viewport.WithWidth(w))
		// tab height - viewport height = Y Offset
		t.viewport.YPosition = t.h - h
		t.viewport.LeftGutterFunc = func(info viewport.GutterContext) string {
			switch {

			case info.Soft:
				return "   | "
			case info.Index >= info.TotalLines:
				return "  ~| "
			default:
				return fmt.Sprintf("%4d | ", info.Index+1)
			}
		}
		t.viewport.SetContent(content)
		t.ready = true
	} else {
		t.viewport.SetWidth(w)
		t.viewport.SetHeight(h)
		t.viewport.SetContent(content)
	}
}
