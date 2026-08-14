package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := &facts.Engine{
		Data: map[string]facts.Fact{
			"hostname": {
				Facter: func() (any, error) { return os.Hostname() },
			},
			"count": {
				Facter: func() (any, error) { return rand.Int(), nil },
			},
		},
	}
	go e.Collect(ctx)

	p := tea.NewProgram(tui.New(
		tui.AsClient,
		tui.WithEngine(e),
		(&tui.MainTab{}).Load,
		(&tui.MainTab{}).Load,
	))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Critical error starting scry: %v", err)
		os.Exit(1)
	}
}
