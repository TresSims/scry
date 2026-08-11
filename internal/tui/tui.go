package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	Checks   map[string]bool
	hostname string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error

	m.hostname, err = os.Hostname()
	if err != nil {
		panic("WTF?")
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	s := fmt.Sprintf("Bubble Tea View from %s\n\n", m.hostname)

	for check, status := range m.Checks {
		s += fmt.Sprintf("%s: ", check)

		if status {
			s += "Passed."
		} else {
			s += "Failed!"
		}

		s += "\n\n"
	}

	s += "\n"

	return tea.NewView(s)
}
