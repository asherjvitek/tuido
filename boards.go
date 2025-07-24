package main

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	// "slices"
)

var (
	style = lg.NewStyle().
		AlignHorizontal(lg.Center).
		AlignVertical(lg.Center)

	boardStyle = lg.NewStyle().
			Width(20).
			Height(10).
			AlignHorizontal(lg.Center).
			AlignVertical(lg.Center).
			Border(lg.RoundedBorder())

	selectedStyle = boardStyle.BorderForeground(lg.Color("2"))
)

type boards struct {
	boards   []board
	selected int
	width    int
	height   int
}

func (m boards) Init() tea.Cmd {
	return nil
}

func (m boards) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			var cmd tea.Cmd
			cmd = func() tea.Msg {
				return changeScreenBoard { boardId: m.boards[m.selected].id }
			}
			return m, cmd
		case "N":
			var cmd tea.Cmd
			cmd = func() tea.Msg {
				return newBoard {}
			}
			return m, cmd
		case "right", "l":
			if m.selected == len(m.boards)-1 {
				break
			}

			m.selected++
		case "left", "h":
			if m.selected == 0 {
				break
			}

			m.selected--
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m boards) View() string {

	text := make([]string, len(m.boards))

	for i, b := range m.boards {
		if i == m.selected {
			text[i] = selectedStyle.Render(b.name)
		} else {
			text[i] = boardStyle.Render(b.name)
		}
	}

	return lg.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lg.Center).
		AlignVertical(lg.Center).
		Render(lg.JoinHorizontal(lg.Center, text...))
}
