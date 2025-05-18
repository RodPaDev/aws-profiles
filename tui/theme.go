package tui

import "github.com/charmbracelet/lipgloss"

var Colors = struct {
	Text      lipgloss.Color
	Primary   lipgloss.Color
	OnPrimary lipgloss.Color
}{
	Text:      lipgloss.Color("#FFFFFF"),
	Primary:   lipgloss.Color("#FF9900"),
	OnPrimary: lipgloss.Color("#3b2301"),
}

var Icons = struct {
	SidebarCursorActive   string
	SidebarCursorInactive string
}{
	SidebarCursorActive:   "▶",
	SidebarCursorInactive: " ",
}
