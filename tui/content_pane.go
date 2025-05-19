package tui

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/state"
)

type ContentPaneModel struct {
	State       *state.State
	DebugString string
	size        layout.PaneSize
	profile     state.ProfilePosition
}

func (m ContentPaneModel) Init() tea.Cmd {
	return nil
}

func (m ContentPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	}
	return m, nil
}

func (m ContentPaneModel) View() string {
	if selected, ok := m.State.GetCurrentProfile(); ok {
		if marshall, err := json.MarshalIndent(selected, "", "  "); err == nil {
			return string(marshall)
		}
	}
	return "no profile or no unmarshall"
}
