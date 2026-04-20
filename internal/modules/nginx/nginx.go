// Package nginx manages the Nginx web server.
// It shows all site configs, lets you enable/disable sites, edit configs,
// check SSL cert expiry, and install nginx if it isn't there yet.
package nginx

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// DataMsg is sent when we finish loading nginx info.
type DataMsg struct {
	installed bool
	running   bool
	sites     []SiteInfo
	err       string
}

// SiteInfo represents one nginx site config file.
type SiteInfo struct {
	Name    string // filename like "mysite.com"
	Enabled bool   // is there a symlink in sites-enabled?
	SSL     string // cert expiry info, or "no SSL"
}

// Mode tracks which sub-view we're in.
type Mode int

const (
	ModeList   Mode = iota // list of sites
	ModeConfig             // viewing/editing a config file
)

// Model holds the nginx screen state.
type Model struct {
	installed bool
	running   bool
	sites     []SiteInfo
	cursor    int
	loaded    bool
	mode      Mode
	configContent string // raw text of the selected config
	message  string
}

func New() Model { return Model{} }

// LoadCmd loads nginx status and all site configs.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchNginxData()
	}
}

// fetchNginxData checks if nginx is installed, running, and reads all site configs.
func fetchNginxData() DataMsg {
	data := DataMsg{}

	// Is nginx installed?
	data.installed = core.CommandExists("nginx")
	if !data.installed {
		return data
	}

	// Is it running right now?
	statusResult := core.RunCommand("systemctl is-active nginx 2>/dev/null")
	data.running = statusResult.Output == "active"

	// Get list of available sites from sites-available
	// Each file in this dir is one site config
	availDir := "/etc/nginx/sites-available"
	entries, err := os.ReadDir(availDir)
	if err != nil {
		data.err = "Cannot read " + availDir + ": " + err.Error()
		return data
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		// Check if this site is enabled by looking for a symlink in sites-enabled
		enabledPath := "/etc/nginx/sites-enabled/" + name
		_, linkErr := os.Lstat(enabledPath)
		enabled := linkErr == nil // symlink exists = enabled

		// Check for SSL cert for this site
		// Common cert locations: /etc/letsencrypt/live/<domain>/
		sslInfo := checkSSL(name)

		data.sites = append(data.sites, SiteInfo{
			Name:    name,
			Enabled: enabled,
			SSL:     sslInfo,
		})
	}

	return data
}

// checkSSL looks for a Let's Encrypt cert for the given site name
// and returns the expiry date, or "no SSL" if none found.
func checkSSL(siteName string) string {
	// Try to find a cert matching this site name
	certPath := fmt.Sprintf("/etc/letsencrypt/live/%s/cert.pem", siteName)
	result := core.RunCommand(fmt.Sprintf(
		"openssl x509 -enddate -noout -in %s 2>/dev/null | cut -d= -f2", certPath,
	))
	if result.Err != nil || result.Output == "" {
		return "no SSL"
	}
	return "SSL: " + result.Output
}

