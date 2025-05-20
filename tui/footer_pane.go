package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/state"
)

type FooterPaneModel struct {
	DebugString string
	size        layout.PaneSize
	State       *state.State
}

func (m *FooterPaneModel) Init() tea.Cmd {
	return nil
}

func (m *FooterPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	}
	return m, nil
}

func (m *FooterPaneModel) View() string {
	return fmt.Sprintf("%s: %dx%d", m.DebugString, m.size.Width, m.size.Height)
}
