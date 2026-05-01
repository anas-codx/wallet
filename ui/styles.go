package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginRight(2).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#F25D94")).
			Bold(true)

	ItemStyle = lipgloss.NewStyle().PaddingLeft(4)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170")).
				Bold(true)

	CheckMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")

	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
)
