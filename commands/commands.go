package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ErrorMsg struct {
	Error error
}

// Used to change to which board that you would like to be on
type ChangeScreenBoard struct {
	BoardId int
	Name string
}

func ChangeScreenBoardCmd(boardId int, name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoard{BoardId: boardId, Name: name}
	}
}

// Used to change back to the main boards page
type ChangeScreenBoardsMsg struct{ CurrentBoardId int }

func ChangeScreenBoardsCmd(boardId int) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreenBoardsMsg{CurrentBoardId: boardId}
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

