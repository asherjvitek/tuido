package boards

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"tuido/commands"
	"tuido/board"
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

type Model struct {
	Boards   []board.Model
	selected int
	width    int
	height   int
}

func (m Model) New() Model {
	return Model{
		Boards: []board.Model {
			board.New(1),
		},
		selected: 0,
		width: 0,
		height: 0,
	}
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			var cmd tea.Cmd
			cmd = func() tea.Msg {
				return commands.ChangeScreenBoard { BoardId: m.Boards[m.selected].Id }
			}
			return m, cmd
		case "N":
			var cmd tea.Cmd
			cmd = func() tea.Msg {
				return commands.NewBoard {}
			}
			return m, cmd
		case "right", "l":
			if m.selected == len(m.Boards)-1 {
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
