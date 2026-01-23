package ui

import "github.com/charmbracelet/lipgloss"

// Capacitor brand colors
var (
	// Primary colors - Capacitor blue
	capBlue  = lipgloss.Color("#119EFF")
	capCyan  = lipgloss.Color("#73B7F6")
	capDark  = lipgloss.Color("#16161D")
	capLight = lipgloss.Color("#ECEDEE")
	capGray  = lipgloss.Color("#4A4A5A")

	// Status colors
	successColor = lipgloss.Color("#4ADE80")
	errorColor   = lipgloss.Color("#F87171")
	warnColor    = lipgloss.Color("#FBBF24")
	mutedColor   = lipgloss.Color("#64748B")

	// Platform colors
	iosColor     = lipgloss.Color("#0A84FF")
	androidColor = lipgloss.Color("#34D399")
	webColor     = lipgloss.Color("#F97316") // Orange for web
)

// CapacitorLogo returns the logo for welcome screen
func CapacitorLogo() string {
	textStyle := lipgloss.NewStyle().Foreground(capLight).Bold(true)

	lines := []string{
		"",
		textStyle.Render("lazycap"),
		mutedStyle.Render("Capacitor Dashboard"),
	}

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

// LogoCompact returns a compact inline logo for the header
func LogoCompact() string {
	bolt := lipgloss.NewStyle().Foreground(capBlue).Bold(true).Render("⚡")
	name := lipgloss.NewStyle().Foreground(capLight).Bold(true).Render("lazycap")
	return bolt + " " + name
}

// LogoSmall returns a minimal logo for header (alias for LogoCompact)
func LogoSmall() string {
	return LogoCompact()
}

// Styles
var (
	// Project name in header
	projectStyle = lipgloss.NewStyle().
			Foreground(capCyan)

	// Section titles
	titleStyle = lipgloss.NewStyle().
			Foreground(capBlue).
			Bold(true).
			MarginBottom(1)

	// Active pane border
	activePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(capBlue).
			Padding(1, 2)

	// Inactive pane border
	inactivePaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(capGray).
				Padding(1, 2)

	// Device status
	onlineStyle = lipgloss.NewStyle().
			Foreground(successColor)

	offlineStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	failedStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Log pane
	logPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(capGray).
			Padding(0, 1)

	activeLogPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(capBlue).
				Padding(0, 1)

	logEmptyStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Help bar
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(capCyan).
			Bold(true)

	// Badges
	iosBadge = lipgloss.NewStyle().
			Foreground(iosColor).
			Bold(true)

	androidBadge = lipgloss.NewStyle().
			Foreground(androidColor).
			Bold(true)

	webBadge = lipgloss.NewStyle().
			Foreground(webColor).
			Bold(true)

	// Upgrade notice
	upgradeStyle = lipgloss.NewStyle().
			Foreground(warnColor).
			Bold(true)

	// Muted text
	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Error
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Tab styles for process tabs
	activeTabStyle = lipgloss.NewStyle().
			Foreground(capDark).
			Background(capBlue).
			Padding(0, 1).
			MarginRight(1).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(capLight).
				Background(capGray).
				Padding(0, 1).
				MarginRight(1)
)

// StatusDot returns a colored dot
func StatusDot(online bool) string {
	if online {
		return onlineStyle.Render("●")
	}
	return offlineStyle.Render("○")
}

// PlatformBadge returns styled platform text
func PlatformBadge(platform string) string {
	switch platform {
	case "ios":
		return iosBadge.Render("iOS")
	case "android":
		return androidBadge.Render("Android")
	case "web":
		return webBadge.Render("Web")
	}
	return platform
}
