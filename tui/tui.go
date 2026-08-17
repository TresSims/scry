package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TresSims/scry/facts"
	"github.com/TresSims/scry/tui/internal/styles"
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
	facts facts.Cache
}

// FactsMsg carries a new [facts.Cache] into the model. The fact engine runs
// outside of bubbletea, so its results arrive through [tea.Program.Send].
type FactsMsg struct {
	Facts facts.Cache
}

type Option func(*Model)

func New(opts ...Option) *Model {
	model := &Model{}

	for _, opt := range opts {
		opt(model)
	}

	model.t = 0

	model.style = styles.MainStyle

	return model
}

// WithFacts seeds the model with a [facts.Cache] so the first frame has
// content to render before the engine's next collection pass.
func WithFacts(f facts.Cache) Option {
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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

		headerHeight := lipgloss.Height(m.renderHeader())

		for _, t := range m.tabs {
			t.SetSize(m.w, m.h-headerHeight)
		}

	case FactsMsg:
		m.facts = msg.Facts
	}

	tabModel, cmd := m.tabs[m.t].Update(msg)
	cmds = append(cmds, cmd)

	// Error discarded because [Tab] explicitly implements [tea.Model]
	m.tabs[m.t], _ = tabModel.(Tab)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	header := m.renderHeader()

	tabView := m.tabs[m.t].View()

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, header, tabView.Content))

	v.AltScreen = true

	return v
}

func (m Model) renderHeader() string {
	tabList := []string{}

	for i, tab := range m.tabs {
		if i == m.t {
			tabList = append(tabList, styles.PrimaryTabStyle.Render(tab.Name()))

			continue
		}

		tabList = append(tabList, styles.TabStyle.Render(tab.Name()))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Left, tabList...)
	tabsWidth := lipgloss.Width(tabs)
	tabsHeight := lipgloss.Height(tabs)

	gap := styles.GapTabStyle.Height(tabsHeight).
		Render(strings.Repeat(" ", max(0, m.w-tabsWidth)))

	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, gap)
}
