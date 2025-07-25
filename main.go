package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	screen tea.Model
	boards boards
	width  int
	height int
}

func (m *model) changeToBoards() (tea.Model, tea.Cmd) {
	m.boards.height = m.height
	m.boards.width = m.width

	// I know that I am missing something but it would appear that the name
	// updates are not flowing through but other updates to the screen model are.....
	// I am sure that I could fix this another way but for now this is fine
	switch screen := m.screen.(type) {
	case board:
		for i, b := range m.boards.boards {
			if b.id == screen.id {
				m.boards.boards[i] = screen
			}
		}
	}

	m.screen = m.boards

	return m, nil
}

func (m *model) changeToBoard(msg changeScreenBoard) (tea.Model, tea.Cmd) {
	for _, board := range m.boards.boards {
		if board.id == msg.boardId {
			board.height = m.height
			board.width = m.width
			m.screen = board

			return m, nil
		}
	}

	// should this panic or something?
	return m, nil
}

func (m *model) newBoard() (tea.Model, tea.Cmd) {
	nextId := 0
	for _, b := range m.boards.boards {
		if nextId < b.id {
			nextId = b.id
		}
	}

	nextId++
	newBoard := board{
		id:     nextId,
		name:   "New Board",
		lists:  make([]list, 0),
		input:  getTextInput(),
		height: m.height,
		width:  m.width,
	}
	m.boards.boards = append(m.boards.boards, newBoard)
	m.screen = newBoard

	return m, nil
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		}
	case changeScreenBoards:
		return m.changeToBoards()
	case changeScreenBoard:
		return m.changeToBoard(msg)
	case newBoard:
		return m.newBoard()
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
	boards := getBoards()
	model := model{screen: boards, boards: boards}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Something is wrong %v", err)
		os.Exit(1)
	}
}
