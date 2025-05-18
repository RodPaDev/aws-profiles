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
		case "j", tea.KeyDown.String():
			if m.State.Selection.Index < len(m.profileKeys)-1 {
				m.State.Selection.Index += 1
			} else {
				m.State.Selection.Index = 0
			}
		case "k", tea.KeyUp.String():
			if m.State.Selection.Index > 0 {
				m.State.Selection.Index -= 1
			} else {
				m.State.Selection.Index = len(m.profileKeys) - 1
			}
		}
	}
	return m, nil
}

func (m *SidebarModel) View() string {

	lines := make([]string, len(m.profileKeys))

	for idx, profile := range m.profileKeys {
		prefix := " "
		// default color is white
		color := lipgloss.Color("white")
		if m.State.Selection.Index == idx {
			prefix = ">"
			color = lipgloss.Color("205")
		}
		lines[idx] = lipgloss.NewStyle().
			Foreground(color).
			Render(fmt.Sprintf("%s %s", prefix, profile))
	}

	return lipgloss.NewStyle().
		Height(m.size.Height).
		Width(m.size.Width).
		Render(lipgloss.JoinVertical(lipgloss.Top, lines...,
		))
}
