package main

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facter"
	"github.com/TresSims/scry/tui"
)

const (
	TimeKey     = "time"
	TimezoneKey = "timezone"
)

type PluginSample struct{}

func CurrentTime(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			now := time.Now()

			publish(now)
		case <-ctx.Done():
			return nil
		}
	}
}

func Timezone(ctx context.Context, publish func(val any)) error {
	publish(time.Now().Location())

	return nil
}

func (p *PluginSample) Facters() map[string]facter.Facter {
	facters := map[string]facter.Facter{}

	facters[TimeKey] = CurrentTime
	facters[TimezoneKey] = Timezone

	return facters
}

type SampleTab struct {
	w int
	h int

	f facter.Cache
}

func (t *SampleTab) Init() tea.Cmd {
	return nil
}

func (t *SampleTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.FactsMsg:
		t.f = msg.Facts
	}

	return t, nil
}

func (t *SampleTab) View() tea.View {
	var v tea.View

	style := lipgloss.NewStyle().
		Width(t.w).
		Height(t.h).
		Border(lipgloss.RoundedBorder()).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render(fmt.Sprintf("%s\n%s", t.f[TimeKey], t.f[TimezoneKey]))

	v.SetContent(style)

	return v
}

func (t *SampleTab) Name() string {
	return "Time"
}

func (t *SampleTab) SetSize(w, h int) {
	t.w = w
	t.h = h
}

func (p *PluginSample) Tabs() []tui.Tab {
	tabs := []tui.Tab{&SampleTab{}}

	return tabs
}

var Bundle = PluginSample{}
