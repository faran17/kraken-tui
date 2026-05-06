package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faran17/kraken-tui/internal/app"
)

func main() {
	// Redirect debug logs so they don't corrupt the TUI
	f, err := tea.LogToFile("debug.log", "kraken")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: cannot open debug log: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	apiKey := os.Getenv("GEMINI_API_KEY")

	m, err := app.New(apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
