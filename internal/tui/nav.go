// nav.go defines the sidebar navigation — the list of screens on the left.
// Each screen has a name, an icon, and a keyboard shortcut.
package tui

// Screen is just a number that represents which page we're on.
// We use numbers (not strings) so switching screens is a simple comparison.
type Screen int

const (
	ScreenDashboard Screen = iota // 0 — home screen
	ScreenFiles                   // 1 — file manager
	ScreenNode                    // 2 — PM2 / Node apps
	ScreenNginx                   // 3 — Nginx + Cloudflare
	ScreenGit                     // 4 — Git repos
	ScreenPackages                // 5 — apt/npm packages
	ScreenMonitor                 // 6 — system monitor
	ScreenSecurity                // 7 — firewall/SSH
)

// NavItem is one entry in the sidebar menu
type NavItem struct {
	Icon     string // emoji icon shown next to the name
	Label    string // name shown in the sidebar
	Shortcut string // keyboard shortcut shown in help bar
	Target   Screen // which screen this item goes to
}

// NavItems is the full list of sidebar menu entries, in display order.
// To add a new screen: add a Screen constant above AND an entry here.
var NavItems = []NavItem{
	{"󰕮", "Dashboard", "1", ScreenDashboard},
	{"", "Files", "2", ScreenFiles},
	{"", "Node / PM2", "3", ScreenNode},
	{"", "Nginx", "4", ScreenNginx},
	{"", "Git", "5", ScreenGit},
	{"", "Packages", "6", ScreenPackages},
	{"", "Monitor", "7", ScreenMonitor},
	{"", "Security", "8", ScreenSecurity},
}

// renderSidebar draws the left sidebar.
// activeScreen is which screen is currently showing so we can highlight it.
func renderSidebar(activeScreen Screen, height int) string {
	// App name at the top of the sidebar
	header := TitleStyle.Render("  VPS Manager")

	items := ""
	for _, item := range NavItems {
		// Build the text: icon + label
		line := item.Icon + "  " + item.Label

		if item.Target == activeScreen {
			// This is the currently active screen — highlight it
			items += SidebarItemActiveStyle.Width(20).Render(line) + "\n"
		} else {
			items += SidebarItemStyle.Width(20).Render(line) + "\n"
		}
	}

	// Put header and items together inside the sidebar container
	_ = height // we could use this to pad the sidebar to full height later
	return SidebarStyle.Width(22).Render(header + "\n\n" + items)
}
