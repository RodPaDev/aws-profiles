package tui

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/state"
)

type SidebarPaneModel struct {
	cursor      state.ProfilePosition
	State       *state.State
	DebugString string
	size        layout.PaneSize
	profileKeys []string
}

func (m *SidebarPaneModel) Init() tea.Cmd {
	for k := range m.State.ProfileMap {
		m.profileKeys = append(m.profileKeys, k)
	}

	slices.Sort(m.profileKeys)
	return nil
}

func (m *SidebarPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeySpace.String(), tea.KeyEnter.String():
			if m.State.Selection.Index == m.cursor.Index {
				m.State.Selection.Index = -1
				m.State.Selection.Key = ""
			} else {
				m.State.Selection.Index = m.cursor.Index
				m.State.Selection.Key = m.profileKeys[m.cursor.Index]
			}
		case tea.KeyBackspace.String():
			m.State.Selection.Index = -1
			m.State.Selection.Key = ""
		case "j", tea.KeyDown.String():
			if m.cursor.Index < len(m.profileKeys)-1 {
				m.cursor.Index += 1
			} else {
				m.cursor.Index = 0
			}
		case "k", tea.KeyUp.String():
			if m.cursor.Index > 0 {
				m.cursor.Index -= 1
			} else {
				m.cursor.Index = len(m.profileKeys) - 1
			}
		case "h", tea.KeyLeft.String():
			m.cursor.Index = 0
		case "l", tea.KeyRight.String():
			m.cursor.Index = len(m.profileKeys) - 1
		}

	}
	return m, nil
}

func (m *SidebarPaneModel) View() string {
	lines := make([]string, len(m.profileKeys))

	for idx, profile := range m.profileKeys {
		prefix := Icons.SidebarCursorInactive
		style := lipgloss.NewStyle().Width(m.size.Width)

		isActiveIndex := m.cursor.Index == idx
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
