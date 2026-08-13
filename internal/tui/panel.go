package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Mode int

const (
	Client Mode = iota
	Server
)

// [Model] implements [tea.Model]
type Model struct {
	// viewport [w]idth, [h]eight
	w, h int

	// index of the current tab
	t int

	// client or server mode, effectively whether to fetch data from the
	// local machine or a remote one
	mode Mode

	// tabs has a list of tabs that can be rendered
	tabs []Tab

	// the default stlye
	style lipgloss.Style
}

type Option func(*Model)

func New(opts ...Option) *Model {
	model := &Model{}

	for _, opt := range opts {
		opt(model)
	}

	model.t = 0

	model.style = lipgloss.NewStyle().
		Foreground(lipgloss.White).
		Background(lipgloss.Black)

	return model
}

func AsClient(m *Model) {
	m.mode = Client
}

func AsServer(m *Model) {
	m.mode = Server
}

func (_ Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j":
			if m.t < len(m.tabs)-1 {
				m.t += 1
			}
			return m, nil
		case "k":
			if m.t > 0 {
				m.t -= 1
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
	}

	return m, nil
}

func (m Model) View() tea.View {
	header := m.renderHeader()
	headerHeight := lipgloss.Height(header)

	tab := m.tabs[m.t].Render(m.w, m.h-headerHeight, m.style)
	tabHeight := lipgloss.Height(tab)

	// the tab should fill the hight, but we pad it out here just in case
	spacer := m.style.Height(m.h - headerHeight - tabHeight).Render()

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Top, header, tab, spacer))
}

func (m Model) renderHeader() string {
	tabStyle := m.style.
		Padding(2).
		Border(lipgloss.RoundedBorder())

	primaryTabStyle := tabStyle.BorderBottom(false).
		Foreground(lipgloss.Blue)

	renderedTabs := []string{}

	for i, tab := range m.tabs {
		if i == m.t {
			renderedTabs = append(renderedTabs, primaryTabStyle.Render(tab.Name()))

			continue
		}

		renderedTabs = append(renderedTabs, tabStyle.Render(tab.Name()))
	}

	return lipgloss.JoinHorizontal(0.2, renderedTabs...)
}
