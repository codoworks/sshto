package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("69")
	ColorSecondary = lipgloss.Color("241")
	ColorSuccess   = lipgloss.Color("42")
	ColorWarning   = lipgloss.Color("214")
	ColorDanger    = lipgloss.Color("196")

	// Group colors
	GroupColors = map[string]lipgloss.Color{
		"red":     lipgloss.Color("196"),
		"green":   lipgloss.Color("42"),
		"yellow":  lipgloss.Color("214"),
		"blue":    lipgloss.Color("69"),
		"magenta": lipgloss.Color("165"),
		"cyan":    lipgloss.Color("51"),
		"white":   lipgloss.Color("255"),
		"gray":    lipgloss.Color("241"),
	}

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorPrimary).
				Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginTop(1)

	GroupTagStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginRight(1)
)

// GroupTag returns a styled group tag
func GroupTag(name, color string) string {
	c, ok := GroupColors[color]
	if !ok {
		c = ColorSecondary
	}
	return GroupTagStyle.Copy().
		Background(c).
		Foreground(lipgloss.Color("0")).
		Render(name)
}
