// Package packages is the package installer screen.
// It lets you search and install apt packages (system), global npm packages (Node),
// and snap packages — all without typing commands.
// Tab switches between the three package managers.
package packages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// Which package manager is active
type Tab int

const (
	TabApt  Tab = iota // apt — system packages
	TabNpm             // npm global packages
	TabSnap            // snap packages
)

// DataMsg carries the loaded package list.
type DataMsg struct {
	tab      Tab
	packages []PackageInfo
	err      string
}

// SearchResultMsg carries results from a package search.
type SearchResultMsg struct {
	results []string
	err     string
}

// PackageInfo is one installed package.
type PackageInfo struct {
	Name    string
	Version string
	Source  string // "apt", "npm", "snap"
}

// Mode tracks what sub-view we're in.
type Mode int

const (
	ModeList    Mode = iota // browsing installed packages
	ModeSearch              // typing a search query
	ModeResults             // showing search results
)

// Model holds the packages screen state.
type Model struct {
	activeTab   Tab
	packages    []PackageInfo  // installed packages for current tab
	cursor      int
	loaded      bool
	mode        Mode
	searchQuery string   // what the user typed in search
	searchResults []string
	message     string
}

func New() Model { return Model{activeTab: TabApt} }

// LoadCmd loads installed packages for the currently active tab.
func (m Model) LoadCmd() tea.Cmd {
	tab := m.activeTab
	return func() tea.Msg {
		return fetchPackages(tab)
	}
}

// fetchPackages gets the list of installed packages for the given tab.
func fetchPackages(tab Tab) DataMsg {
	data := DataMsg{tab: tab}

	switch tab {
	case TabApt:
		// List manually installed apt packages (not auto-dependencies)
		// dpkg-query gives us name and version cleanly
		result := core.RunCommand(
			"dpkg-query -W -f='${Package}\\t${Version}\\n' 2>/dev/null | head -100",
		)
		if result.Err != nil {
			data.err = result.Output
			return data
		}
		for _, line := range strings.Split(result.Output, "\n") {
			parts := strings.Split(line, "\t")
			if len(parts) == 2 && parts[0] != "" {
				data.packages = append(data.packages, PackageInfo{
					Name:    parts[0],
					Version: parts[1],
					Source:  "apt",
				})
			}
		}

	case TabNpm:
		// List globally installed npm packages
		result := core.RunCommand("npm list -g --depth=0 2>/dev/null")
		if result.Err != nil && result.Output == "" {
			data.err = "npm not installed"
			return data
		}
		// Output looks like: "├── pm2@5.3.0" or "└── typescript@5.0.0"
		for _, line := range strings.Split(result.Output, "\n") {
			// Strip tree characters
			clean := strings.TrimLeft(line, "├└─│ ")
			if clean == "" || strings.HasPrefix(clean, "/") {
				continue // skip the path line at the top
			}
			// Split "name@version"
			atIdx := strings.LastIndex(clean, "@")
			if atIdx <= 0 {
				continue
			}
			data.packages = append(data.packages, PackageInfo{
				Name:    clean[:atIdx],
				Version: clean[atIdx+1:],
				Source:  "npm",
			})
		}

	case TabSnap:
		// List installed snap packages
		result := core.RunCommand("snap list 2>/dev/null")
		if result.Err != nil {
			data.err = "snap not installed or no snaps"
			return data
		}
		// Snap output: "Name     Version  Rev  Tracking  Publisher  Notes"
		lines := strings.Split(result.Output, "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue // skip header
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				data.packages = append(data.packages, PackageInfo{
					Name:    fields[0],
					Version: fields[1],
					Source:  "snap",
				})
			}
		}
	}

	return data
}

