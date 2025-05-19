package tui

import (
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/lib"
	"github.com/rodpadev/aws-profiles/state"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(Colors.Primary)
	contentStyle = lipgloss.NewStyle().PaddingLeft(2)
	noStyle      = lipgloss.NewStyle()
)

type ManagedInput struct {
	Name  string
	input textinput.Model
}

type ContentPaneModel struct {
	cursor                state.Cursor
	State                 *state.State
	DebugString           string
	size                  layout.PaneSize
	profile               state.ProfilePosition
	inputs                []ManagedInput
	modifiedInputs        []int
	isInputFocusCaptured  bool
	inputFocusCapturedIdx int
	inputsInitialized     bool
}

func (m *ContentPaneModel) BuildInputs(profile lib.Profile) {
	m.inputFocusCapturedIdx = -1
	m.isInputFocusCaptured = false
	m.inputs = make([]ManagedInput, 5)
	var t ManagedInput

	for i := range m.inputs {
		t.input = textinput.New()
		t.input.Blur()
		t.input.Width = m.size.Width
		t.input.Prompt = ""
		t.input.TextStyle = noStyle

		switch i {
		case 0:
			t.Name = "Profile Name"
			t.input.CharLimit = 32
			t.input.SetValue(profile.Name)
		case 1:
			t.Name = "Access Key"
			t.input.SetValue(profile.Credential.AccessKey)
		case 2:
			t.Name = "Secret Key"
			t.input.EchoMode = textinput.EchoPassword
			t.input.EchoCharacter = 0
			t.input.SetValue(profile.Credential.SecretKey)
		case 3:
			t.Name = "Output"
			t.input.SetValue(profile.Config.Output)
		case 4:
			t.Name = "Region"
			t.input.SetValue(profile.Config.Region)
		}
		m.inputs[i] = t
	}
}

func (m *ContentPaneModel) Init() tea.Cmd {
	return nil
}

func (m *ContentPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size

	case tea.KeyMsg:
		key := msg.String()

		if m.isInputFocusCaptured {
			cmd := m.updateFocusedInput(msg)
			if key == tea.KeyEscape.String() {
				m.modifiedInputs = append(m.modifiedInputs, m.inputFocusCapturedIdx)
				m.isInputFocusCaptured = false
				m.inputFocusCapturedIdx = -1
			}
			return m, cmd
		}

		m.cursor.HandleKey(msg, len(m.inputs))

		switch key {
		case tea.KeyEnter.String():
			m.isInputFocusCaptured = true
			m.inputFocusCapturedIdx = m.cursor.Index

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.inputFocusCapturedIdx {
					m.inputs[i].input.TextStyle = focusedStyle
					cmds[i] = m.inputs[i].input.Focus()
				} else {
					m.inputs[i].input.Blur()
					m.inputs[i].input.TextStyle = noStyle
				}
			}
			return m, tea.Batch(cmds...)

		case tea.KeyEscape.String():
			m.State.Selection.Index = -1
			m.State.Selection.Key = ""
		}
	}

	return m, nil
}

func (m *ContentPaneModel) updateFocusedInput(msg tea.Msg) tea.Cmd {
	idx := m.inputFocusCapturedIdx
	updatedInput, cmd := m.inputs[idx].input.Update(msg)
	m.inputs[idx].input = updatedInput
	return cmd
}

func (m *ContentPaneModel) View() string {
	selectedProfile, ok := m.State.GetCurrentProfile()
	if !ok {
		return "Select a profile to edit"
	}

	if ok && !m.inputsInitialized {
		m.BuildInputs(selectedProfile)
		m.inputsInitialized = true
	}

	var lines []string
	for i, field := range m.inputs {
		nameStyle := noStyle

		if slices.Contains(m.modifiedInputs, i) {
			nameStyle = nameStyle.Bold(true).Underline(true)
		}

		var name string
		if m.cursor.Index == i && m.inputFocusCapturedIdx != i {
			name = nameStyle.
				Background(Colors.Primary).
				Foreground(Colors.OnPrimary).
				Render(field.Name)
		} else {
			name = nameStyle.Render(field.Name)
		}

		if m.inputFocusCapturedIdx == i {
			field.input.TextStyle = focusedStyle
		} else {
			field.input.TextStyle = noStyle
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Render(name+": "),
			field.input.View(),
		)

		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Top, lines...)
	return contentStyle.Width(m.size.Width).Height(m.size.Height).Render(content)
}
