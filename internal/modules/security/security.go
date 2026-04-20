// Package security is the security control panel.
// It shows UFW firewall rules, SSH authorized keys, Fail2ban jail status,
// and the list of system user accounts.
// Tab switches between these four security areas.
package security

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// Which security section is active
type Tab int

const (
	TabFirewall Tab = iota // UFW firewall rules
	TabSSH                 // SSH authorized keys
	TabFail2ban            // Fail2ban jail status
	TabUsers               // system user accounts
)

// DataMsg carries all loaded security data.
type DataMsg struct {
	firewallActive bool
	firewallRules  []FirewallRule
	sshKeys        []SSHKey
	fail2banJails  []Fail2banJail
	users          []UserInfo
	err            string
}

// FirewallRule is one UFW rule (allow or deny for a port).
type FirewallRule struct {
	Number  string // rule number
	To      string // destination (port/protocol)
	Action  string // "ALLOW" or "DENY"
	From    string // "Anywhere" or specific IP
}

// SSHKey is one line from ~/.ssh/authorized_keys.
type SSHKey struct {
	Type    string // "ssh-rsa", "ssh-ed25519", etc.
	Comment string // usually user@machine
	Short   string // first 20 chars of the key for display
}

// Fail2banJail is one jail entry from fail2ban.
type Fail2banJail struct {
	Name    string
	Enabled bool
	Banned  int // number of currently banned IPs
}

// UserInfo is one system user account.
type UserInfo struct {
	Name  string
	UID   string
	Shell string
	Home  string
}

// Model holds the security screen state.
type Model struct {
	activeTab Tab
	data      DataMsg
	cursor    int
	loaded    bool
	message   string
}

func New() Model { return Model{activeTab: TabFirewall} }

// LoadCmd loads all security data in the background.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchSecurityData()
	}
}

// fetchSecurityData collects firewall rules, SSH keys, fail2ban, and users.
func fetchSecurityData() DataMsg {
	data := DataMsg{}

	// --- UFW Firewall ---
	// Check if UFW is active
	ufwStatus := core.RunCommand("ufw status 2>/dev/null")
	data.firewallActive = strings.Contains(ufwStatus.Output, "Status: active")

	// Parse the numbered rule list
	rulesResult := core.RunCommand("ufw status numbered 2>/dev/null")
	if rulesResult.Err == nil {
		for _, line := range strings.Split(rulesResult.Output, "\n") {
			// Rule lines look like: "[ 1] 22/tcp                     ALLOW IN    Anywhere"
			if !strings.HasPrefix(line, "[") {
				continue
			}
			// Strip the brackets and number
			line = strings.TrimLeft(line, "[ ")
			parts := strings.SplitN(line, "]", 2)
			if len(parts) < 2 {
				continue
			}
			num := strings.TrimSpace(parts[0])
			rest := strings.Fields(parts[1])
			if len(rest) < 3 {
				continue
			}
			rule := FirewallRule{
				Number: num,
				To:     rest[0],
				Action: rest[1],
			}
			if len(rest) > 3 {
				rule.From = strings.Join(rest[3:], " ")
			} else {
				rule.From = "Anywhere"
			}
			data.firewallRules = append(data.firewallRules, rule)
		}
	}

	// --- SSH Authorized Keys ---
	// Read the authorized_keys file for the current user
	keysResult := core.RunCommand("cat ~/.ssh/authorized_keys 2>/dev/null")
	if keysResult.Err == nil && keysResult.Output != "" {
		for _, line := range strings.Split(keysResult.Output, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			comment := ""
			if len(parts) >= 3 {
				comment = parts[2]
			}
			short := ""
			if len(parts[1]) > 20 {
				short = parts[1][:20] + "…"
			} else {
				short = parts[1]
			}
			data.sshKeys = append(data.sshKeys, SSHKey{
				Type:    parts[0],
				Comment: comment,
				Short:   short,
			})
		}
	}

	// --- Fail2ban ---
	// Get the status of all jails
	fail2banResult := core.RunCommand("fail2ban-client status 2>/dev/null")
	if fail2banResult.Err == nil {
		// Output: "Jail list:\t sshd, nginx-http-auth"
		for _, line := range strings.Split(fail2banResult.Output, "\n") {
			if strings.Contains(line, "Jail list:") {
				// Extract the jail names
				jailList := strings.TrimPrefix(line, "|- Jail list:")
				jailList = strings.TrimSpace(jailList)
				for _, jailName := range strings.Split(jailList, ",") {
					jailName = strings.TrimSpace(jailName)
					if jailName == "" {
						continue
					}
					// Get banned count for this jail
					jailStatus := core.RunCommand(fmt.Sprintf(
						"fail2ban-client status %s 2>/dev/null | grep 'Currently banned'",
						jailName,
					))
					banned := 0
					if jailStatus.Err == nil {
						_ = fmt.Sscanf(strings.TrimSpace(strings.Split(jailStatus.Output, ":")[1]), "%d", &banned) // intentionally ignoring parse error for fail2ban stats
					}
					data.fail2banJails = append(data.fail2banJails, Fail2banJail{
						Name:    jailName,
						Enabled: true,
						Banned:  banned,
					})
				}
				break
			}
		}
	}

	// --- System Users ---
	// Read /etc/passwd and show real accounts (UID >= 1000, plus root)
	usersResult := core.RunCommand("awk -F: '$3 >= 1000 || $1 == \"root\" {print $1, $3, $7, $6}' /etc/passwd 2>/dev/null")
	if usersResult.Err == nil {
		for _, line := range strings.Split(usersResult.Output, "\n") {
			parts := strings.Fields(line)
			if len(parts) < 4 {
				continue
			}
			data.users = append(data.users, UserInfo{
				Name:  parts[0],
				UID:   parts[1],
				Shell: parts[2],
				Home:  parts[3],
			})
		}
	}

	return data
}