// Update handles input on the packages screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case DataMsg:
		// Only apply if this message is for our current tab
		// (prevents stale data from a previous tab showing up)
		if msg.tab == m.activeTab {
			m.packages = msg.packages
			m.loaded = true
			m.cursor = 0
			if msg.err != "" {
				m.message = msg.err
			}
		}
		return m, nil

	case SearchResultMsg:
		m.searchResults = msg.results
		m.mode = ModeResults
		m.cursor = 0
		if msg.err != "" {
			m.message = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		// In search mode — collect typed characters
		if m.mode == ModeSearch {
			switch msg.String() {
			case "esc":
				m.mode = ModeList
				m.searchQuery = ""
			case "enter":
				// Run the search
				query := m.searchQuery
				tab := m.activeTab
				return m, func() tea.Msg {
					return runSearch(query, tab)
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
			default:
				// Only add printable single characters
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}
			}
			return m, nil
		}

		// In results mode
		if m.mode == ModeResults {
			switch msg.String() {
			case "esc", "backspace":
				m.mode = ModeList
				m.searchResults = nil
				m.searchQuery = ""
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.searchResults)-1 {
					m.cursor++
				}
			case "enter", "i":
				// Install the selected search result
				if len(m.searchResults) > 0 {
					pkgName := strings.Fields(m.searchResults[m.cursor])[0]
					tab := m.activeTab
					m.message = "Installing " + pkgName + "..."
					m.mode = ModeList
					return m, func() tea.Msg {
						return installPackage(pkgName, tab)
					}
				}
			}
			return m, nil
		}

		// Normal list mode
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.packages)-1 {
				m.cursor++
			}

		// Switch tabs with left/right or number keys
		case "left", "h":
			if m.activeTab > TabApt {
				m.activeTab--
				m.loaded = false
				m.cursor = 0
				return m, m.LoadCmd()
			}
		case "right", "l":
			if m.activeTab < TabSnap {
				m.activeTab++
				m.loaded = false
				m.cursor = 0
				return m, m.LoadCmd()
			}

		// Tab key also cycles between apt/npm/snap
		case "tab":
			m.activeTab = Tab((int(m.activeTab) + 1) % 3)
			m.loaded = false
			m.cursor = 0
			return m, m.LoadCmd()

		// Search for a package to install
		case "/":
			m.mode = ModeSearch
			m.searchQuery = ""

		// Remove selected package
		case "D":
			if len(m.packages) > 0 {
				pkg := m.packages[m.cursor]
				m.message = fmt.Sprintf("Press X to confirm remove: %s", pkg.Name)
			}

		// Confirm remove
		case "X":
			if len(m.packages) > 0 {
				pkg := m.packages[m.cursor]
				tab := m.activeTab
				m.message = "Removing " + pkg.Name + "..."
				return m, func() tea.Msg {
					return removePackage(pkg.Name, tab)
				}
			}

		// Update all packages in the current tab
		case "u":
			tab := m.activeTab
			m.message = "Updating packages..."
			return m, func() tea.Msg {
				return updatePackages(tab)
			}

		// Refresh the list
		case "r":
			m.loaded = false
			return m, m.LoadCmd()
		}
	}
	return m, nil
}

// runSearch searches for packages matching the query in the active tab's manager.
func runSearch(query string, tab Tab) SearchResultMsg {
	var result core.CommandResult
	switch tab {
	case TabApt:
		result = core.RunCommand(fmt.Sprintf("apt-cache search %q 2>/dev/null | head -30", query))
	case TabNpm:
		result = core.RunCommand(fmt.Sprintf("npm search %s --no-description 2>/dev/null | head -20", query))
	case TabSnap:
		result = core.RunCommand(fmt.Sprintf("snap find %q 2>/dev/null | head -20", query))
	}
	if result.Err != nil {
		return SearchResultMsg{err: result.Output}
	}
	lines := strings.Split(result.Output, "\n")
	var filtered []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			filtered = append(filtered, l)
		}
	}
	return SearchResultMsg{results: filtered}
}

// installPackage installs a package using the right manager for the active tab.
func installPackage(name string, tab Tab) DataMsg {
	var result core.CommandResult
	switch tab {
	case TabApt:
		result = core.RunCommand(fmt.Sprintf("apt-get install -y %s 2>&1", name))
	case TabNpm:
		result = core.RunCommand(fmt.Sprintf("npm install -g %s 2>&1", name))
	case TabSnap:
		result = core.RunCommand(fmt.Sprintf("snap install %s 2>&1", name))
	}
	if result.Err != nil {
		return DataMsg{tab: tab, err: "Install failed: " + truncate(result.Output, 80)}
	}
	// Reload the package list
	return fetchPackages(tab)
}

