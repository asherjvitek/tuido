package boards


import (
	"fmt"
	"slices"
	"tuido/commands"
	"tuido/data"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) navigate(dest int) {
	if dest < 0 || dest > len(m.Boards)-1 {
		return
	}

	m.selected = dest
}

func (m *Model) moveBoard(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > len(m.Boards)-1 || len(m.Boards) == 0 {
		return m, nil
	}

	a := m.Boards[m.selected]
	aPos := a.Position
	b := m.Boards[dest]

	a.Position = b.Position
	b.Position = aPos

	m.Boards[dest] = a
	m.Boards[m.selected] = b

	m.selected = dest

	err := data.UpdateBoard(a)

	if err != nil {
		panic(err)
	}

	err = data.UpdateBoard(b)

	if err != nil {
		panic(err)
	}

	return m, nil
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

func (m Model) deleteBoard() (tea.Model, tea.Cmd) {
	if len(m.Boards) == 0 {
		return m, nil
	}

	board := m.Boards[m.selected]
	err := data.DeleteBoard(board)

	if err != nil {
		panic(err)
	}

	m.Boards = slices.Delete(m.Boards, m.selected, m.selected+1)

	if m.selected >= len(m.Boards) {
		m.selected = len(m.Boards) - 1
	}

	return m, nil

}

type initMsg []data.Board
type selectedBoardIdMsg int

func SelectedBoardIdCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return selectedBoardIdMsg(boardId)
	}
}