// Update handles input on the security screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DataMsg:
		m.data = msg
		m.loaded = true
		m.cursor = 0
		if msg.err != "" {
			m.message = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		// Figure out how many items are in the current tab list
		listLen := m.currentListLen()

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < listLen-1 {
				m.cursor++
			}

		// Switch tabs with left/right
		case "left", "h":
			if m.activeTab > TabFirewall {
				m.activeTab--
				m.cursor = 0
			}
		case "right", "l":
			if m.activeTab < TabUsers {
				m.activeTab++
				m.cursor = 0
			}

		// Tab key cycles forward
		case "tab":
			m.activeTab = Tab((int(m.activeTab) + 1) % 4)
			m.cursor = 0

		// Enable UFW (if on firewall tab)
		case "e":
			if m.activeTab == TabFirewall && !m.data.firewallActive {
				result := core.RunCommand("ufw --force enable 2>&1")
				if result.Err != nil {
					m.message = "Enable failed: " + result.Output
				} else {
					m.message = "Firewall enabled"
					return m, m.LoadCmd()
				}
			}

		// Add a common safe firewall rule — allow SSH (so you don't lock yourself out)
		case "s":
			if m.activeTab == TabFirewall {
				result := core.RunCommand("ufw allow ssh 2>&1")
				if result.Err != nil {
					m.message = "Failed: " + result.Output
				} else {
					m.message = "SSH access allowed (port 22)"
					return m, m.LoadCmd()
				}
			}

		// Allow HTTP (port 80)
		case "w":
			if m.activeTab == TabFirewall {
				result := core.RunCommand("ufw allow http 2>&1")
				if result.Err != nil {
					m.message = "Failed: " + result.Output
				} else {
					m.message = "HTTP allowed (port 80)"
					return m, m.LoadCmd()
				}
			}

		// Allow HTTPS (port 443)
		case "W":
			if m.activeTab == TabFirewall {
				result := core.RunCommand("ufw allow https 2>&1")
				if result.Err != nil {
					m.message = "Failed: " + result.Output
				} else {
					m.message = "HTTPS allowed (port 443)"
					return m, m.LoadCmd()
				}
			}

		// Delete a firewall rule
		case "D":
			if m.activeTab == TabFirewall && len(m.data.firewallRules) > 0 {
				rule := m.data.firewallRules[m.cursor]
				m.message = fmt.Sprintf("Press X to delete rule %s: %s %s from %s",
					rule.Number, rule.To, rule.Action, rule.From)
			}

		// Confirm firewall rule delete
		case "X":
			if m.activeTab == TabFirewall && len(m.data.firewallRules) > 0 {
				rule := m.data.firewallRules[m.cursor]
				result := core.RunCommand(fmt.Sprintf("ufw --force delete %s 2>&1", rule.Number))
				if result.Err != nil {
					m.message = "Delete failed: " + result.Output
				} else {
					m.message = "Rule " + rule.Number + " deleted"
					return m, m.LoadCmd()
				}
			}

		// Unban all IPs in selected fail2ban jail
		case "u":
			if m.activeTab == TabFail2ban && len(m.data.fail2banJails) > 0 {
				jail := m.data.fail2banJails[m.cursor]
				result := core.RunCommand(fmt.Sprintf("fail2ban-client set %s unbanip all 2>&1", jail.Name))
				if result.Err != nil {
					m.message = "Unban failed: " + result.Output
				} else {
					m.message = "Unbanned all IPs in " + jail.Name
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

// currentListLen returns how many items are in the current tab's list.
// Used to clamp the cursor correctly when switching tabs.
func (m Model) currentListLen() int {
	switch m.activeTab {
	case TabFirewall:
		return len(m.data.firewallRules)
	case TabSSH:
		return len(m.data.sshKeys)
	case TabFail2ban:
		return len(m.data.fail2banJails)
	case TabUsers:
		return len(m.data.users)
	}
	return 0
}

// View renders the security screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Loading security info...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Security") + "\n\n")

	// Tab bar
	tabNames := []string{"Firewall", "SSH Keys", "Fail2ban", "Users"}
	tabLine := "  "
	for i, name := range tabNames {
		if Tab(i) == m.activeTab {
			tabLine += style.SidebarItemActiveStyle.Padding(0, 2).Render(name)
		} else {
			tabLine += style.BoxStyle.Padding(0, 2).Render(style.MutedStyle.Render(name))
		}
		if i < len(tabNames)-1 {
			tabLine += " "
		}
	}
	b.WriteString(tabLine + "\n\n")

	// Render the active tab content
	switch m.activeTab {
	case TabFirewall:
		b.WriteString(m.viewFirewall(width))
	case TabSSH:
		b.WriteString(m.viewSSH(width))
	case TabFail2ban:
		b.WriteString(m.viewFail2ban(width))
	case TabUsers:
		b.WriteString(m.viewUsers(width))
	}

	// Status message or help
	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render(m.helpText()))
	}

	return b.String()
}

