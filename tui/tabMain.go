package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/TresSims/scry/facter"
	"github.com/TresSims/scry/facter/facts"
)

// [MainTab] implements [Tab]
type MainTab struct {
	w, h int

	f facter.Cache

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

	var (
		renderedList string
		logList      []facts.SyslogLine
		ok           bool
		netInfo      facts.NetInfo
		cpuInfo      facts.CPUInfo
		memInfo      facts.MemInfo
	)

	netInfo, ok = t.f[facter.NetInfoKey].(facts.NetInfo)
	if !ok {
		netInfo = facts.NetInfo{
			Rx: -1,
			Tx: -1,
		}
	}

	cpuInfo, ok = t.f[facter.CpuInfoKey].(facts.CPUInfo)
	if !ok {
		cpuInfo = facts.CPUInfo{
			Name:  "unknown",
			Cores: -1,
			Clock: -1.0,
		}
	}

	memInfo, ok = t.f[facter.MemInfoKey].(facts.MemInfo)
	if !ok {
		memInfo = facts.MemInfo{
			Free:  "-1",
			Total: "-1",
		}
	}

	table := table.New().
		Row("Host Info", "System Info", "Network Info").
		Row(
			fmt.Sprintf("hostname: %s", t.f[facter.HostnameKey]),
			fmt.Sprintf("%s: %d cores at %f", cpuInfo.Name, cpuInfo.Cores, cpuInfo.Clock),
			fmt.Sprintf("%s", t.f[facter.IpKey]),
		).
		Row(
			fmt.Sprintf("Uptime: %s", t.f[facter.UptimeKey]),
			fmt.Sprintf("RAM Total: %s, Free: %s", memInfo.Total, memInfo.Free),
			fmt.Sprintf("Connectivity: %t", t.f[facter.ConnectivityKey]),
		).
		Row(
			fmt.Sprintf("%s", t.f[facter.SystemKey]),
			fmt.Sprintf("Network Rx: %d / Tx: %d", netInfo.Rx, netInfo.Tx),
			fmt.Sprintf("Scry Connections: %d", t.f[facter.SubscriptionsKey]),
		).
		Width(t.w)

	renderedTable := table.Render()

	logListUnknown := t.f[facter.JournalKey]
	if logList, ok = logListUnknown.([]facts.SyslogLine); !ok {
		renderedList = "Error reading logList"
	}

	if renderedList != "Error reading logList" {
		for _, line := range logList {
			renderedList += line.String() + "\n"
		}
	}

	t.renderViewport(t.w, t.h-lipgloss.Height(renderedTable), renderedList)

	key := "--- " + lipgloss.NewStyle().Foreground(lipgloss.Blue).Render("Journal") + " "
	hr := lipgloss.NewStyle().Width(t.w).Render(key + strings.Repeat("-", max(0, t.w-lipgloss.Width(key))))

	v.SetContent(lipgloss.JoinVertical(
		lipgloss.Top,
		renderedTable,
		hr,
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
		t.viewport.GotoBottom()
		t.ready = true
	} else {
		// Follow the live stream unless the user has scrolled up
		atBottom := t.viewport.AtBottom()

		t.viewport.SetWidth(w)
		t.viewport.SetHeight(h)
		t.viewport.SetContent(content)

		if atBottom {
			t.viewport.GotoBottom()
		}
	}
}
