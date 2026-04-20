// Package node manages Node.js apps running under PM2.
// PM2 is a process manager — it keeps your Node apps alive and restarts them if they crash.
// This screen lists all PM2 apps and lets you start, stop, restart, and watch logs.
package node

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// DataMsg is sent when PM2 app list is loaded.
type DataMsg struct {
	apps []AppInfo
	pm2  bool   // is PM2 installed at all?
	node bool   // is Node.js installed?
	err  string
}

// AppInfo holds info about one PM2-managed app.
type AppInfo struct {
	ID     string // PM2 app ID number
	Name   string // app name
	Status string // "online", "stopped", "errored"
	CPU    string // CPU percentage
	Memory string // memory usage
	Uptime string // how long it's been running
}

// Model holds state for the Node/PM2 screen.
type Model struct {
	apps    []AppInfo
	cursor  int
	loaded  bool
	pm2     bool   // PM2 is installed
	node    bool   // Node is installed
	message string // status message after an action
	logMode bool   // are we watching logs right now?
	logs    string // last fetched log output
}

func New() Model { return Model{} }

// LoadCmd fetches PM2 app list in the background.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		data := DataMsg{}

		// Check if Node and PM2 are installed
		data.node = core.CommandExists("node")
		data.pm2 = core.CommandExists("pm2")

		if !data.pm2 {
			return data // PM2 not installed, nothing more to do
		}

		// `pm2 jlist` outputs JSON — easier to parse than the table format
		// Each line contains app info separated by pipes when we use this format
		result := core.RunCommand("pm2 list --no-color 2>/dev/null")
		if result.Err != nil {
			data.err = result.Output
			return data
		}

		// Parse the PM2 table output.
		// The table looks like:
		// │ id │ name │ ... │ status │ cpu │ mem │
		// We skip header/separator lines and extract real rows.
		lines := strings.Split(result.Output, "\n")
		for _, line := range lines {
			// Real data rows start with "│" and contain actual values
			if !strings.HasPrefix(line, "│") {
				continue
			}
			// Split on │ and clean up each cell
			cells := strings.Split(line, "│")
			if len(cells) < 10 {
				continue
			}

			// Skip header row (contains "id" literally)
			id := strings.TrimSpace(cells[1])
			if id == "id" || id == "" {
				continue
			}

			app := AppInfo{
				ID:     id,
				Name:   strings.TrimSpace(cells[2]),
				Status: strings.TrimSpace(cells[9]),  // status column
				CPU:    strings.TrimSpace(cells[11]), // cpu column
				Memory: strings.TrimSpace(cells[12]), // mem column
				Uptime: strings.TrimSpace(cells[8]),  // uptime column
			}
			if app.Name != "" && app.ID != "" {
				data.apps = append(data.apps, app)
			}
		}

		data.pm2 = true
		return data
	}
}

// Update handles input on the Node/PM2 screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DataMsg:
		m.apps = msg.apps
		m.pm2 = msg.pm2
		m.node = msg.node
		m.loaded = true
		if msg.err != "" {
			m.message = "Error: " + msg.err
		}
		return m, nil

	case tea.KeyMsg:
		// Exit log mode with any key
		if m.logMode {
			m.logMode = false
			m.logs = ""
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.apps)-1 {
				m.cursor++
			}

		// Start selected app
		case "s":
			if len(m.apps) > 0 {
				app := m.apps[m.cursor]
				result := core.RunCommand("pm2 start " + app.ID)
				if result.Err != nil {
					m.message = "Failed to start: " + result.Output
				} else {
					m.message = "Started: " + app.Name
					return m, m.LoadCmd()
				}
			}

		// Stop selected app
		case "S":
			if len(m.apps) > 0 {
				app := m.apps[m.cursor]
				result := core.RunCommand("pm2 stop " + app.ID)
				if result.Err != nil {
					m.message = "Failed to stop: " + result.Output
				} else {
					m.message = "Stopped: " + app.Name
					return m, m.LoadCmd()
				}
			}

		// Restart selected app
		case "R":
			if len(m.apps) > 0 {
				app := m.apps[m.cursor]
				result := core.RunCommand("pm2 restart " + app.ID)
				if result.Err != nil {
					m.message = "Failed to restart: " + result.Output
				} else {
					m.message = "Restarted: " + app.Name
					return m, m.LoadCmd()
				}
			}

		// View last 50 lines of logs for selected app
		case "l":
			if len(m.apps) > 0 {
				app := m.apps[m.cursor]
				result := core.RunCommand("pm2 logs " + app.ID + " --lines 50 --nostream 2>/dev/null")
				m.logMode = true
				m.logs = result.Output
				if m.logs == "" {
					m.logs = "(no logs yet)"
				}
			}

		// Refresh
		case "r":
			m.loaded = false
			return m, m.LoadCmd()
		}
	}
	return m, nil
}

// View renders the Node/PM2 screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Loading PM2 apps...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Node / PM2") + "\n\n")

	// Show install status
	nodeStatus := style.SuccessStyle.Render("✓ installed")
	if !m.node {
		nodeStatus = style.DangerStyle.Render("✗ not installed")
	}
	pm2Status := style.SuccessStyle.Render("✓ installed")
	if !m.pm2 {
		pm2Status = style.DangerStyle.Render("✗ not installed")
	}
	b.WriteString(fmt.Sprintf("  Node.js: %s    PM2: %s\n\n", nodeStatus, pm2Status))

	// If we're in log view mode, show logs instead of the app list
	if m.logMode {
		b.WriteString(style.BoldStyle.Render("  Logs (press any key to go back):\n\n"))
		// Show last N lines that fit the screen
		logLines := strings.Split(m.logs, "\n")
		maxLines := 25
		start := 0
		if len(logLines) > maxLines {
			start = len(logLines) - maxLines
		}
		for _, line := range logLines[start:] {
			b.WriteString("  " + style.MutedStyle.Render(line) + "\n")
		}
		return b.String()
	}

	// Not installed — show install instructions
	if !m.pm2 {
		b.WriteString(style.MutedStyle.Render("  PM2 is not installed.\n\n"))
		b.WriteString("  To install: " + style.BoldStyle.Render("npm install -g pm2") + "\n")
		b.WriteString("  Or press " + style.BoldStyle.Render("i") + " to install now (requires npm)\n")
		return b.String()
	}

	// App list header
	b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-4s %-20s %-10s %-8s %-10s\n",
		"ID", "Name", "Status", "CPU", "Memory")))
	b.WriteString(style.MutedStyle.Render("  " + strings.Repeat("─", width-6) + "\n"))

	if len(m.apps) == 0 {
		b.WriteString(style.MutedStyle.Render("\n  No PM2 apps running. Start one with: pm2 start app.js\n"))
	}

	for i, app := range m.apps {
		// Status indicator with color
		statusStyled := ""
		switch app.Status {
		case "online":
			statusStyled = style.SuccessStyle.Render("● online   ")
		case "stopped":
			statusStyled = style.MutedStyle.Render("○ stopped  ")
		default:
			statusStyled = style.DangerStyle.Render("✗ " + app.Status + "  ")
		}

		line := fmt.Sprintf("  %-4s %-20s %s %-8s %-10s",
			app.ID, truncate(app.Name, 20), statusStyled, app.CPU, app.Memory)

		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}

	// Status message or help
	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  s: start  S: stop  R: restart  l: logs  r: refresh"))
	}

	return b.String()
}

// truncate shortens a string to maxLen and adds "…" if it was cut.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
