// Package style holds all colors and lipgloss styles.
// Every module imports this. The tui package also imports this.
// Because style imports nothing from this project, there is no import cycle.
package style

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("#7C3AED")
	ColorAccent    = lipgloss.Color("#10B981")
	ColorWarning   = lipgloss.Color("#F59E0B")
	ColorDanger    = lipgloss.Color("#EF4444")
	ColorMuted     = lipgloss.Color("#6B7280")
	ColorBorder    = lipgloss.Color("#374151")
	ColorSidebarBg = lipgloss.Color("#1F2937")
	ColorText      = lipgloss.Color("#F9FAFB")
	ColorSubtext   = lipgloss.Color("#9CA3AF")
)

var SidebarStyle = lipgloss.NewStyle().
	Background(ColorSidebarBg).
	Padding(1, 2)

var SidebarItemStyle = lipgloss.NewStyle().
	Foreground(ColorSubtext).
	PaddingLeft(1)

var SidebarItemActiveStyle = lipgloss.NewStyle().
	Foreground(ColorText).
	Background(ColorPrimary).
	PaddingLeft(1).
	Bold(true)

var TitleStyle = lipgloss.NewStyle().
	Foreground(ColorPrimary).
	Bold(true).
	MarginBottom(1)

var BoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorBorder).
	Padding(0, 1)

var SuccessStyle = lipgloss.NewStyle().Foreground(ColorAccent)
var DangerStyle  = lipgloss.NewStyle().Foreground(ColorDanger)
var WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning)
var MutedStyle   = lipgloss.NewStyle().Foreground(ColorMuted)
var BoldStyle    = lipgloss.NewStyle().Foreground(ColorText).Bold(true)

var HelpStyle = lipgloss.NewStyle().
	Foreground(ColorMuted).
	BorderTop(true).
	BorderForeground(ColorBorder)
