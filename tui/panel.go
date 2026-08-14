package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/tui/styles"
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
		Foreground(lipgloss.White)

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

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, header, tab))

	v.AltScreen = true

	return v
}

func (m Model) renderHeader() string {
	tabStyle := m.style.
		Border(styles.Tab).
		Padding(0, 1, 0)

	primaryTabStyle := tabStyle.Border(styles.ActiveTab).
		Foreground(lipgloss.Blue)

	tabGapStyle := tabStyle.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	tabList := []string{}

	for i, tab := range m.tabs {
		if i == m.t {
			tabList = append(tabList, primaryTabStyle.Render(tab.Name()))

			continue
		}

		tabList = append(tabList, tabStyle.Render(tab.Name()))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Left, tabList...)
	tabsWidth := lipgloss.Width(tabs)
	tabsHeight := lipgloss.Height(tabs)

	gap := tabGapStyle.Height(tabsHeight).
		Render(strings.Repeat(" ", max(0, m.w-tabsWidth)))

	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, gap)
}
