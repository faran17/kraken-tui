package styles

import "github.com/charmbracelet/lipgloss"

// ── Palette ──────────────────────────────────────────────────────────────────

const (
	// Ocean depths (backgrounds)
	ColorBg         = "#070D1A"
	ColorBgPanel    = "#0D1B2A"
	ColorBgHover    = "#122336"
	ColorBgSelected = "#1A3350"

	// Ocean water (blues)
	ColorOceanDark   = "#1B3A5C"
	ColorOceanMid    = "#1E6FBF"
	ColorOceanBright = "#3EA6FF"
	ColorOceanLight  = "#93C5FD"

	// Kraken orange (accent / warnings / user messages)
	ColorKrakenOrange     = "#FF6B35"
	ColorKrakenOrangeDark = "#C04A14"
	ColorKrakenOrangeGlow = "#FF9262"

	// Kraken green (success / AI messages / done items)
	ColorKrakenGreen     = "#00E5A0"
	ColorKrakenGreenDark = "#00A372"
	ColorKrakenGreenDim  = "#1A5C47"

	// Text
	ColorTextPrimary = "#E8F4FD"
	ColorTextSecond  = "#94A3B8"
	ColorTextDim     = "#4A5568"
	ColorTextMuted   = "#2D3748"

	// Status / system
	ColorError   = "#F87171"
	ColorWarning = "#FBBF24"
	ColorInfo    = ColorOceanBright
)

// ── Base Styles ───────────────────────────────────────────────────────────────

var (
	Logo = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorKrakenOrange)).
		Bold(true).
		MarginLeft(1)

	Base = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBg)).
		Foreground(lipgloss.Color(ColorTextPrimary))

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextDim))

	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextMuted))

	Bold = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Bold(true)
)

// ── Panel Borders ─────────────────────────────────────────────────────────────

var oceanBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "╰",
	BottomRight: "╯",
}

// PanelInactive is the base panel style (unfocused).
func PanelInactive(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(oceanBorder).
		BorderForeground(lipgloss.Color(ColorOceanDark)).
		Background(lipgloss.Color(ColorBgPanel)).
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Width(width).
		Height(height).
		Padding(0, 1)
}

// PanelActive is the focused panel style — glows with Kraken orange.
func PanelActive(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(oceanBorder).
		BorderForeground(lipgloss.Color(ColorKrakenOrange)).
		Background(lipgloss.Color(ColorBgPanel)).
		Foreground(lipgloss.Color(ColorTextPrimary)).
		Width(width).
		Height(height).
		Padding(0, 1)
}

// ── Title / Header ────────────────────────────────────────────────────────────

func PanelTitle(active bool) lipgloss.Style {
	c := ColorOceanBright
	if active {
		c = ColorKrakenOrange
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c)).
		Bold(true).
		Padding(0, 1)
}

var AppTitle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorKrakenGreen)).
	Bold(true)

// ── File Browser ──────────────────────────────────────────────────────────────

var (
	FileDir = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorOceanBright)).
		Bold(true)

	FileRegular = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextPrimary))

	FileHidden = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextDim))

	FileExec = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreen))

	FileSelected = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorBgSelected)).
			Foreground(lipgloss.Color(ColorKrakenOrangeGlow)).
			Bold(true)

	FileSize = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextSecond)).
			Width(8).
			Align(lipgloss.Right)

	FilePerm = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextDim)).
			Width(10)

	FilePath = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorOceanLight)).
			Italic(true)

	FilePrompt = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenOrange)).
			Bold(true)

	FileConfirmDanger = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorError)).
				Bold(true)
)

// ── Chat ──────────────────────────────────────────────────────────────────────

var (
	ChatUserMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenOrange)).
			Bold(true)

	ChatUserBubble = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorKrakenOrangeDark)).
			Foreground(lipgloss.Color(ColorTextPrimary)).
			Padding(0, 1).
			MarginBottom(1)

	ChatAIMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreen)).
			Bold(true)

	ChatAIBubble = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorKrakenGreenDim)).
			Foreground(lipgloss.Color(ColorTextPrimary)).
			Padding(0, 1).
			MarginBottom(1)

	ChatSystemMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextDim)).
			Italic(true)

	ChatSessionTab = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color(ColorTextDim))

	ChatSessionTabActive = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(lipgloss.Color(ColorKrakenGreen)).
				Bold(true).
				Underline(true)

	ChatSpinner = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreen))

	ChatInput = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color(ColorOceanDark)).
			Padding(0, 1)
)

// ── Todo ──────────────────────────────────────────────────────────────────────

var (
	TodoDone = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreenDark)).
			Strikethrough(true)

	TodoPending = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextPrimary))

	TodoSelected = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorBgSelected)).
			Foreground(lipgloss.Color(ColorKrakenOrangeGlow)).
			Bold(true)

	TodoCheckDone = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreen)).
			Bold(true)

	TodoCheckPending = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorTextDim))

	TodoInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenOrange))

	TodoCount = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorTextSecond)).
			Italic(true)
)

// ── Status / Help Bar ─────────────────────────────────────────────────────────

var (
	StatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorOceanDark)).
			Foreground(lipgloss.Color(ColorTextPrimary)).
			Padding(0, 1)

	StatusKey = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorOceanMid)).
			Foreground(lipgloss.Color(ColorTextPrimary)).
			Padding(0, 1).
			Bold(true)

	StatusVal = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorOceanDark)).
			Foreground(lipgloss.Color(ColorTextSecond)).
			Padding(0, 1)

	StatusErr = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorBg)).
			Foreground(lipgloss.Color(ColorError)).
			Padding(0, 1).
			Bold(true)

	StatusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorKrakenGreen)).
			Padding(0, 1)
)

// ── Helpers ───────────────────────────────────────────────────────────────────

// HelpPill renders a key + description pill for the help bar.
func HelpPill(key, desc string) string {
	k := StatusKey.Render(key)
	v := StatusVal.Render(desc)
	return k + v
}
