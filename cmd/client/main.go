package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

func main() {
	p := tea.NewProgram(tui.New(
		tui.AsClient,
		tui.MainTab{}.Load,
		tui.MainTab{}.Load,
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := facts.Engine{
		Data: map[string]facts.Fact{
			"hostname": {
				Facter: func() (any, error) { return os.Hostname() },
			},
		},
	}
	go e.Collect(ctx)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Critical error starting scry: %v", err)
		os.Exit(1)
	}
}
