package main

import (
	"fmt"
	"os"
	"tuido/board"
	"tuido/boards"
	"tuido/commands"
	"tuido/data"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	screen tea.Model
	width  int
	height int
}

func (m model) windowSizeMsg() tea.Msg {
	return tea.WindowSizeMsg{
		Height: m.height,
		Width:  m.width,
	}
}

func (m *model) changeToBoards(boardId int) (tea.Model, tea.Cmd) {
	m.screen = boards.Model{Selected: 0, Height: m.height, Width: m.width}

	return m, m.screen.Init()
}

func (m *model) changeToBoard(msg commands.ChangeScreenBoard) (tea.Model, tea.Cmd) {
	m.screen = board.Model{Id: msg.BoardId, Name: msg.Name, Height: m.height, Width: m.width}

	return m, m.screen.Init()
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		s, ok := m.screen.(board.Model)
		if ok && s.Editing {
			break
		}

		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		}
	case commands.ChangeScreenBoardsMsg:
		return m.changeToBoards(msg.CurrentBoardId)
	case commands.ChangeScreenBoard:
		return m.changeToBoard(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.screen, cmd = m.screen.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.screen.View()
}

func main() {
	err := data.InitDatabase()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{screen: boards.Model{}}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Something is wrong %v", err)
		os.Exit(1)
	}
}
