package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/state"
)

type FooterPaneModel struct {
	size        layout.PaneSize
	DebugString string
	State       *state.State
}

func (m *FooterPaneModel) Init() tea.Cmd {
	return nil
}

func (m *FooterPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEscape.String():
			m.State.IsCommandModeActive = false
			return m, nil
		}
	}

	return m, nil
}

func (m *FooterPaneModel) View() string {
	if m.State.IsCommandModeActive {
		return fmt.Sprintf(":")
	}

	return fmt.Sprintf("Enter command mode with ':'")

}
