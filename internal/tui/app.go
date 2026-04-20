// app.go is the main BubbleTea application.
// BubbleTea works in three steps that repeat forever:
//   1. Init()   — runs once at startup, can start background tasks
//   2. Update() — receives a keypress or message, returns new state + optional command
//   3. View()   — reads current state, returns the string to draw on screen
//
// Think of it like a game loop: update state → draw screen → wait for input → repeat.
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// Each module provides its own Model that we embed here
	"github.com/vpsmanager/vps/internal/modules/dashboard"
	"github.com/vpsmanager/vps/internal/modules/files"
	"github.com/vpsmanager/vps/internal/modules/git"
	"github.com/vpsmanager/vps/internal/modules/monitor"
	"github.com/vpsmanager/vps/internal/modules/nginx"
	"github.com/vpsmanager/vps/internal/modules/node"
	"github.com/vpsmanager/vps/internal/modules/packages"
	"github.com/vpsmanager/vps/internal/modules/security"
)

// AppModel is the root model for the whole application.
// It holds the current screen and all module sub-models.
type AppModel struct {
	currentScreen Screen // which screen is showing right now
	width         int    // terminal width (updated on resize)
	height        int    // terminal height (updated on resize)

	// One sub-model per module — each manages its own state
	dashboardModel dashboard.Model
	filesModel     files.Model
	nodeModel      node.Model
	nginxModel     nginx.Model
	gitModel       git.Model
	packagesModel  packages.Model
	monitorModel   monitor.Model
	securityModel  security.Model
}

// NewApp creates a fresh AppModel with all modules initialized.
// Call this once at startup.
func NewApp() AppModel {
	return AppModel{
		currentScreen:  ScreenDashboard, // start on the dashboard
		dashboardModel: dashboard.New(),
		filesModel:     files.New(),
		nodeModel:      node.New(),
		nginxModel:     nginx.New(),
		gitModel:       git.New(),
		packagesModel:  packages.New(),
		monitorModel:   monitor.New(),
		securityModel:  security.New(),
	}
}

// Init runs once when the app starts.
// We return a command that loads the dashboard data immediately.
func (a AppModel) Init() tea.Cmd {
	return a.dashboardModel.LoadCmd()
}

// Update is called every time something happens: keypress, resize, data loaded.
// It returns the new state of the app + optionally a background command to run.
func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Terminal was resized — store new dimensions
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	// A key was pressed
	case tea.KeyMsg:
		switch msg.String() {

		// Quit the app
		case "q", "ctrl+c":
			return a, tea.Quit

		// Switch screens with number keys
		case "1":
			return a.switchTo(ScreenDashboard)
		case "2":
			return a.switchTo(ScreenFiles)
		case "3":
			return a.switchTo(ScreenNode)
		case "4":
			return a.switchTo(ScreenNginx)
		case "5":
			return a.switchTo(ScreenGit)
		case "6":
			return a.switchTo(ScreenPackages)
		case "7":
			return a.switchTo(ScreenMonitor)
		case "8":
			return a.switchTo(ScreenSecurity)

		// Tab cycles forward through screens
		case "tab":
			next := (int(a.currentScreen) + 1) % len(NavItems)
			return a.switchTo(Screen(next))

		// Shift+Tab cycles backward
		case "shift+tab":
			prev := (int(a.currentScreen) - 1 + len(NavItems)) % len(NavItems)
			return a.switchTo(Screen(prev))

		default:
			// Pass the keypress down to whichever module is active
			return a.updateActiveModule(msg)
		}

	default:
		// Pass any other message (like data-loaded events) to the active module
		return a.updateActiveModule(msg)
	}
}

// switchTo changes the current screen and fires that module's load command.
func (a AppModel) switchTo(s Screen) (tea.Model, tea.Cmd) {
	a.currentScreen = s
	// Tell the new module to load/refresh its data
	cmd := a.loadCmdForScreen(s)
	return a, cmd
}

// loadCmdForScreen returns the data-loading command for a given screen.
// Each module has its own LoadCmd that fetches fresh data from the system.
func (a AppModel) loadCmdForScreen(s Screen) tea.Cmd {
	switch s {
	case ScreenDashboard:
		return a.dashboardModel.LoadCmd()
	case ScreenFiles:
		return a.filesModel.LoadCmd()
	case ScreenNode:
		return a.nodeModel.LoadCmd()
	case ScreenNginx:
		return a.nginxModel.LoadCmd()
	case ScreenGit:
		return a.gitModel.LoadCmd()
	case ScreenPackages:
		return a.packagesModel.LoadCmd()
	case ScreenMonitor:
		return a.monitorModel.LoadCmd()
	case ScreenSecurity:
		return a.securityModel.LoadCmd()
	}
	return nil
}

// updateActiveModule passes messages to whichever module is currently showing.
// Each module handles its own keyboard input and updates its own state.
func (a AppModel) updateActiveModule(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch a.currentScreen {
	case ScreenDashboard:
		a.dashboardModel, cmd = a.dashboardModel.Update(msg)
	case ScreenFiles:
		a.filesModel, cmd = a.filesModel.Update(msg)
	case ScreenNode:
		a.nodeModel, cmd = a.nodeModel.Update(msg)
	case ScreenNginx:
		a.nginxModel, cmd = a.nginxModel.Update(msg)
	case ScreenGit:
		a.gitModel, cmd = a.gitModel.Update(msg)
	case ScreenPackages:
		a.packagesModel, cmd = a.packagesModel.Update(msg)
	case ScreenMonitor:
		a.monitorModel, cmd = a.monitorModel.Update(msg)
	case ScreenSecurity:
		a.securityModel, cmd = a.securityModel.Update(msg)
	}
	return a, cmd
}

// View draws the entire screen as a string.
// BubbleTea calls this after every Update() and renders whatever we return.
func (a AppModel) View() string {
	// If we don't know the terminal size yet, show a loading message
	if a.width == 0 {
		return "Loading VPS Manager..."
	}

	// Draw the sidebar on the left
	sidebar := renderSidebar(a.currentScreen, a.height)

	// Draw the active module's content on the right
	contentWidth := a.width - 24 // 24 = sidebar width + some padding
	content := a.viewActiveModule(contentWidth)

	// Place sidebar and content side by side
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)

	// Add the help bar at the bottom
	help := renderHelpBar(a.width)

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}

// viewActiveModule asks the current module to render itself.
func (a AppModel) viewActiveModule(width int) string {
	switch a.currentScreen {
	case ScreenDashboard:
		return a.dashboardModel.View(width)
	case ScreenFiles:
		return a.filesModel.View(width)
	case ScreenNode:
		return a.nodeModel.View(width)
	case ScreenNginx:
		return a.nginxModel.View(width)
	case ScreenGit:
		return a.gitModel.View(width)
	case ScreenPackages:
		return a.packagesModel.View(width)
	case ScreenMonitor:
		return a.monitorModel.View(width)
	case ScreenSecurity:
		return a.securityModel.View(width)
	}
	return ""
}

// renderHelpBar draws the bottom bar showing keyboard shortcuts.
func renderHelpBar(width int) string {
	shortcuts := "  1-8: switch screen  •  tab: next  •  q: quit  •  r: refresh"
	return HelpStyle.Width(width).Render(fmt.Sprintf("%-*s", width-2, shortcuts))
}
