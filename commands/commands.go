package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Used to change to which board that you would like to be on
type ChangeScreenBoard struct {
	BoardId int
}

// Used to change back to the main boards page
type ChangeScreenBoards struct{}

type NewBoard struct{}

type BoardUpdated struct{}

func SaveBoard() tea.Msg {
	return BoardUpdated{}
}
