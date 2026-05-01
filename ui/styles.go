package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	MainColor      = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	ErrorColor     = lipgloss.Color("#EF4444")
	BgColor        = lipgloss.Color("#2B2D42")
	TextColor      = lipgloss.Color("#EDF2F4")
	DimTextColor   = lipgloss.Color("#8D99AE")

	// Base Layout
	WindowStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MainColor).
			Background(BgColor)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(MainColor).
			Padding(0, 1).
			Bold(true).
			MarginBottom(1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(DimTextColor).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(DimTextColor).
			MarginTop(1).
			PaddingTop(1)

	// Menu & Lists
	ItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true).
				PaddingLeft(1)

	// Components
	TitleStyle = lipgloss.NewStyle().
			Foreground(MainColor).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Italic(true).
			PaddingLeft(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(MainColor).
			Bold(true).
			MarginTop(1)

	// Icons
	IconBullet   = "•"
	IconSelected = "➜"
	IconCheck    = "✓"
	IconLock     = "🔒"
	IconKey      = "🔑"
)
