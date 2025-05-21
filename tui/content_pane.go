package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/lib"
	"github.com/rodpadev/aws-profiles/state"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(Colors.Primary)
	contentStyle = lipgloss.NewStyle().PaddingLeft(1)
	noStyle      = lipgloss.NewStyle()
)

type ManagedInput struct {
	Name  string
	input textinput.Model
	field string
}

type ContentPaneModel struct {
	cursor             state.Cursor
	State              *state.State
	size               layout.PaneSize
	profile            state.ProfilePosition
	inputs             []ManagedInput
	isInputFocusLocked bool
	focusedInputIndex  int
	inputsInitialized  bool
}

func (m *ContentPaneModel) BuildInputs(profile lib.Profile) {
	m.focusedInputIndex = -1
	m.isInputFocusLocked = false
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
			t.field = "name"
			t.Name = "Profile Name"
			t.input.CharLimit = 32
			t.input.SetValue(profile.Name)
		case 1:
			t.field = "aws_access_key_id"
			t.Name = "Access Key"
			t.input.SetValue(profile.Credential.AccessKey)
		case 2:
			t.field = "aws_secret_access_key"
			t.Name = "Secret Key"
			t.input.EchoMode = textinput.EchoPassword
			t.input.EchoCharacter = 0
			t.input.SetValue(profile.Credential.SecretKey)
		case 3:
			t.field = "output"
			t.Name = "Output"
			t.input.SetValue(profile.Config.Output)
		case 4:
			t.field = "region"
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

		if m.isInputFocusLocked {
			cmd := m.updateFocusedInput(msg)
			if key == tea.KeyEscape.String() {
				modifiedValue := m.inputs[m.focusedInputIndex].input.Value()
				profile, ok := m.State.GetCurrentProfile()
				if ok && modifiedValue != profile.GetField(m.inputs[m.focusedInputIndex].field) {

					key := m.State.Selection.Key
					edited, exists := m.State.EditedProfileMap[key]
					if !exists {
						edited = profile
					}

					edited.SetField(m.inputs[m.focusedInputIndex].field, modifiedValue)

					m.State.EditedProfileMap[key] = edited
				}
				m.isInputFocusLocked = false
				m.focusedInputIndex = -1
			}
			return m, cmd
		}

		m.cursor.HandleKey(msg, len(m.inputs))

		switch key {
		case tea.KeyEnter.String():
			m.isInputFocusLocked = true
			m.focusedInputIndex = m.cursor.Index

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusedInputIndex {
					m.inputs[i].input.TextStyle = focusedStyle
					cmds[i] = m.inputs[i].input.Focus()
				} else {
					m.inputs[i].input.Blur()
					m.inputs[i].input.TextStyle = noStyle
				}
			}
			return m, tea.Batch(cmds...)

		case tea.KeyEscape.String():
			// exit content pane
			// todo: this should be handled in the main loop
			m.State.Selection.Index = -1
			m.State.Selection.Key = ""
		}
	}

	return m, nil
}

func (m *ContentPaneModel) updateFocusedInput(msg tea.Msg) tea.Cmd {
	idx := m.focusedInputIndex
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
	for idx, input := range m.inputs {
		nameStyle := noStyle
		prefix := " "

		if m.State.IsFieldEdited(m.State.Selection.Key, input.field) {
			nameStyle = nameStyle.Bold(true).Underline(true)
			prefix = Icons.ModifiedMarker
		}

		var name string
		if m.cursor.Index == idx && m.focusedInputIndex != idx {
			name = nameStyle.
				Background(Colors.Primary).
				Foreground(Colors.OnPrimary).
				Render(input.Name)
		} else {
			name = nameStyle.Render(input.Name)
		}

		if m.focusedInputIndex == idx {
			input.input.TextStyle = focusedStyle
		} else {
			input.input.TextStyle = noStyle
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Render(
				fmt.Sprintf("%s %s: ", prefix, name),
			),
			input.input.View(),
		)

		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Top, lines...)
	return contentStyle.Width(m.size.Width).Height(m.size.Height).Render(content)
}
