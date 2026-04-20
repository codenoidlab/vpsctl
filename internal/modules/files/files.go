// Package files is the file manager screen.
// It shows a two-panel view: left = directory tree, right = file list.
// You navigate with arrow keys and can copy, move, rename, and delete files.
package files

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// DirLoadedMsg is sent when we finish reading a directory's contents.
type DirLoadedMsg struct {
	path    string
	entries []FileEntry
	err     error
}

// FileEntry is one file or folder in the current directory.
type FileEntry struct {
	Name  string
	IsDir bool
	Size  string // human-readable size like "4.2K"
	Perms string // like "rwxr-xr-x"
}

// Model holds the state for the file manager.
type Model struct {
	currentPath string      // the directory we're looking at
	entries     []FileEntry // files in that directory
	cursor      int         // which entry is highlighted
	loaded      bool
	message     string // status message shown at the bottom
}

// New creates a fresh file manager starting at the home directory.
func New() Model {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}
	return Model{
		currentPath: home,
		cursor:      0,
	}
}

// LoadCmd returns a command that reads the current directory.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return loadDir(m.currentPath)
	}
}

// loadDir reads all files in a directory and returns them as FileEntry list.
func loadDir(path string) DirLoadedMsg {
	entries, err := os.ReadDir(path)
	if err != nil {
		return DirLoadedMsg{path: path, err: err}
	}

	var files []FileEntry

	// Always add ".." so you can navigate up
	files = append(files, FileEntry{Name: "..", IsDir: true})

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		size := ""
		if !entry.IsDir() {
			size = formatSize(info.Size())
		}

		files = append(files, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  size,
			Perms: info.Mode().String(),
		})
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].Name == ".." {
			return true // ".." always first
		}
		if files[j].Name == ".." {
			return false
		}
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir // dirs before files
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return DirLoadedMsg{path: path, entries: files}
}

// formatSize converts bytes to a human-readable string like "4.2K" or "1.3M".
func formatSize(bytes int64) string {
	switch {
	case bytes >= 1024*1024*1024:
		return fmt.Sprintf("%.1fG", float64(bytes)/(1024*1024*1024))
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1fM", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1fK", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// Update handles keyboard input for the file manager.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Directory finished loading
	case DirLoadedMsg:
		if msg.err != nil {
			m.message = "Error: " + msg.err.Error()
			return m, nil
		}
		m.currentPath = msg.path
		m.entries = msg.entries
		m.cursor = 0
		m.loaded = true
		m.message = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		// Move cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Move cursor down
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		// Enter a directory or open a file
		case "enter":
			if len(m.entries) == 0 {
				break
			}
			selected := m.entries[m.cursor]
			if selected.IsDir {
				// Navigate into the directory
				newPath := filepath.Join(m.currentPath, selected.Name)
				m.loaded = false
				return m, func() tea.Msg { return loadDir(newPath) }
			}
			// Open file in less (read-only viewer)
			// In a real terminal app we'd launch an external command here
			m.message = "Use 'e' to edit, currently view-only for files"

		// Go up one directory (same as selecting "..")
		case "backspace", "h":
			parent := filepath.Dir(m.currentPath)
			if parent != m.currentPath { // stop at filesystem root
				m.loaded = false
				return m, func() tea.Msg { return loadDir(parent) }
			}

		// Copy a file
		case "c":
			if len(m.entries) > 0 && m.entries[m.cursor].Name != ".." {
				src := filepath.Join(m.currentPath, m.entries[m.cursor].Name)
				dest := src + ".copy"
				result := core.RunCommand(fmt.Sprintf("cp -r %q %q", src, dest))
				if result.Err != nil {
					m.message = "Copy failed: " + result.Output
				} else {
					m.message = "Copied to: " + dest
					return m, m.LoadCmd() // refresh directory listing
				}
			}

		// Delete a file (asks for confirmation via message)
		case "d":
			if len(m.entries) > 0 && m.entries[m.cursor].Name != ".." {
				name := m.entries[m.cursor].Name
				m.message = fmt.Sprintf("Press D to confirm delete: %s", name)
			}

		// Confirm delete
		case "D":
			if len(m.entries) > 0 && m.entries[m.cursor].Name != ".." {
				target := filepath.Join(m.currentPath, m.entries[m.cursor].Name)
				result := core.RunCommand(fmt.Sprintf("rm -rf %q", target))
				if result.Err != nil {
					m.message = "Delete failed: " + result.Output
				} else {
					m.message = "Deleted: " + target
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

// View renders the file manager screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Loading files...")
	}

	var b strings.Builder

	// Screen title with current path
	b.WriteString(style.TitleStyle.Render("Files") + "\n")
	b.WriteString(style.MutedStyle.Render("  " + m.currentPath) + "\n\n")

	// File list
	// Visible area — show a window of entries around the cursor
	maxVisible := 20
	start := 0
	if m.cursor > maxVisible-3 {
		start = m.cursor - maxVisible + 3
	}

	for i := start; i < len(m.entries) && i < start+maxVisible; i++ {
		entry := m.entries[i]

		// Build the icon + name display
		icon := "  " // file icon
		if entry.IsDir {
			icon = " " // folder icon
		}
		if entry.Name == ".." {
			icon = "  "
		}

		name := entry.Name
		if entry.IsDir && entry.Name != ".." {
			name = style.BoldStyle.Render(name + "/")
		}

		// Size column, right-aligned
		size := fmt.Sprintf("%6s", entry.Size)

		line := fmt.Sprintf("%s %s", icon, name)

		if i == m.cursor {
			// Highlighted line
			b.WriteString(style.SidebarItemActiveStyle.Width(width - 12).Render(line))
			b.WriteString(style.MutedStyle.Render(size))
		} else {
			b.WriteString("  " + line)
			if entry.Size != "" {
				b.WriteString(style.MutedStyle.Render(fmt.Sprintf("%*s", width-len(line)-10, size)))
			}
		}
		b.WriteString("\n")
	}

	// Show how many more entries there are if the list is scrolled
	if len(m.entries) > maxVisible {
		b.WriteString(style.MutedStyle.Render(fmt.Sprintf("\n  ... %d total entries", len(m.entries))))
	}

	// Status message or help text at the bottom
	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  ↑↓: navigate  enter: open  backspace: up  c: copy  d: delete  r: refresh"))
	}

	return b.String()
}
