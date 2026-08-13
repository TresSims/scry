package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/TresSims/scry/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.New(
		tui.AsClient,
		tui.WithMainTab,
		tui.WithMainTab,
	))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Critical error starting scry: %v", err)
		os.Exit(1)
	}
}
