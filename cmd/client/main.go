package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/TresSims/scry/tui"
)

func main() {
	p := tea.NewProgram(tui.New(
		tui.AsClient,
		tui.MainTab{}.Load,
		tui.MainTab{}.Load,
	))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Critical error starting scry: %v", err)
		os.Exit(1)
	}
}
