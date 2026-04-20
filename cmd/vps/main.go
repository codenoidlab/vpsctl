// main.go is the entry point for the VPS Manager.
// It just boots BubbleTea with our app model and hands off control.
// Keep this file tiny — all real logic lives in internal/tui and internal/modules.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpsmanager/vps/internal/tui"
)

func main() {
	// Print version if asked
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("VPS Manager v0.1.0")
		os.Exit(0)
	}

	// Create the root app model
	app := tui.NewApp()

	// Start BubbleTea.
	// WithAltScreen() uses the terminal's alternate screen buffer —
	// this means the TUI fills the whole terminal and disappears cleanly on exit.
	program := tea.NewProgram(app, tea.WithAltScreen())

	// Run until the user quits
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running VPS Manager:", err)
		os.Exit(1)
	}
}