// viewFirewall renders the UFW firewall rules list.
func (m Model) viewFirewall(width int) string {
	var b strings.Builder

	// Firewall status
	if m.data.firewallActive {
		b.WriteString("  UFW: " + style.SuccessStyle.Render("● active") + "\n\n")
	} else {
		b.WriteString("  UFW: " + style.DangerStyle.Render("● inactive") +
			style.MutedStyle.Render("  (press e to enable)") + "\n\n")
	}

	if len(m.data.firewallRules) == 0 {
		b.WriteString(style.MutedStyle.Render("  No rules yet. Add some:\n"))
		b.WriteString(style.MutedStyle.Render("  s = allow SSH  w = allow HTTP  W = allow HTTPS\n"))
		return b.String()
	}

	b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-5s %-22s %-8s %s\n", "#", "Port/Service", "Action", "From")))
	b.WriteString(style.MutedStyle.Render("  "+strings.Repeat("─", width-6)) + "\n")

	for i, rule := range m.data.firewallRules {
		actionStyled := ""
		if strings.Contains(rule.Action, "ALLOW") {
			actionStyled = style.SuccessStyle.Render("ALLOW   ")
		} else {
			actionStyled = style.DangerStyle.Render("DENY    ")
		}
		line := fmt.Sprintf("  %-5s %-22s %s %s",
			rule.Number, truncate(rule.To, 21), actionStyled, rule.From)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

// viewSSH renders the SSH authorized keys list.
func (m Model) viewSSH(width int) string {
	var b strings.Builder
	b.WriteString(style.BoldStyle.Render("  Authorized keys (~/.ssh/authorized_keys)\n\n"))

	if len(m.data.sshKeys) == 0 {
		b.WriteString(style.MutedStyle.Render("  No SSH keys found.\n"))
		b.WriteString(style.MutedStyle.Render("  Add a key: echo 'ssh-ed25519 AAA... user@host' >> ~/.ssh/authorized_keys\n"))
		return b.String()
	}

	for i, key := range m.data.sshKeys {
		comment := key.Comment
		if comment == "" {
			comment = "(no comment)"
		}
		line := fmt.Sprintf("  %-16s %-25s %s",
			key.Type, truncate(comment, 24), key.Short)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

// viewFail2ban renders the Fail2ban jail list.
func (m Model) viewFail2ban(width int) string {
	var b strings.Builder

	if !core.CommandExists("fail2ban-client") {
		b.WriteString(style.MutedStyle.Render("  Fail2ban is not installed.\n"))
		b.WriteString("  Install with: " + style.BoldStyle.Render("apt install fail2ban") + "\n")
		return b.String()
	}

	if len(m.data.fail2banJails) == 0 {
		b.WriteString(style.MutedStyle.Render("  No active fail2ban jails.\n"))
		return b.String()
	}

	b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-25s %-10s %s\n", "Jail", "Status", "Banned IPs")))
	b.WriteString(style.MutedStyle.Render("  "+strings.Repeat("─", width-6)) + "\n")

	for i, jail := range m.data.fail2banJails {
		status := style.SuccessStyle.Render("enabled   ")
		if !jail.Enabled {
			status = style.MutedStyle.Render("disabled  ")
		}
		bannedStr := ""
		if jail.Banned > 0 {
			bannedStr = style.WarningStyle.Render(fmt.Sprintf("%d banned", jail.Banned))
		} else {
			bannedStr = style.MutedStyle.Render("none")
		}
		line := fmt.Sprintf("  %-25s %s %s", truncate(jail.Name, 24), status, bannedStr)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

// viewUsers renders the system user account list.
func (m Model) viewUsers(width int) string {
	var b strings.Builder
	b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-16s %-6s %-20s %s\n", "User", "UID", "Shell", "Home")))
	b.WriteString(style.MutedStyle.Render("  "+strings.Repeat("─", width-6)) + "\n")

	for i, user := range m.data.users {
		nameStyled := user.Name
		if user.Name == "root" {
			nameStyled = style.DangerStyle.Render(user.Name) // root in red as a reminder
		}
		line := fmt.Sprintf("  %-16s %-6s %-20s %s",
			nameStyled, user.UID, truncate(user.Shell, 19), user.Home)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

// helpText returns the relevant help line for the active tab.
func (m Model) helpText() string {
	base := "  tab: switch section  r: refresh"
	switch m.activeTab {
	case TabFirewall:
		return "  s:allow SSH  w:allow HTTP  W:allow HTTPS  D→X:delete rule  e:enable UFW  r:refresh"
	case TabSSH:
		return base
	case TabFail2ban:
		return "  u: unban all in selected jail  r: refresh"
	case TabUsers:
		return base
	}
	return base
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
