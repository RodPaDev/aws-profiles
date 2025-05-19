package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rodpadev/aws-profiles/layout"
	"github.com/rodpadev/aws-profiles/lib"
	"github.com/rodpadev/aws-profiles/state"
	"github.com/rodpadev/aws-profiles/tui"
)

const (
	leftPaneWidth int = 25
	footerHeight  int = 4
)

var AppState state.State

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

		computedLayout := layout.ComputeLayout(m.width, m.height)

		m.leftPane, _ = m.leftPane.Update(layout.SizeMsg{
			Size: layout.PaneSize{Width: computedLayout.LeftPaneWidth, Height: computedLayout.MainHeight - 2},
		})
		m.rightPane, _ = m.rightPane.Update(layout.SizeMsg{
			Size: layout.PaneSize{Width: computedLayout.RightPaneWidth, Height: computedLayout.MainHeight - 2},
		})
		m.footer, _ = m.footer.Update(layout.SizeMsg{
			Size: layout.PaneSize{Width: computedLayout.ScreenWidth - 2, Height: computedLayout.FooterHeight},
		})

	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyBackspace.String():
			AppState.Selection.Index = -1
			AppState.Selection.Key = ""
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	}

	var cmds []tea.Cmd

	if AppState.Selection.Key == "" {
		m.leftPane, cmds = updateModel(m.leftPane, msg, cmds)
	} else {
		m.rightPane, cmds = updateModel(m.rightPane, msg, cmds)
	}
	m.footer, cmds = updateModel(m.footer, msg, cmds)

	return m, tea.Batch(cmds...)
}

func updateModel(child tea.Model, msg tea.Msg, cmds []tea.Cmd) (tea.Model, []tea.Cmd) {
	newModel, cmd := child.Update(msg)
	return newModel, append(cmds, cmd)
}

func (m layoutModel) View() string {
	computedLayout := layout.ComputeLayout(m.width, m.height)

	right := lipgloss.NewStyle().
		Width(computedLayout.RightPaneWidth).
		Height(computedLayout.MainHeight - 2).
		Border(lipgloss.NormalBorder()).
		Render(m.rightPane.View())

	left := lipgloss.NewStyle().
		Width(computedLayout.LeftPaneWidth).
		Height(computedLayout.MainHeight - 2).
		Border(lipgloss.NormalBorder()).
		Render(m.leftPane.View())

	footer := lipgloss.NewStyle().
		Width(computedLayout.ScreenWidth - 2).
		Height(computedLayout.FooterHeight).
		Border(lipgloss.NormalBorder()).
		BorderBottom(false).
		BorderRight(false).
		BorderLeft(false).
		Render(m.footer.View())

	topRow := lipgloss.NewStyle().
		Height(computedLayout.MainHeight).
		Width(computedLayout.ScreenWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, left, right))

	mainRender := lipgloss.JoinVertical(lipgloss.Top, topRow, footer)

	if m.width < 75 || m.height < 25 {
		mainRender = lipgloss.Place(m.width, computedLayout.HeightAvailable, lipgloss.Center, lipgloss.Center, "Too Small! Expand your terminal")
	}

	return lipgloss.NewStyle().
		Height(computedLayout.HeightAvailable).
		Width(m.width - 2).
		Border(lipgloss.NormalBorder()).
		Render(mainRender)

}

func initModel() layoutModel {
	return layoutModel{
		leftPane: &tui.SidebarPaneModel{
			DebugString: "left",
			State:       &AppState,
		},
		rightPane: &tui.ContentPaneModel{
			DebugString: "right",
			State:       &AppState,
		},
		footer: dummyLayout{
			debugString: "footer",
		},
	}
}

type dummyLayout struct {
	debugString string
	size        layout.PaneSize
}

func (m dummyLayout) Init() tea.Cmd {
	return nil
}

func (m dummyLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.SizeMsg:
		m.size = msg.Size
	}
	return m, nil
}

func (m dummyLayout) View() string {
	return fmt.Sprintf("%s: %dx%d", m.debugString, m.size.Width, m.size.Height)
}

func main() {

	data, err := lib.LoadAWSProfileData()
	if err != nil {
		log.Fatal(err)
	}

	parsed, profileCount := lib.ParseAWSProfileData(data)
	if err != nil {
		log.Fatal(err)
	}

	profileMap, err := lib.BuildProfileMap(parsed)

	if err != nil {
		log.Fatal(err)
	}

	AppState = state.State{
		ProfileMap:     profileMap,
		ProfileMapSize: profileCount,
	}

	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
