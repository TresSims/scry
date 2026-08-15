package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui/styles"
)

// [Model] implements [tea.Model]
type Model struct {
	// viewport [w]idth, [h]eight
	w, h int

	// index of the current tab
	t int

	// tabs has a list of tabs that can be rendered
	tabs []Tab

	// the default stlye
	style lipgloss.Style

	// the most recent facts collected by the engine
	facts facts.Snapshot
}

// FactsMsg carries a new [facts.Snapshot] into the model. The fact engine runs
// outside of bubbletea, so its results arrive through [tea.Program.Send].
type FactsMsg struct {
	Facts facts.Snapshot
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

// WithFacts seeds the model with a [facts.Snapshot] so the first frame has
// content to render before the engine's next collection pass.
func WithFacts(f facts.Snapshot) Option {
	return func(m *Model) {
		m.facts = f
	}
}

// WithTab adds a [Tab] to the [Model]
func WithTab(t Tab) Option {
	return func(m *Model) {
		m.tabs = append(m.tabs, t)
	}
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

	case FactsMsg:
		m.facts = msg.Facts
	}

	return m, nil
}

func (m Model) View() tea.View {
	header := m.renderHeader()
	headerHeight := lipgloss.Height(header)

	tab := m.tabs[m.t].Render(m.w, m.h-headerHeight, m.facts, m.style)

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
