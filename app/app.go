package app

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
	screen   tea.Model
	provider data.Provider
	width    int
	height   int
}

func (m model) windowSizeMsg() tea.Msg {
	return tea.WindowSizeMsg{
		Height: m.height,
		Width:  m.width,
	}
}

func (m *model) changeToBoards(boardId int) (tea.Model, tea.Cmd) {
	m.screen = boards.Model{Provider: m.provider}
	return m, tea.Batch(m.screen.Init(), boards.SelectedBoardIdCmd(boardId), m.windowSizeMsg)
}

func (m *model) changeToBoard(msg commands.ChangeScreenBoard) (tea.Model, tea.Cmd) {
	m.screen = board.Model{
		Board:    data.Board(msg),
		Provider: m.provider,
	}

	return m, tea.Batch(m.screen.Init(), m.windowSizeMsg)
}

func Run(provider data.Provider) {
	model := model{
		screen: boards.Model{Provider: provider}, 
		provider: provider,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Something is wrong %v", err)
		os.Exit(1)
	}
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
	case error:
		// Not sure what the right way to do this is yet but we will get there.
		panic(msg)
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
