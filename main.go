package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	leftPaneWidth int = 25
	footerHeight  int = 4
)

type layoutModel struct {
	width, height int
	leftPane      tea.Model
	rightPane     tea.Model
	footer        tea.Model
}

func (m layoutModel) Init() tea.Cmd {
	return tea.Batch(
		m.leftPane.Init(),
		m.rightPane.Init(),
		m.footer.Init(),
	)
}

func (m layoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	}

	var cmds []tea.Cmd

	m.leftPane, cmds = updateModel(m.leftPane, msg, cmds)
	m.rightPane, cmds = updateModel(m.rightPane, msg, cmds)
	m.footer, cmds = updateModel(m.footer, msg, cmds)

	return m, tea.Batch(cmds...)
}

func updateModel(child tea.Model, msg tea.Msg, cmds []tea.Cmd) (tea.Model, []tea.Cmd) {
	newModel, cmd := child.Update(msg)
	return newModel, append(cmds, cmd)
}

func (m layoutModel) View() string {
	heightWithoutBorder := m.height - 2
	footerVisibleHeight := footerHeight - 1
	mainHeight := heightWithoutBorder - footerVisibleHeight - 2

	right := lipgloss.NewStyle().
		Width(m.width - leftPaneWidth - 6).
		Height(mainHeight - 2).
		Border(lipgloss.NormalBorder()).
		Render(m.rightPane.View())

	left := lipgloss.NewStyle().
		Width(leftPaneWidth).
		Height(mainHeight - 2).
		Border(lipgloss.NormalBorder()).
		Render(m.leftPane.View())

	footer := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(footerHeight).
		Border(lipgloss.NormalBorder()).
		BorderBottom(false).
		BorderRight(false).
		BorderLeft(false).
		Render(m.footer.View())

	topRow := lipgloss.NewStyle().Height(mainHeight).Width(m.width - 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
	)

	mainRender := lipgloss.JoinVertical(lipgloss.Top, topRow, footer)

	if m.width < 75 || m.height < 25 {
		mainRender = lipgloss.Place(m.width, heightWithoutBorder, lipgloss.Center, lipgloss.Center, "Too Small! Expand your terminal")
	}

	return lipgloss.NewStyle().
		Height(heightWithoutBorder).
		Width(m.width - 2).
		Border(lipgloss.NormalBorder()).
		Render(mainRender)

}

func initModel() layoutModel {
	return layoutModel{
		leftPane: dummyLayout{
			debugString: "left",
		},
		rightPane: dummyLayout{
			debugString: "right",
		},
		footer: dummyLayout{
			debugString: "footer",
		},
	}
}

type dummyLayout struct {
	debugString string
}

func (m dummyLayout) Init() tea.Cmd {
	return nil
}

func (m dummyLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m dummyLayout) View() string {
	return m.debugString
}

func main() {

	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
