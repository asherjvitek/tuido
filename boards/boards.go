package boards

import (
	"fmt"
	"slices"
	"tuido/commands"
	"tuido/data"
	"tuido/util"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

var (
	boardStyle = lg.NewStyle().
			Width(20).
			Height(10).
			AlignHorizontal(lg.Center).
			AlignVertical(lg.Center).
			Border(lg.RoundedBorder())

	selectedStyle = boardStyle.BorderForeground(lg.Color("2"))
)

type Model struct {
	Boards   []data.Board
	Selected int
	Width    int
	Height   int
}

func (m Model) NextId() int {
	nextId := 0
	for _, b := range m.Boards {
		if nextId < b.Id {
			nextId = b.Id
		}
	}

	return nextId + 1
}

func (m *Model) navigate(dest int) {
	if dest < 0 || dest > len(m.Boards)-1 {
		return
	}

	m.Selected = dest
}

func (m *Model) moveBoard(dest int) tea.Model {
	if dest < 0 || dest > len(m.Boards)-1 || len(m.Boards) == 0 {
		return m
	}

	a := m.Boards[m.Selected]
	b := m.Boards[dest]

	a.Position = dest
	b.Position = m.Selected

	m.Boards[dest] = a
	m.Boards[m.Selected] = b

	m.Selected = dest

	err := data.UpdateBoard(a)

	if err != nil {
		util.Error(fmt.Sprintf("Error updating board %d/%s", a.Id, a.Name), err)
	}

	err = data.UpdateBoard(b)

	if err != nil {
		util.Error(fmt.Sprintf("Error updating board %d/%s", a.Id, a.Name), err)
	}

	return m
}

func (m *Model) newBoard() {
	board, err := data.NewBoard()

	if err != nil {
		util.Error("Error creating new board", err)
	}

	m.Boards = append(m.Boards, board)
}

func (m *Model) deleteBoard() {
	if len(m.Boards) == 0 {
		return
	}

	err := data.DeleteBoard(m.Boards[m.Selected])

	if err != nil {
		util.Error("Error deleting board", err)
	}

	m.Boards = slices.Delete(m.Boards, m.Selected, m.Selected+1)

	if m.Selected >= len(m.Boards)-1 {
		m.Selected--
	}
}

type initMsg struct {
	boards []data.Board
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		model, err := data.GetBoards()

		if err != nil {
			util.Error("Error getting boards", err)
		}

		return initMsg{boards: model}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case "enter":
			if len(m.Boards) == 0 {
				return m, nil
			}
			
			board := m.Boards[m.Selected]
			return m, commands.ChangeScreenBoardCmd(board.Id, board.Name)
		case "N":
			m.newBoard()
		case "D":
			m.deleteBoard()
		case "right", "l":
			m.navigate(m.Selected + 1)
		case "left", "h":
			m.navigate(m.Selected - 1)
		case "shift+right", "L":
			m.moveBoard(m.Selected + 1)
		case "shift+left", "H":
			m.moveBoard(m.Selected - 1)
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case initMsg:
		m.Boards = msg.boards
	}

	return m, nil
}

func (m Model) View() string {

	text := make([]string, len(m.Boards))

	for i, b := range m.Boards {
		if i == m.Selected {
			text[i] = selectedStyle.Render(b.Name)
		} else {
			text[i] = boardStyle.Render(b.Name)
		}
	}

	return lg.NewStyle().
		Width(m.Width).
		Height(m.Height).
		AlignHorizontal(lg.Center).
		AlignVertical(lg.Center).
		Render(lg.JoinHorizontal(lg.Center, text...))
}
