package styles

import "charm.land/lipgloss/v2"

var (
	ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	TabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	TabStyle = MainStyle.
			Border(ActiveTabBorder).
			Padding(0, 1, 0)

	PrimaryTabStyle = TabStyle.Border(ActiveTabBorder).
			Foreground(lipgloss.Blue)

	GapTabStyle = TabStyle.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
)
