package main

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

const (
	netIfaceCol  = 10
	netRateCol   = 12
	netTotalsCol = 28

	netMinBarWidth = 8
)

var (
	netRxStyle      = lipgloss.NewStyle().Foreground(lipgloss.Blue)
	netTxStyle      = lipgloss.NewStyle().Foreground(lipgloss.Magenta)
	netSummaryStyle = lipgloss.NewStyle().Faint(true)
)

// [NetInfoTab] implements [tui.Tab] and renders per-interface throughput
// with bars scaled to the busiest interface in the current sample
type NetInfoTab struct {
	w, h int

	f facts.Cache

	bar progress.Model
}

func NewNetInfoTab() *NetInfoTab {
	return &NetInfoTab{
		bar: progress.New(progress.WithDefaultBlend()),
	}
}

func (t *NetInfoTab) Init() tea.Cmd {
	return nil
}

func (t *NetInfoTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.FactsMsg:
		t.f = msg.Facts
	}

	return t, nil
}

func (t *NetInfoTab) View() tea.View {
	var v tea.View

	devs, ok := t.f["network"].([]NetDev)
	if !ok {
		v.SetContent(lipgloss.NewStyle().Width(t.w).Height(t.h).
			Render("waiting for network data..."))

		return v
	}

	rows := []string{}

	if summary := renderTcpSummary(t.f); summary != "" {
		rows = append(rows, summary, "")
	}

	barWidth := max(netMinBarWidth, netBarWidth(t.w))
	t.bar.SetWidth(max(0, barWidth))

	var peak float64
	for _, d := range devs {
		if d.Interface == "lo" {
			continue
		}

		peak = max(peak, d.RxRate, d.TxRate)
	}
	peak = max(peak, 1)

	for _, d := range devs {
		if d.Interface == "lo" {
			continue
		}

		rows = append(rows, t.renderNetRow(d, peak))
	}

	v.SetContent(lipgloss.NewStyle().
		Width(t.w).
		Height(t.h).
		Render(lipgloss.JoinVertical(lipgloss.Top, rows...)))

	return v
}

func (t *NetInfoTab) Name() string {
	return "Network"
}

func (t *NetInfoTab) SetSize(w, h int) {
	t.w = w
	t.h = h

	t.bar.SetWidth(max(0, netBarWidth(w)))
}

// netBarWidth returns the width left for a rate bar after the label and
// value columns
func netBarWidth(w int) int {
	fixed := netIfaceCol + 2*(netRateCol+4) + netTotalsCol + 6

	return w - fixed
}

func (t *NetInfoTab) renderNetRow(d NetDev, peak float64) string {
	iface := lipgloss.NewStyle().Width(netIfaceCol).Render(d.Interface)

	rx := lipgloss.JoinHorizontal(
		lipgloss.Left,
		netRxStyle.Render("↓"),
		t.bar.ViewAs(d.RxRate/peak),
		lipgloss.NewStyle().Width(netRateCol).Render(humanizeRate(d.RxRate)),
	)

	tx := lipgloss.JoinHorizontal(
		lipgloss.Left,
		netTxStyle.Render("↑"),
		t.bar.ViewAs(d.TxRate/peak),
		lipgloss.NewStyle().Width(netRateCol).Render(humanizeRate(d.TxRate)),
	)

	totals := netSummaryStyle.Width(netTotalsCol).
		Render(fmt.Sprintf("rx %s tx %s", humanizeBytes(float64(d.RxBytes)), humanizeBytes(float64(d.TxBytes))))

	return lipgloss.JoinHorizontal(lipgloss.Top, iface, rx, "  ", tx, "  ", totals)
}

func renderTcpSummary(f facts.Cache) string {
	states, ok := f["tcp_states"].(TcpStates)
	if !ok || len(states) == 0 {
		return ""
	}

	names := make([]string, 0, len(states))
	for name, count := range states {
		if count > 0 {
			names = append(names, name)
		}
	}
	slices.Sort(names)

	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, fmt.Sprintf("%d %s", states[name], name))
	}

	return netSummaryStyle.Render("TCP " + strings.Join(parts, " · "))
}

func humanizeBytes(bytes float64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%.0f B", bytes)
	}

	names := []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	val := bytes / unit
	for _, name := range names {
		if val < unit {
			return fmt.Sprintf("%.1f %s", val, name)
		}

		val /= unit
	}

	return fmt.Sprintf("%.1f %s", val, names[len(names)-1])
}

func humanizeRate(bps float64) string {
	return humanizeBytes(bps) + "/s"
}
