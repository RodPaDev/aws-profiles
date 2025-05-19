package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/lib"
	"github.com/rodpadev/aws-profiles/state"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	contentStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type ManagedInput struct {
	Name  string
	idx   int
	input textinput.Model
}

type ContentPaneModel struct {
	cursor      state.Cursor
	State       *state.State
	DebugString string
	size        layout.PaneSize
	profile     state.ProfilePosition
	inputs      []ManagedInput
}

func (m *ContentPaneModel) BuildInputs(profile lib.Profile) {

	m.inputs = make([]ManagedInput, 5)
	var t ManagedInput

	for i := range m.inputs {
		t.idx = i
		t.input = textinput.New()
		t.input.Width = m.size.Width
		t.input.Prompt = ""

		switch i {
		case 0:
			t.Name = "Profile Name"
			t.input.Focus()
			t.input.PromptStyle = focusedStyle
			t.input.TextStyle = focusedStyle
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
		if m.cursor.HandleKey(msg, len(m.inputs)) {
			break
		}
	}

	return m, nil
}

func (m *ContentPaneModel) View() string {

	selectedProfile, ok := m.State.GetCurrentProfile()
	if !ok {
		return "Select a profile to edit"
	}

	m.BuildInputs(selectedProfile)
	var b strings.Builder

	for i := range m.inputs {
		if m.cursor.Index == i {
			m.inputs[i].input.Focus()
			m.inputs[i].input.PromptStyle = focusedStyle
			m.inputs[i].input.TextStyle = focusedStyle
		} else {
			m.inputs[i].input.Blur()
			m.inputs[i].input.PromptStyle = lipgloss.NewStyle()
			m.inputs[i].input.TextStyle = lipgloss.NewStyle()
		}
		b.WriteString(fmt.Sprintf("%s: %s", m.inputs[i].Name, m.inputs[i].input.View()))
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	mainRender := contentStyle.Width(m.size.Width).Height(m.size.Height).Render(b.String())
	return mainRender
}
