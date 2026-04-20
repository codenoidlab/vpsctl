// Package dashboard is the home screen.
// It shows CPU usage, RAM usage, disk usage, uptime, and running services.
// Think of it as a quick health check at a glance.
package dashboard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
	"github.com/vpsmanager/vps/pkg/osadapter"
)

// DataLoadedMsg is sent when background data fetching is done.
// BubbleTea uses messages to pass data from goroutines back to the UI.
type DataLoadedMsg struct {
	cpuPercent  int
	ramTotal    int
	ramUsed     int
	diskTotal   string
	diskUsed    string
	diskPercent int
	uptime      string
	services    []ServiceStatus
}

// ServiceStatus holds whether a system service is running or stopped.
type ServiceStatus struct {
	Name    string
	Running bool
}

// Model is the state for the dashboard screen.
// BubbleTea requires each screen to have its own Model + Update + View.
type Model struct {
	loaded bool          // have we fetched data yet?
	data   DataLoadedMsg // the last fetched data
}

// New creates a fresh dashboard model with no data yet.
func New() Model {
	return Model{}
}

// LoadCmd returns a BubbleTea command that fetches all dashboard data.
// It runs in a background goroutine and sends a DataLoadedMsg when done.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchDashboardData()
	}
}

// fetchDashboardData does the actual system calls to collect data.
// This runs in a goroutine (via LoadCmd) so the UI doesn't freeze.
func fetchDashboardData() DataLoadedMsg {
	data := DataLoadedMsg{}

	// --- CPU usage ---
	// `top -bn1` runs top once (non-interactive) and gives us CPU stats
	// We grep for the "Cpu(s)" line and pull out the idle percentage
	cpuResult := core.RunCommand("top -bn1 | grep 'Cpu(s)' | awk '{print $8}' | cut -d. -f1")
	if cpuResult.Err == nil && cpuResult.Output != "" {
		idle := 0
		if _, err := fmt.Sscanf(cpuResult.Output, "%d", &idle); err == nil {
			data.cpuPercent = 100 - idle // usage = 100 - idle
		}
	}

	// --- RAM usage (from /proc/meminfo) ---
	mem := osadapter.ReadMemInfo()
	data.ramTotal = mem.TotalMB
	data.ramUsed = mem.UsedMB

	// --- Disk usage ---
	// `df -h /` shows disk usage for the root partition in human-readable form
	diskResult := core.RunCommand("df -h / | tail -1 | awk '{print $2, $3, $5}'")
	if diskResult.Err == nil {
		parts := strings.Fields(diskResult.Output)
		if len(parts) >= 3 {
			data.diskTotal = parts[0]
			data.diskUsed = parts[1]
			// parts[2] looks like "45%" — strip the % and parse
			pct := strings.TrimSuffix(parts[2], "%")
			_ = fmt.Sscanf(pct, "%d", &data.diskPercent) // intentionally ignoring parse error for disk percentage
		}
	}

	// --- Uptime ---
	uptimeSeconds := osadapter.ReadUptime()
	data.uptime = osadapter.FormatUptime(uptimeSeconds)

	// --- Service statuses ---
	// Check whether common services are active using systemctl
	serviceNames := []string{"nginx", "mysql", "postgresql", "redis", "docker", "fail2ban"}
	for _, name := range serviceNames {
		// systemctl is-active returns exit code 0 if running
		result := core.RunCommand(fmt.Sprintf("systemctl is-active %s 2>/dev/null", name))
		running := result.Output == "active"
		// Only include services that are installed (active or failed, not "not-found")
		if result.Output != "not-found" && result.Output != "" {
			data.services = append(data.services, ServiceStatus{
				Name:    name,
				Running: running,
			})
		}
	}

	return data
}

// Update handles messages for the dashboard.
// The main thing it handles is receiving the DataLoadedMsg from LoadCmd.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Data finished loading — store it and mark as loaded
	case DataLoadedMsg:
		m.data = msg
		m.loaded = true
		return m, nil

	// 'r' key refreshes the dashboard data
	case tea.KeyMsg:
		if msg.String() == "r" {
			m.loaded = false
			return m, m.LoadCmd()
		}
	}
	return m, nil
}

// View renders the dashboard as a string.
// width is the available horizontal space after the sidebar.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Loading dashboard data...")
	}

	var b strings.Builder

	// Screen title
	b.WriteString(style.TitleStyle.Render("Dashboard") + "\n\n")

	// --- Resource bars row ---
	b.WriteString(renderBar("CPU", m.data.cpuPercent, 100, width/3))
	b.WriteString("\n")
	ramPct := 0
	if m.data.ramTotal > 0 {
		ramPct = (m.data.ramUsed * 100) / m.data.ramTotal
	}
	ramLabel := fmt.Sprintf("%dMB / %dMB", m.data.ramUsed, m.data.ramTotal)
	b.WriteString(renderBarLabel("RAM", ramPct, ramLabel, width/3))
	b.WriteString("\n")
	diskLabel := fmt.Sprintf("%s / %s", m.data.diskUsed, m.data.diskTotal)
	b.WriteString(renderBarLabel("Disk", m.data.diskPercent, diskLabel, width/3))
	b.WriteString("\n\n")

	// --- Uptime ---
	b.WriteString(style.BoldStyle.Render("Uptime: ") + m.data.uptime + "\n\n")

	// --- Services list ---
	b.WriteString(style.BoldStyle.Render("Services:\n"))
	if len(m.data.services) == 0 {
		b.WriteString(style.MutedStyle.Render("  No services detected\n"))
	}
	for _, svc := range m.data.services {
		if svc.Running {
			b.WriteString("  " + style.SuccessStyle.Render("● ") + svc.Name + "\n")
		} else {
			b.WriteString("  " + style.DangerStyle.Render("● ") + style.MutedStyle.Render(svc.Name) + "\n")
		}
	}

	b.WriteString("\n" + style.MutedStyle.Render("  press r to refresh"))

	return b.String()
}

// renderBar draws a percentage bar with a label.
// Example:  CPU  [████████░░]  45%
func renderBar(label string, value int, max int, barWidth int) string {
	return renderBarLabel(label, value, fmt.Sprintf("%d%%", value), barWidth)
}

// renderBarLabel draws a bar where the right side shows a custom label instead of just %.
func renderBarLabel(label string, percent int, rightLabel string, barWidth int) string {
	if barWidth < 10 {
		barWidth = 20 // minimum sensible width
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// How many filled blocks inside the bar
	innerWidth := barWidth - 10 // account for label + brackets + spaces
	if innerWidth < 4 {
		innerWidth = 4
	}
	filled := (percent * innerWidth) / 100

	// Choose bar color based on usage level
	bar := strings.Repeat("█", filled) + strings.Repeat("░", innerWidth-filled)
	barStyled := ""
	switch {
	case percent >= 90:
		barStyled = style.DangerStyle.Render(bar) // red when critical
	case percent >= 70:
		barStyled = style.WarningStyle.Render(bar) // amber when high
	default:
		barStyled = style.SuccessStyle.Render(bar) // green when ok
	}

	// Format: "CPU  [████░░░░]  512MB / 2048MB"
	return fmt.Sprintf("  %-6s [%s]  %s", label, barStyled, rightLabel)
}
