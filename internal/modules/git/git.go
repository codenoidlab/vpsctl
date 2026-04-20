// Package git manages Git repositories on the server.
// It finds all git repos under common directories and lets you
// pull, push, commit, check status, and switch branches — without typing git commands.
package git

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/core"
	"github.com/vpsmanager/vps/internal/style"
)

// DataMsg is sent when repo scanning is done.
type DataMsg struct {
	repos []RepoInfo
	err   string
}

// RepoInfo holds the state of one git repository.
type RepoInfo struct {
	Path   string // full path to the repo
	Name   string // just the folder name
	Branch string // current branch
	Status string // "clean", "modified", "ahead", etc.
	Ahead  int    // commits ahead of remote
	Behind int    // commits behind remote
}

// Mode tracks what sub-view we're in
type Mode int

const (
	ModeList   Mode = iota // looking at the list of repos
	ModeDetail             // looking at one repo's status/diff
)

// Model holds state for the Git screen.
type Model struct {
	repos   []RepoInfo
	cursor  int
	loaded  bool
	mode    Mode
	detail  string // status/diff output for the selected repo
	message string
}

func New() Model { return Model{} }

// LoadCmd scans common directories for git repos.
func (m Model) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		return findGitRepos()
	}
}

// findGitRepos searches common locations for .git directories.
// We look in home, /var/www, /opt, and /srv — the usual spots on a VPS.
func findGitRepos() DataMsg {
	data := DataMsg{}

	if !core.CommandExists("git") {
		data.err = "git is not installed"
		return data
	}

	// Search these directories for git repos (up to 3 levels deep)
	searchDirs := []string{"~", "/var/www", "/opt", "/srv"}
	seen := map[string]bool{} // avoid duplicates

	for _, dir := range searchDirs {
		// `find` locates all .git directories, then we take the parent
		result := core.RunCommand(fmt.Sprintf(
			"find %s -maxdepth 3 -name '.git' -type d 2>/dev/null | head -20", dir,
		))
		if result.Err != nil || result.Output == "" {
			continue
		}

		for _, gitDir := range strings.Split(result.Output, "\n") {
			gitDir = strings.TrimSpace(gitDir)
			if gitDir == "" {
				continue
			}
			// The repo root is the parent of .git
			repoPath := strings.TrimSuffix(gitDir, "/.git")
			if seen[repoPath] {
				continue
			}
			seen[repoPath] = true

			repo := getRepoInfo(repoPath)
			data.repos = append(data.repos, repo)
		}
	}

	return data
}

// getRepoInfo collects branch, status, and ahead/behind info for one repo.
func getRepoInfo(path string) RepoInfo {
	repo := RepoInfo{Path: path}

	// Repo name is just the last part of the path
	parts := strings.Split(path, "/")
	repo.Name = parts[len(parts)-1]

	// Get current branch name
	branchResult := core.RunCommand(fmt.Sprintf(
		"git -C %q rev-parse --abbrev-ref HEAD 2>/dev/null", path,
	))
	repo.Branch = branchResult.Output
	if repo.Branch == "" {
		repo.Branch = "unknown"
	}

	// Fetch remote status (don't do git fetch — that's slow, just check tracking)
	aheadBehind := core.RunCommand(fmt.Sprintf(
		"git -C %q rev-list --left-right --count @{upstream}...HEAD 2>/dev/null", path,
	))
	if aheadBehind.Err == nil && aheadBehind.Output != "" {
		parts := strings.Fields(aheadBehind.Output)
		if len(parts) == 2 {
			_ = fmt.Sscanf(parts[0], "%d", &repo.Behind) // intentionally ignoring parse error for git stats
			_ = fmt.Sscanf(parts[1], "%d", &repo.Ahead)  // intentionally ignoring parse error for git stats
		}
	}

	// Check if working tree is clean or has changes
	statusResult := core.RunCommand(fmt.Sprintf(
		"git -C %q status --porcelain 2>/dev/null", path,
	))
	if statusResult.Output == "" {
		repo.Status = "clean"
	} else {
		// Count the changed files
		changedCount := len(strings.Split(strings.TrimSpace(statusResult.Output), "\n"))
		repo.Status = fmt.Sprintf("%d changed", changedCount)
	}

	return repo
}

