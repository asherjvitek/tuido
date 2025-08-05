package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"tuido/data"
)

// Used to change to which board that you would like to be on
type ChangeScreenBoard data.Board
func ChangeScreenBoardCmd(board data.Board) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoard(board)
	}
}

// Used to change back to the main boards page
type ChangeScreenBoards struct{ CurrentBoardId int }

func ChangeScreenBoardsCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoards{CurrentBoardId: boardId}
	}
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}
