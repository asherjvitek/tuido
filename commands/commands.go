package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Used to change to which board that you would like to be on
type ChangeScreenBoard struct {
	BoardId int
}

func ChangeScreenBoardCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoard{BoardId: boardId}
	}
}

// Used to change back to the main boards page
type ChangeScreenBoards struct{ CurrentBoardId int }

func ChangeScreenBoardsCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoards{CurrentBoardId: boardId}
	}
}


type NewBoard struct{}

func NewBoardMsg() tea.Msg {
	return NewBoard{}
}

type SaveData struct{}

func SaveDataMsg() tea.Msg {
	return SaveData{}
}
