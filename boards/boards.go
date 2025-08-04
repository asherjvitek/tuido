package boards

import (
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
	selected int
	width    int
	height   int
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
			return m, commands.ChangeScreenBoardCmd(m.Boards[m.selected])
		case "N":
			return m.createBoard()
		case "right", "l":
			m.navigate(m.selected + 1)
		case "left", "h":
			m.navigate(m.selected - 1)
		case "shift+right", "L":
			return m.moveBoard(m.selected + 1)
		case "shift+left", "H":
			return m.moveBoard(m.selected - 1)
		case "D":
			return m.deleteBoard()
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
				m.selected = i
				break
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	text := make([]string, len(m.Boards))

	for i, b := range m.Boards {
		if i == m.selected {
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
