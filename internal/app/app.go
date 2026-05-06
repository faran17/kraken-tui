// Package app contains the root Bubble Tea model that composes the file browser,
// Gemini chat, and todo list panels into a single cohesive interface.
package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/faran17/kraken-tui/internal/chat"
	"github.com/faran17/kraken-tui/internal/filebrowser"
	"github.com/faran17/kraken-tui/internal/terminal"
	"github.com/faran17/kraken-tui/internal/todo"
	"github.com/faran17/kraken-tui/pkg/styles"
)

var krakenLogo = styles.Logo.Render(strings.TrimPrefix(`
    __ __                 __                
   / //_/  ____   ____   / /_   ___   ____  
  / ,<    / __ \ / __ \ / ,<   / _ \ / __ \ 
 / /| |  / /  / / /_/ // /| | /  __// / / / 
/_/ |_| /_/    \__,_/ /_/ |_| \___//_/ /_/  
`, "\n"))

// ── Constants ─────────────────────────────────────────────────────────────────

// Panel indices for routing keyboard events and drawing active borders.
const (
	PanelFiles    = iota // index 0: File Browser
	PanelChat            // index 1: Gemini Chat
	PanelTodo            // index 2: Todo List
	PanelTerminal        // index 3: Terminal
	panelCount           // total number of panels (4)
)

// ── Model ─────────────────────────────────────────────────────────────────────

// Model is the root (compositor) Bubble Tea model. It owns the three child panels
// and handles global hotkeys (like Tab and Ctrl+C).
type Model struct {
	width  int // terminal width
	height int // terminal height

	activePanel int // which panel currently receives key events

	// The child components
	files filebrowser.Model
	chat  chat.Model
	todo  todo.Model
	term  terminal.Model

	// Global UI elements
	spinner spinner.Model // loading spinner used before first render
	ready   bool          // becomes true after receiving the first WindowSizeMsg
}

// New constructs the root model. It initializes the three child panels and
// passes the API key through to the chat model.
func New(apiKey string) (Model, error) {
	// A generic loading spinner used if startup takes a moment.
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.ChatSpinner

	// Initialize the chat panel. It handles its own persistence and API setup.
	chatModel, err := chat.New(apiKey)
	if err != nil {
		return Model{}, fmt.Errorf("chat init: %w", err)
	}

	// Initialize the todo panel. It loads items from disk immediately.
	todoModel, err := todo.New()
	if err != nil {
		return Model{}, fmt.Errorf("todo init: %w", err)
	}

	// Initialize the terminal panel.
	termModel := terminal.New()

	return Model{
		activePanel: PanelFiles, // start with the file browser focused
		files:       filebrowser.New(),
		chat:        chatModel,
		todo:        todoModel,
		term:        termModel,
		spinner:     s,
	}, nil
}

// ── Tea interface ─────────────────────────────────────────────────────────────

// Init returns a batch of all startup commands needed by the child panels.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick, // start the global loading spinner
		m.chat.Init(),  // start cursor blink in chat input
		m.todo.Init(),  // (currently no-op)
		m.files.Init(), // (currently no-op)
		m.term.Init(),  // start terminal
	)
}

// Update handles global key presses and routes all other messages to the
// currently active child panel.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// ── Keyboard events ──────────────────────────────────────────────────────
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			// Global quit sequence
			return m, tea.Quit

		case "tab":
			// Cycle active panel forwards: 0 -> 1 -> 2 -> 0
			m.activePanel = (m.activePanel + 1) % panelCount
			return m, nil

		case "shift+tab":
			// Cycle active panel backwards: 0 -> 2 -> 1 -> 0
			// Adding panelCount before modulo prevents negative results in Go.
			m.activePanel = (m.activePanel - 1 + panelCount) % panelCount
			return m, nil
		}

	// ── Terminal resize events ───────────────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Calculate available dimensions for the child panels
		w, h := m.panelDimensions()

		// Propagate resize downwards
		m.files = m.files.SetSize(w, h)
		m.chat = m.chat.SetSize(w, h)
		m.todo = m.todo.SetSize(w, h)
		m.term = m.term.SetSize(m.width, m.terminalHeight())
		return m, nil

	// ── Spinner ticks ────────────────────────────────────────────────────────
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// ── Routing ──────────────────────────────────────────────────────────────
	// If the message wasn't consumed globally (like Tab or Ctrl+C), route it
	// to the active child panel.
	switch m.activePanel {
	case PanelFiles:
		newFiles, cmd := m.files.Update(msg)
		m.files = newFiles
		cmds = append(cmds, cmd)

	case PanelChat:
		newChat, cmd := m.chat.Update(msg)
		m.chat = newChat
		cmds = append(cmds, cmd)

	case PanelTodo:
		newTodo, cmd := m.todo.Update(msg)
		m.todo = newTodo
		cmds = append(cmds, cmd)

	case PanelTerminal:
		newTerm, cmd := m.term.Update(msg)
		m.term = newTerm
		cmds = append(cmds, cmd)
	}

	// Because child updates might return cmds (like file browser search keystrokes
	// or chat token streaming), we return a batch of all accumulated commands.
	return m, tea.Batch(cmds...)
}

