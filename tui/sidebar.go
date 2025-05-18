package tui

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/state"
)

type SidebarModel struct {
	State       state.State
	DebugString string
	size        layout.PaneSize
	profileKeys []string
}

func (m *SidebarModel) Init() tea.Cmd {
	for k := range m.State.ProfileMap {
		m.profileKeys = append(m.profileKeys, k)
	}

	slices.Sort(m.profileKeys)
	return nil
}

func (m *SidebarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "spacebar":
			m.State.Selection = m.State.Cursor
		case "j", tea.KeyDown.String():
			if m.State.Cursor.Index < len(m.profileKeys)-1 {
				m.State.Cursor.Index += 1
			} else {
				m.State.Cursor.Index = 0
			}
		case "k", tea.KeyUp.String():
			if m.State.Cursor.Index > 0 {
				m.State.Cursor.Index -= 1
			} else {
				m.State.Cursor.Index = len(m.profileKeys) - 1
			}
		case "h", tea.KeyLeft.String():
			m.State.Cursor.Index = 0
		case "l", tea.KeyRight.String():
			m.State.Cursor.Index = len(m.profileKeys) - 1
		}

	}
	return m, nil
}

func (m *SidebarModel) View() string {
	lines := make([]string, len(m.profileKeys))

	for idx, profile := range m.profileKeys {
		prefix := Icons.SidebarCursorInactive
		style := lipgloss.NewStyle().Width(m.size.Width)

		isActiveIndex := m.State.Cursor.Index == idx
		isCurrentSelection := m.State.Selection.Index == idx

		var fg, bg lipgloss.Color

		switch {
		case isActiveIndex && isCurrentSelection:
			prefix = Icons.SidebarCursorActive
			fg = lipgloss.Color(Colors.OnPrimary)
			bg = lipgloss.Color(Colors.Primary)
		case isActiveIndex:
			prefix = Icons.SidebarCursorActive
			fg = lipgloss.Color(Colors.Primary)
		case isCurrentSelection:
			fg = lipgloss.Color(Colors.OnPrimary)
			bg = lipgloss.Color(Colors.Primary)
		default:
			fg = lipgloss.Color(Colors.Text)
		}

		style = style.Foreground(fg)

		if isCurrentSelection {
			style = style.Background(bg)
		} else {
			style = style.UnsetBackground()
		}

		lines[idx] = style.Render(fmt.Sprintf(" %s %s", prefix, profile))
	}

	return lipgloss.NewStyle().
		Height(m.size.Height).
		Width(m.size.Width).
		Render(lipgloss.JoinVertical(lipgloss.Top, lines...))
}