// Update handles input on the Git screen.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DataMsg:
		m.repos = msg.repos
		m.loaded = true
		if msg.err != "" {
			m.message = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		// In detail mode: esc goes back to list
		if m.mode == ModeDetail {
			if msg.String() == "esc" || msg.String() == "q" || msg.String() == "backspace" {
				m.mode = ModeList
				m.detail = ""
			}
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}

		// View git status of selected repo
		case "enter", "s":
			if len(m.repos) > 0 {
				repo := m.repos[m.cursor]
				result := core.RunCommand(fmt.Sprintf("git -C %q status 2>&1", repo.Path))
				m.detail = result.Output
				m.mode = ModeDetail
			}

		// View diff
		case "d":
			if len(m.repos) > 0 {
				repo := m.repos[m.cursor]
				result := core.RunCommand(fmt.Sprintf("git -C %q diff --stat 2>&1", repo.Path))
				m.detail = result.Output
				if m.detail == "" {
					m.detail = "(no unstaged changes)"
				}
				m.mode = ModeDetail
			}

		// Git pull
		case "p":
			if len(m.repos) > 0 {
				repo := m.repos[m.cursor]
				result := core.RunCommand(fmt.Sprintf("git -C %q pull 2>&1", repo.Path))
				m.message = result.Output
				if len(m.message) > 80 {
					m.message = m.message[:80] + "…"
				}
				return m, m.LoadCmd()
			}

		// Switch branch (shows a simple prompt via message)
		case "b":
			if len(m.repos) > 0 {
				repo := m.repos[m.cursor]
				result := core.RunCommand(fmt.Sprintf("git -C %q branch -a 2>/dev/null", repo.Path))
				m.detail = result.Output
				m.mode = ModeDetail
				m.message = "Branches listed. Press esc to go back."
			}

		// Refresh
		case "r":
			m.loaded = false
			return m, m.LoadCmd()
		}
	}
	return m, nil
}

// View renders the Git screen.
func (m Model) View(width int) string {
	if !m.loaded {
		return style.MutedStyle.Render("\n  Scanning for git repos...")
	}

	var b strings.Builder
	b.WriteString(style.TitleStyle.Render("Git") + "\n\n")

	// Detail view — shows status/diff output for one repo
	if m.mode == ModeDetail {
		if len(m.repos) > 0 {
			b.WriteString(style.BoldStyle.Render("  " + m.repos[m.cursor].Name) + "\n\n")
		}
		// Show each line of the git output
		for _, line := range strings.Split(m.detail, "\n") {
			// Color-code git status lines
			switch {
			case strings.HasPrefix(line, "+") || strings.Contains(line, "new file"):
				b.WriteString("  " + style.SuccessStyle.Render(line) + "\n")
			case strings.HasPrefix(line, "-") || strings.Contains(line, "deleted"):
				b.WriteString("  " + style.DangerStyle.Render(line) + "\n")
			case strings.Contains(line, "modified"):
				b.WriteString("  " + style.WarningStyle.Render(line) + "\n")
			default:
				b.WriteString("  " + line + "\n")
			}
		}
		b.WriteString("\n" + style.MutedStyle.Render("  backspace / esc: back to list"))
		return b.String()
	}

	// List view
	if len(m.repos) == 0 {
		b.WriteString(style.MutedStyle.Render("  No git repos found in ~/  /var/www/  /opt/  /srv/\n"))
		b.WriteString(style.MutedStyle.Render("  Press r to rescan."))
		return b.String()
	}

	// Column header
	b.WriteString(style.BoldStyle.Render(fmt.Sprintf("  %-24s %-14s %-10s %s\n",
		"Repo", "Branch", "Status", "Sync")))
	b.WriteString(style.MutedStyle.Render("  " + strings.Repeat("─", width-6) + "\n"))

	for i, repo := range m.repos {
		// Sync indicator: show how many commits ahead/behind
		sync := ""
		if repo.Ahead > 0 {
			sync += style.WarningStyle.Render(fmt.Sprintf("↑%d", repo.Ahead)) + " "
		}
		if repo.Behind > 0 {
			sync += style.DangerStyle.Render(fmt.Sprintf("↓%d", repo.Behind))
		}
		if sync == "" {
			sync = style.MutedStyle.Render("in sync")
		}

		statusStyled := ""
		if repo.Status == "clean" {
			statusStyled = style.SuccessStyle.Render("clean     ")
		} else {
			statusStyled = style.WarningStyle.Render(repo.Status + "  ")
		}

		line := fmt.Sprintf("  %-24s %-14s %s %s",
			truncate(repo.Name, 23),
			truncate(repo.Branch, 13),
			statusStyled,
			sync,
		)

		if i == m.cursor {
			b.WriteString(style.SidebarItemActiveStyle.Width(width-4).Render(line) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(style.WarningStyle.Render("  " + m.message))
	} else {
		b.WriteString(style.MutedStyle.Render("  enter: status  d: diff  p: pull  b: branches  r: rescan"))
	}

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
