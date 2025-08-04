package boards

import (
	"fmt"
	"tuido/commands"
	"tuido/data"

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
	width    int
	height   int
}

func (m *Model) navigate(dest int) {
	if dest < 0 || dest > len(m.Boards)-1 {
		return
	}

	m.Selected = dest
}

func (m *Model) moveBoard(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > len(m.Boards)-1 || len(m.Boards) == 0 {
		return m, nil
	}

	a := m.Boards[m.Selected]
	aPos := a.Position
	b := m.Boards[dest]

	a.Position = b.Position
	b.Position = aPos

	m.Boards[dest] = a
	m.Boards[m.Selected] = b

	m.Selected = dest

	err := data.UpdateBoard(a)

	if err != nil {
		panic(err)
	}


	err = data.UpdateBoard(b)

	if err != nil {
		panic(err)
	}


	return m, commands.SaveDataMsg
}

func (m *Model) createBoard() (tea.Model, tea.Cmd) {
	pos, err := data.GetPosition(m.Boards, len(m.Boards))

	if err != nil {
		panic(err)
	}

	board := data.Board{
		Name:     "New Board",
		Position: pos,
	}

	err = data.InsertBoard(&board)

	if err != nil {
		panic(fmt.Errorf("Failed to insert new board: %v", err))
	}

	pos, err = data.GetPosition([]data.List{}, 0)

	if err != nil {
		panic(err)
	}

	list := data.List{
		BoardId:  board.BoardId,
		Name:     "New List",
		Position: pos,
	}

	err = data.InsertList(&list)

	if err != nil {
		panic(fmt.Errorf("Failed to insert new list: %v", err))
	}

	item := data.Item{
		ListId:   list.ListId,
		Text:     "New Item",
		Position: pos,
	}

	err = data.InsertItem(&item)

	if err != nil {
		panic(fmt.Errorf("Failed to insert new item: %v", err))
	}

	return m, commands.ChangeScreenBoardCmd(board)
}

type initMsg []data.Board
type selectedBoardIdMsg int

func SelectedBoardIdCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return selectedBoardIdMsg(boardId)
	}
}

func (m Model) Init() tea.Cmd {
	boards, err := data.Boards()
	if err != nil {
		panic(err)
	}

	return func() tea.Msg { return initMsg(boards) }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, commands.ChangeScreenBoardCmd(m.Boards[m.Selected])
		case "N":
			return m.createBoard()
		case "right", "l":
			m.navigate(m.Selected + 1)
		case "left", "h":
			m.navigate(m.Selected - 1)
		case "shift+right", "L":
			return m.moveBoard(m.Selected + 1)
		case "shift+left", "H":
			return m.moveBoard(m.Selected - 1)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case initMsg:
		m.Boards = msg
	case selectedBoardIdMsg:
		boardId := int(msg)
		for i, b := range m.Boards {
			if b.BoardId == boardId {
				m.Selected = i
				break
			}
		}
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
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lg.Center).
		AlignVertical(lg.Center).
		Render(lg.JoinHorizontal(lg.Center, text...))
}
