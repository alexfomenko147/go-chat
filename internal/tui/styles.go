package tui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
			Padding(0, 1)

	OrgPanelStyle = lipgloss.NewStyle().
			Width(20).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5865F2"))

	ChannelPanelStyle = lipgloss.NewStyle().
				Width(24).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#57F287"))

	ChatPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ED4245"))

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FEE75C"))

	StatusStyle = lipgloss.NewStyle().
			Height(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5865F2"))

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5865F2")).
			Padding(0, 1).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5865F2"))

	MessageStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(80)

	SenderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#57F287"))

	TimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Width(10).
			Align(lipgloss.Right)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#57F287"))

	DimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FEE75C"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ED4245"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#57F287"))

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)
