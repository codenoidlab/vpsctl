// styles.go re-exports everything from the style package.
// Modules should import "internal/style" directly.
// This file stays for any tui-specific layout styles.
package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/vpsmanager/vps/internal/style"
)

// Re-export styles so existing code in this package still works
var (
	SidebarStyle           = style.SidebarStyle
	SidebarItemStyle       = style.SidebarItemStyle
	SidebarItemActiveStyle = style.SidebarItemActiveStyle
	TitleStyle             = style.TitleStyle
	BoxStyle               = style.BoxStyle
	SuccessStyle           = style.SuccessStyle
	DangerStyle            = style.DangerStyle
	WarningStyle           = style.WarningStyle
	MutedStyle             = style.MutedStyle
	BoldStyle              = style.BoldStyle
	HelpStyle              = style.HelpStyle
)

// ContentStyle is only used inside the tui layout, not in modules
var ContentStyle = lipgloss.NewStyle().Padding(1, 2)