// View constructs the final string drawn to the terminal.
func (m Model) View() string {
	if !m.ready {
		// Terminal hasn't provided a WindowSizeMsg yet; show loading.
		return "\n  " + m.spinner.View() + "  Loading Kraken…"
	}

	w, h := m.panelDimensions()

	// Render each panel as a distinct block of text.
	filesPanel := m.renderPanel(PanelFiles, "󰉋  Files", m.files.View(), w, h)
	chatPanel := m.renderPanel(PanelChat, "  Gemini", m.chat.View(), w, h)
	todoPanel := m.renderPanel(PanelTodo, "  Tasks", m.todo.View(), w, h)

	// Join the three main panels horizontally side-by-side.
	mainBody := lipgloss.JoinHorizontal(lipgloss.Top, filesPanel, chatPanel, todoPanel)

	// Render the bottom terminal panel.
	termPanel := m.renderPanel(PanelTerminal, "  Terminal", m.term.View(), m.width-2, m.terminalHeight())

	// Render the top and bottom bars.
	header := m.renderHeader()
	status := m.renderStatus(m.width)

	// We'll show the ASCII logo at the top if there is enough vertical space
	// to avoid completely destroying the UI on extremely small screens.
	logoStr := ""
	if m.height >= 30 {
		logoStr = krakenLogo + "\n"
	}

	// Join everything vertically.
	return lipgloss.JoinVertical(lipgloss.Left, logoStr, header, mainBody, termPanel, status)
}

// ── Rendering helpers ─────────────────────────────────────────────────────────

// renderPanel wraps a child's raw view output in a stylized border.
func (m Model) renderPanel(idx int, title, content string, w, h int) string {
	active := m.activePanel == idx
	titleStr := styles.PanelTitle(active).Render(title)

	// Join the title text with the main content text.
	inner := lipgloss.JoinVertical(lipgloss.Left, titleStr, content)

	// Apply the bright orange border if active, otherwise dim gray.
	if active {
		return styles.PanelActive(w, h).Render(inner)
	}
	return styles.PanelInactive(w, h).Render(inner)
}

// renderHeader draws the top bar containing the app title and the panel tabs.
func (m Model) renderHeader() string {
	logo := styles.AppTitle.Render("🐙 KRAKEN")
	sub := styles.Dim.Render(" — AI · Files · Tasks")

	tabs := []string{"[Files]", "[Chat]", "[Tasks]"}
	var tabStr strings.Builder
	for i, t := range tabs {
		// Highlight the tab corresponding to the active panel
		if i == m.activePanel {
			tabStr.WriteString(styles.ChatSessionTabActive.Render(t))
		} else {
			tabStr.WriteString(styles.ChatSessionTab.Render(t))
		}
	}

	left := logo + sub
	right := tabStr.String()

	// Calculate spacing needed to push the tabs to the right edge
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}
	line := left + strings.Repeat(" ", gap) + right

	return styles.StatusBar.Width(m.width).Render(line)
}

// renderStatus draws the bottom bar displaying context-aware keybindings.
func (m Model) renderStatus(width int) string {
	// Base commands always available
	pills := []string{
		styles.HelpPill("Tab", "switch panel"),
		styles.HelpPill("^C", "quit"),
	}

	// Panel-specific commands
	switch m.activePanel {
	case PanelFiles:
		pills = append(pills,
			styles.HelpPill("↑↓", "navigate"),
			styles.HelpPill("Enter", "open"),
			styles.HelpPill("n", "new file"),
			styles.HelpPill("N", "new dir"),
			styles.HelpPill("r", "rename"),
			styles.HelpPill("d", "delete"),
			styles.HelpPill("y/p", "copy/paste"),
			styles.HelpPill("x", "cut"),
			styles.HelpPill(".", "hidden"),
		)
	case PanelChat:
		pills = append(pills,
			styles.HelpPill("Enter", "send"),
			styles.HelpPill("Alt+N", "new session"),
			styles.HelpPill("Alt+←/→", "sessions"),
		)
	case PanelTodo:
		pills = append(pills,
			styles.HelpPill("n", "new task"),
			styles.HelpPill("Space", "toggle"),
			styles.HelpPill("d", "delete"),
		)
	case PanelTerminal:
		pills = append(pills,
			styles.HelpPill("Enter", "run command"),
			styles.HelpPill("PgUp/PgDn", "scroll history"),
		)
	}

	bar := strings.Join(pills, " ")
	return styles.StatusBar.Width(width).Render(bar)
}

// terminalHeight calculates how tall the terminal panel should be (fixed 12 lines for now)
func (m Model) terminalHeight() int {
	return 12
}

// panelDimensions calculates the width and height available to each of the three
// top panels, taking into account the header bar (1 line), status bar (1 line),
// the terminal panel, and optionally the ASCII logo.
func (m Model) panelDimensions() (int, int) {
	headerH := 1
	statusH := 1
	termH := m.terminalHeight() + 2 // include borders

	logoH := 0
	if m.height >= 30 {
		logoH = lipgloss.Height(krakenLogo)
	}

	// Divide the width evenly by 3, subtract 2 for border strokes
	panelW := (m.width / 3) - 2
	panelH := m.height - headerH - statusH - termH - logoH - 2

	// Enforce minimum dimensions to prevent layout panics during extreme resizing
	if panelW < 10 {
		panelW = 10
	}
	if panelH < 5 {
		panelH = 5
	}
	return panelW, panelH
}
