package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

//Used to change to which board that you would like to be on
type changeScreenBoard struct {
	boardId int
}

// Used to change back to the main boards page
type changeScreenBoards struct {
}

type newBoard struct {
}

type boardUpdated struct {

}

func saveBoard() tea.Msg {
	return boardUpdated{}
}
