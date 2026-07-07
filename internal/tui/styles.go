package tui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(0, 1)

	OrgPanelStyle = lipgloss.NewStyle().
			Width(22).
			Border(lipgloss.RoundedBorder())

	OrgPanelFocusedStyle = lipgloss.NewStyle().
				Width(22).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#5865F2"))

	ChannelPanelStyle = lipgloss.NewStyle().
				Width(26).
				Border(lipgloss.RoundedBorder())

	ChannelPanelFocusedStyle = lipgloss.NewStyle().
					Width(26).
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#57F287"))

	ChatPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FEE75C"))

	DimmedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#3a3a3a"))

	StatusStyle = lipgloss.NewStyle().
			Height(1).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5865F2")).
			Padding(0, 1).
			Bold(true)

	SelectedOrgStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#5865F2")).
				Padding(0, 1).
				Bold(true)

	SelectedChannelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#57F287")).
				Padding(0, 1).
				Bold(true)

	OrgItemStyle = lipgloss.NewStyle().Padding(0, 1)

	ChannelItemStyle = lipgloss.NewStyle().Padding(0, 1)

	SenderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#57F287"))

	TimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Width(8).
			Align(lipgloss.Right)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FEE75C"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ED4245"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#57F287"))

	ModeBadgeInput = lipgloss.NewStyle().
			Background(lipgloss.Color("#FEE75C")).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)

	ModeBadgeNav = lipgloss.NewStyle().
			Background(lipgloss.Color("#5865F2")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(60)

	ChannelHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#ED4245")).
				Padding(0, 1).
				Bold(true)

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)
