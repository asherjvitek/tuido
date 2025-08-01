package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"tuido/board"
	"tuido/boards"
	"tuido/commands"
	"tuido/data"
)

type model struct {
	screen tea.Model
	boards boards.Model
	width  int
	height int
}

func (m model) windowSizeMsg() tea.Msg {
	return tea.WindowSizeMsg{
		Height: m.height,
		Width:  m.width,
	}
}

func (m *model) updateBoardsModel() {
	// I know that I am missing something but it would appear that the name
	// updates are not flowing through but other updates to the screen model are.....
	// I am sure that I could fix this another way but for now this is fine
	switch screen := m.screen.(type) {
	case board.Model:
		for i, b := range m.boards.Boards {
			if b.Id == screen.Id {
				m.boards.Boards[i] = screen
				tea.SetWindowTitle(fmt.Sprintf("tuido - %s", screen.Name))
				return
			}
		}
	}
}

func (m *model) changeToBoards(boardId int) (tea.Model, tea.Cmd) {
	m.updateBoardsModel()

	for i, v := range m.boards.Boards {
		if v.Id == boardId {
			m.boards.Selected = i
		}
	}
	m.screen = m.boards

	return m, m.windowSizeMsg
}

func (m *model) changeToBoard(msg commands.ChangeScreenBoard) (tea.Model, tea.Cmd) {
	for _, board := range m.boards.Boards {
		if board.Id == msg.BoardId {
			m.screen = board

			return m, m.windowSizeMsg
		}
	}

	tea.SetWindowTitle("tuido - Boards")

	// should this panic or something?
	return m, nil
}

func (m *model) newBoard() (tea.Model, tea.Cmd) {
	newBoard := board.New(m.boards.NextId())
	m.boards.Boards = append(m.boards.Boards, newBoard)
	m.screen = newBoard

	return m, m.windowSizeMsg
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
	case commands.ChangeScreenBoards:
		return m.changeToBoards(msg.CurrentBoardId)
	case commands.ChangeScreenBoard:
		return m.changeToBoard(msg)
	case commands.NewBoard:
		return m.newBoard()
	case commands.SaveData:
		m.updateBoardsModel()
		data.SaveData(m.boards)
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
	boards := data.LoadData()
	model := model{screen: boards, boards: boards}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Something is wrong %v", err)
		os.Exit(1)
	}
}