// removePackage removes a package using the right manager.
func removePackage(name string, tab Tab) DataMsg {
	var result core.CommandResult
	switch tab {
	case TabApt:
		result = core.RunCommand(fmt.Sprintf("apt-get remove -y %s 2>&1", name))
	case TabNpm:
		result = core.RunCommand(fmt.Sprintf("npm uninstall -g %s 2>&1", name))
	case TabSnap:
		result = core.RunCommand(fmt.Sprintf("snap remove %s 2>&1", name))
	}
	if result.Err != nil {
		return DataMsg{tab: tab, err: "Remove failed: " + truncate(result.Output, 80)}
	}
	return fetchPackages(tab)
}

// updatePackages updates all packages in the current tab.
func updatePackages(tab Tab) DataMsg {
	var result core.CommandResult
	switch tab {
	case TabApt:
		// Update package index first, then upgrade
		core.RunCommand("apt-get update 2>&1")
		result = core.RunCommand("apt-get upgrade -y 2>&1 | tail -5")
	case TabNpm:
		result = core.RunCommand("npm update -g 2>&1 | tail -5")
	case TabSnap:
		result = core.RunCommand("snap refresh 2>&1 | tail -5")
	}
	_ = result
	return fetchPackages(tab)
}

// View renders the packages screen.
func (m Model) View(width int) string {
	if !m.loaded {
		tabs := []string{"apt", "npm", "snap"}
		return style.MutedStyle.Render("\n  Loading " + tabs[m.activeTab] + " packages...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Packages") + "\n\n")

	// Tab bar — shows apt | npm | snap with active one highlighted
	tabs := []string{"apt (system)", "npm (global)", "snap"}
	tabLine := "  "
	for i, name := range tabs {
		if Tab(i) == m.activeTab {
			tabLine += style.SidebarItemActiveStyle.Padding(0, 2).Render(name)
		} else {
			tabLine += style.BoxStyle.Padding(0, 2).Render(style.MutedStyle.Render(name))
		}
		if i < len(tabs)-1 {
			tabLine += " "
		}
	}
	b.WriteString(tabLine + "\n\n")

	// Search mode — show the search box
	if m.mode == ModeSearch {
		b.WriteString(style.BoldStyle.Render("  Search: ") + m.searchQuery + "█\n\n")
		b.WriteString(style.MutedStyle.Render("  Type to search, enter to run, esc to cancel"))
		return b.String()
	}

	// Results mode — show search results
	if m.mode == ModeResults {
		b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  Search results for: %s\n\n", m.searchQuery)))
		if len(m.searchResults) == 0 {
			b.WriteString(style.MutedStyle.Render("  No results found\n"))
		}
		for i, result := range m.searchResults {
			line := "  " + truncate(result, width-6)
			if i == m.cursor {
				b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString("\n" + style.MutedStyle.Render("  enter/i: install selected  esc: back"))
		return b.String()
	}

	// Package list
	b.WriteString(style.BoldStyle.Render(fmt.Sprintf(
		"  %-35s %s\n", "Package", "Version",
	)))
	b.WriteString(style.MutedStyle.Render("  "+strings.Repeat("─", width-6)) + "\n")

	if len(m.packages) == 0 {
		b.WriteString(style.MutedStyle.Render("\n  No packages found. Press / to search and install one.\n"))
	}

	// Show a window of packages around the cursor
	maxVisible := 18
	start := 0
	if m.cursor > maxVisible-3 {
		start = m.cursor - maxVisible + 3
	}
	for i := start; i < len(m.packages) && i < start+maxVisible; i++ {
		pkg := m.packages[i]
		line := fmt.Sprintf("  %-35s %s", truncate(pkg.Name, 34), pkg.Version)
		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}

	if len(m.packages) > maxVisible {
		b.WriteString(style.MutedStyle.Render(fmt.Sprintf(
			"\n  … %d total packages (scrolling)", len(m.packages),
		)))
	}

	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  /: search+install  D→X: remove  u: update all  tab: switch manager  r: refresh"))
	}

	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
