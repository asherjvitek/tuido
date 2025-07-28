package boards

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"tuido/board"
	"tuido/commands"
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
	Boards   []board.Model
	Selected int
	width    int
	height   int
}

func (m Model) New() Model {
	return Model{
		Boards: []board.Model{
			board.New(1),
		},
		Selected: 0,
		width:    0,
		height:   0,
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
	b := m.Boards[dest]

	m.Boards[dest] = a
	m.Boards[m.Selected] = b

	m.Selected = dest

	return m, commands.SaveDataMsg
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, commands.ChangeScreenBoardCmd(m.Boards[m.Selected].Id)
		case "N":
			return m, commands.NewBoardMsg
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