// Update handles keyboard input on the nginx screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DataMsg:
		m.installed = msg.installed
		m.running = msg.running
		m.sites = msg.sites
		m.loaded = true
		if msg.err != "" {
			m.message = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		// In config view mode, esc goes back to list
		if m.mode == ModeConfig {
			if msg.String() == "esc" || msg.String() == "backspace" {
				m.mode = ModeList
				m.configContent = ""
			}
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.sites)-1 {
				m.cursor++
			}

		// View the raw config file for selected site
		case "enter", "v":
			if len(m.sites) > 0 {
				site := m.sites[m.cursor]
				configPath := "/etc/nginx/sites-available/" + site.Name
				content, err := os.ReadFile(configPath)
				if err != nil {
					m.message = "Cannot read config: " + err.Error()
				} else {
					m.configContent = string(content)
					m.mode = ModeConfig
				}
			}

		// Enable selected site (create symlink in sites-enabled)
		case "e":
			if len(m.sites) > 0 {
				site := m.sites[m.cursor]
				if site.Enabled {
					m.message = site.Name + " is already enabled"
					break
				}
				src := "/etc/nginx/sites-available/" + site.Name
				dst := "/etc/nginx/sites-enabled/" + site.Name
				result := core.RunCommand(fmt.Sprintf("ln -s %q %q 2>&1", src, dst))
				if result.Err != nil {
					m.message = "Enable failed: " + result.Output
				} else {
					m.message = "Enabled: " + site.Name + " (reload nginx with R)"
					return m, m.LoadCmd()
				}
			}

		// Disable selected site (remove symlink from sites-enabled)
		case "d":
			if len(m.sites) > 0 {
				site := m.sites[m.cursor]
				if !site.Enabled {
					m.message = site.Name + " is already disabled"
					break
				}
				dst := "/etc/nginx/sites-enabled/" + site.Name
				result := core.RunCommand(fmt.Sprintf("rm -f %q 2>&1", dst))
				if result.Err != nil {
					m.message = "Disable failed: " + result.Output
				} else {
					m.message = "Disabled: " + site.Name + " (reload nginx with R)"
					return m, m.LoadCmd()
				}
			}

		// Test nginx config for syntax errors
		case "t":
			result := core.RunCommand("nginx -t 2>&1")
			m.message = result.Output
			if len(m.message) > 100 {
				m.message = m.message[:100]
			}

		// Reload nginx (pick up config changes without dropping connections)
		case "R":
			result := core.RunCommand("systemctl reload nginx 2>&1")
			if result.Err != nil {
				m.message = "Reload failed: " + result.Output
			} else {
				m.message = "Nginx reloaded successfully"
				return m, m.LoadCmd()
			}

		// Install nginx if not installed
		case "i":
			if !m.installed {
				m.message = "Installing nginx... (this may take a moment)"
				return m, func() tea.Msg {
					result := core.RunCommand("apt-get install -y nginx 2>&1")
					if result.Err != nil {
						return DataMsg{err: "Install failed: " + result.Output}
					}
					return fetchNginxData()
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

// View renders the nginx screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Loading nginx info...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Nginx") + "\n\n")

	// Not installed — offer to install
	if !m.installed {
		b.WriteString(style.DangerStyle.Render("  ✗ Nginx is not installed\n\n"))
		b.WriteString("  Press " + style.BoldStyle.Render("i") + " to install nginx via apt\n")
		if m.message != "" {
			b.WriteString("\n" + style.WarningStyle.Render("  "+m.message))
		}
		return b.String()
	}

	// Status line
	if m.running {
		b.WriteString("  Status: " + style.SuccessStyle.Render("● running") + "\n\n")
	} else {
		b.WriteString("  Status: " + style.DangerStyle.Render("● stopped") + "\n\n")
	}

	// Config view mode — show the raw config file
	if m.mode == ModeConfig {
		if len(m.sites) > 0 {
			b.WriteString(style.BoldStyle.Render("  "+m.sites[m.cursor].Name) + "\n")
			b.WriteString(style.MutedStyle.Render("  /etc/nginx/sites-available/"+m.sites[m.cursor].Name) + "\n\n")
		}
		// Show config lines with basic syntax highlighting
		for _, line := range strings.Split(m.configContent, "\n") {
			trimmed := strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(trimmed, "#"):
				// Comments in muted gray
				b.WriteString(style.MutedStyle.Render("  "+line) + "\n")
			case strings.Contains(line, "server_name") || strings.Contains(line, "listen"):
				// Key directives in bold
				b.WriteString("  " + style.BoldStyle.Render(line) + "\n")
			case strings.Contains(line, "ssl") || strings.Contains(line, "SSL"):
				// SSL lines in green
				b.WriteString("  " + style.SuccessStyle.Render(line) + "\n")
			default:
				b.WriteString("  " + line + "\n")
			}
		}
		b.WriteString("\n" + style.MutedStyle.Render("  backspace / esc: back to sites list"))
		return b.String()
	}

	// Sites list
	if len(m.sites) == 0 {
		b.WriteString(style.MutedStyle.Render("  No site configs found in /etc/nginx/sites-available/\n"))
		b.WriteString(style.MutedStyle.Render("  Create a config file there to get started.\n"))
	} else {
		b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-30s %-10s %s\n", "Site", "Status", "SSL")))
		b.WriteString(style.MutedStyle.Render("  "+strings.Repeat("─", width-6)) + "\n")

		for i, site := range m.sites {
			statusStr := ""
			if site.Enabled {
				statusStr = style.SuccessStyle.Render("● enabled ")
			} else {
				statusStr = style.MutedStyle.Render("○ disabled")
			}

			sslStr := ""
			if site.SSL == "no SSL" {
				sslStr = style.MutedStyle.Render("no SSL")
			} else {
				sslStr = style.SuccessStyle.Render(site.SSL)
			}

			line := fmt.Sprintf("  %-30s %s  %s", truncate(site.Name, 29), statusStr, sslStr)

			if i == m.cursor {
				b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
	}

	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  enter: view config  e: enable  d: disable  t: test  R: reload  r: refresh"))
	}

	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
