// Package monitor is the system monitor screen.
// It shows live CPU, RAM, disk usage and a process list.
// You can kill and inspect processes from here without typing commands.
package monitor

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
	"github.com/vpsmanager/vps/pkg/osadapter"
)

// DataMsg carries fresh system stats.
type DataMsg struct {
	cpuPercent int
	ramUsed    int
	ramTotal   int
	diskUsed   string
	diskTotal  string
	diskPct    int
	processes  []ProcessInfo
	netStats   string
}

// ProcessInfo is one line from the process list.
type ProcessInfo struct {
	PID     string
	User    string
	CPU     string
	Memory  string
	Command string
}

// Model holds the monitor screen state.
type Model struct {
	data    DataMsg
	cursor  int
	loaded  bool
	message string
}

func New() Model { return Model{} }

// LoadCmd fetches all system stats in the background.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchStats()
	}
}

// fetchStats collects CPU, RAM, disk, and process info.
func fetchStats() DataMsg {
	data := DataMsg{}

	// CPU — same approach as dashboard
	cpuResult := core.RunCommand("top -bn1 | grep 'Cpu(s)' | awk '{print $8}' | cut -d. -f1")
	if cpuResult.Err == nil {
		idle := 0
		if _, err := fmt.Sscanf(cpuResult.Output, "%d", &idle); err == nil {
			data.cpuPercent = 100 - idle
		}
	}

	// RAM from /proc/meminfo
	mem := osadapter.ReadMemInfo()
	data.ramUsed = mem.UsedMB
	data.ramTotal = mem.TotalMB

	// Disk usage
	diskResult := core.RunCommand("df -h / | tail -1 | awk '{print $2, $3, $5}'")
	if diskResult.Err == nil {
		parts := strings.Fields(diskResult.Output)
		if len(parts) >= 3 {
			data.diskTotal = parts[0]
			data.diskUsed = parts[1]
			pct := strings.TrimSuffix(parts[2], "%")
			_, _ = fmt.Sscanf(pct, "%d", &data.diskPct) // intentionally ignoring parse error for disk percentage
		}
	}

	// Process list — top 15 by CPU usage
	// `ps aux` lists all processes; sort by CPU (column 3) descending
	psResult := core.RunCommand("ps aux --sort=-%cpu | head -16")
	if psResult.Err == nil {
		lines := strings.Split(psResult.Output, "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue // skip header
			}
			fields := strings.Fields(line)
			if len(fields) < 11 {
				continue
			}
			// Join the command parts (can have spaces)
			cmd := strings.Join(fields[10:], " ")
			data.processes = append(data.processes, ProcessInfo{
				PID:     fields[1],
				User:    fields[0],
				CPU:     fields[2] + "%",
				Memory:  fields[3] + "%",
				Command: truncate(cmd, 40),
			})
		}
	}

	// Network stats — just a simple snapshot of bytes sent/received
	netResult := core.RunCommand("cat /proc/net/dev | grep -E 'eth0|ens|enp' | head -1 | awk '{print $2, $10}'")
	if netResult.Err == nil && netResult.Output != "" {
		parts := strings.Fields(netResult.Output)
		if len(parts) == 2 {
			var rx, tx int64
			_, _ = fmt.Sscanf(parts[0], "%d", &rx) // intentionally ignoring parse error for network stats
			_, _ = fmt.Sscanf(parts[1], "%d", &tx) // intentionally ignoring parse error for network stats
			data.netStats = fmt.Sprintf("RX: %s  TX: %s",
				formatBytes(rx), formatBytes(tx))
		}
	}

	return data
}

// formatBytes converts bytes to human-readable form.
func formatBytes(b int64) string {
	switch {
	case b >= 1024*1024*1024:
		return fmt.Sprintf("%.1fGB", float64(b)/(1024*1024*1024))
	case b >= 1024*1024:
		return fmt.Sprintf("%.1fMB", float64(b)/(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.1fKB", float64(b)/1024)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

// Update handles input on the monitor screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DataMsg:
		m.data = msg
		m.loaded = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.data.processes)-1 {
				m.cursor++
			}

		// Kill selected process with SIGTERM (graceful)
		case "x":
			if len(m.data.processes) > 0 {
				proc := m.data.processes[m.cursor]
				m.message = fmt.Sprintf("Press X to confirm kill PID %s (%s)", proc.PID, proc.Command)
			}

		// Confirm kill
		case "X":
			if len(m.data.processes) > 0 {
				proc := m.data.processes[m.cursor]
				result := core.RunCommand("kill " + proc.PID)
				if result.Err != nil {
					m.message = "Kill failed (try with sudo): " + proc.PID
				} else {
					m.message = "Killed PID " + proc.PID
					return m, m.LoadCmd()
				}
			}

		// Force kill with SIGKILL (last resort)
		case "K":
			if len(m.data.processes) > 0 {
				proc := m.data.processes[m.cursor]
				result := core.RunCommand("kill -9 " + proc.PID)
				if result.Err != nil {
					m.message = "Force kill failed: " + proc.PID
				} else {
					m.message = "Force killed PID " + proc.PID
					return m, m.LoadCmd()
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

// View renders the monitor screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Collecting system stats...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("System Monitor") + "\n\n")

	// Resource summary at the top
	b.WriteString(renderBar("CPU ", m.data.cpuPercent, width/2) + "\n")
	ramPct := 0
	if m.data.ramTotal > 0 {
		ramPct = (m.data.ramUsed * 100) / m.data.ramTotal
	}
	b.WriteString(renderBarLabel("RAM ", ramPct,
		fmt.Sprintf("%dMB/%dMB", m.data.ramUsed, m.data.ramTotal), width/2) + "\n")
	b.WriteString(renderBarLabel("Disk", m.data.diskPct,
		m.data.diskUsed+"/"+m.data.diskTotal, width/2) + "\n")

	if m.data.netStats != "" {
		b.WriteString("  " + style.MutedStyle.Render("Net: "+m.data.netStats) + "\n")
	}
	b.WriteString("\n")

	// Process list
	b.WriteString(style.BoldStyle.Render(fmt.Sprintf(
		"  %-8s %-12s %-6s %-6s %s\n", "PID", "User", "CPU", "Mem", "Command",
	)))
	b.WriteString(style.MutedStyle.Render("  " + strings.Repeat("─", width-6) + "\n"))

	for i, proc := range m.data.processes {
		line := fmt.Sprintf("  %-8s %-12s %-6s %-6s %s",
			proc.PID,
			truncate(proc.User, 11),
			proc.CPU,
			proc.Memory,
			proc.Command,
		)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  ↑↓: navigate  x: kill (graceful)  K: force kill  r: refresh"))
	}

	return b.String()
}

// renderBar draws a percent bar.
func renderBar(label string, percent int, width int) string {
	return renderBarLabel(label, percent, fmt.Sprintf("%d%%", percent), width)
}

// renderBarLabel draws a bar with a custom right label.
func renderBarLabel(label string, percent int, rightLabel string, barWidth int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	inner := barWidth - 16
	if inner < 4 {
		inner = 10
	}
	filled := (percent * inner) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", inner-filled)
	var barColored string
	switch {
	case percent >= 90:
		barColored = style.DangerStyle.Render(bar)
	case percent >= 70:
		barColored = style.WarningStyle.Render(bar)
	default:
		barColored = style.SuccessStyle.Render(bar)
	}
	return fmt.Sprintf("  %-5s [%s] %s", label, barColored, rightLabel)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
