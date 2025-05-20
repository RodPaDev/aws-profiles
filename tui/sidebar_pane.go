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
	cursor      state.Cursor
	State       *state.State
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

		if m.cursor.HandleKey(msg, len(m.profileKeys)) {
			break
		}

		switch msg.String() {
		case tea.KeySpace.String(), tea.KeyEnter.String():
			if m.State.Selection.Index == m.cursor.Index {
				m.State.Selection.Index = -1
				m.State.Selection.Key = ""
			} else {
				m.State.Selection.Index = m.cursor.Index
				m.State.Selection.Key = m.profileKeys[m.cursor.Index]
			}

		}

	}
	return m, nil
}

func (m *SidebarPaneModel) View() string {
	lines := make([]string, len(m.profileKeys))

	for idx, profile := range m.profileKeys {
		prefix := Icons.SidebarCursorInactive
		style := lipgloss.NewStyle().Width(m.size.Width)
		profileStyle := lipgloss.NewStyle()
		suffix := ""

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
		case m.State.IsProfileEdited(profile):
			profileStyle = style.Bold(true).Underline(true)
			suffix = Icons.ModifiedMarker
		default:
			fg = lipgloss.Color(Colors.Text)
		}

		style = style.Foreground(fg)

		if isCurrentSelection {
			style = style.Background(bg)
		} else {
			style = style.UnsetBackground()
		}

		lines[idx] = style.Render(fmt.Sprintf(" %s %s", prefix, profileStyle.Render(profile+suffix)))
	}

	return lipgloss.NewStyle().
		Height(m.size.Height).
		Width(m.size.Width).
		Render(lipgloss.JoinVertical(lipgloss.Top, lines...))
}
