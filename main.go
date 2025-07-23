package main

import (
	"fmt"
	"os"
	// "os"
	// "strings"
	//
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type list struct {
	title string
	items []string
}

type model struct {
	lists        []list
	width        int
	height       int
	selectedList int
	selectedItem int
	editing      bool
	input        textinput.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case "down":
			if m.selectedItem < len(m.lists[m.selectedList].items)-1 {
				m.selectedItem++
			}
		case "up":
			if m.selectedItem > 0 {
				m.selectedItem--
			}
		case "right":
			if m.selectedList < len(m.lists)-1 {
				m.selectedList++
			}
		case "left":
			if m.selectedList > 0 {
				m.selectedList--
			}
		case "enter":
			if !m.editing {
				m.input.SetValue(m.lists[m.selectedList].items[m.selectedItem])
				m.input.Focus()
				m.editing = true
			}
		case "esc":
			if m.editing {
				m.input.Blur()
				m.lists[m.selectedList].items[m.selectedItem] = m.input.Value()
				m.editing = false
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

var (
	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#ff0000")).
				Align(lipgloss.Center)

	editingStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0000ff")).
			Align(lipgloss.Center)
)

func (m model) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	listWidth := (m.width-border.GetHorizontalBorderSize())/3 - len(m.lists)
	contentWidth := m.width - border.GetHorizontalBorderSize() - border.GetHorizontalPadding()
	contentHeight := m.height - border.GetVerticalBorderSize() - border.GetVerticalPadding()

	lists := make([]string, len(m.lists))
	for li, v := range m.lists {
		title := border.
			BorderBottom(true).
			BorderTop(false).
			BorderRight(false).
			BorderLeft(false).
			Width(listWidth - 2).
			Align(lipgloss.Center).
			Render(v.title)

		items := make([]string, len(v.items))

		for ii, v := range v.items {
			if m.selectedList == li && m.selectedItem == ii {
				if m.editing {
					items[ii] = border.Render(m.input.View())
				} else {
					items[ii] = border.Render(selectedItemStyle.Render(v))
				}
			} else {
				items[ii] = border.Render(v)
			}
		}

		todos := lipgloss.JoinVertical(lipgloss.Left, items...)
		list := lipgloss.JoinVertical(lipgloss.Left, title, todos)

		text := lipgloss.NewStyle().
			Width(contentWidth / 3).
			Height(contentHeight - 3).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Top).
			Render(list)

		lists = append(lists, border.
			Width(listWidth).
			Height(20).
			Render(text))
	}

	view := border.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, lists...))

	return view
}

func main() {
	model := model{
		lists: []list{
			{
				title: "TODO",
				items: []string{
					"Thing to Do 1",
					"Thing to Do 2",
					"Thing to Do 3",
					"Thing to Do 4",
					"Thing to Do 4",
					"Thing to Do 5",
					"Thing to Do 6",
				},
			},
			{
				title: "DOING",
				items: []string{
					"Doing this thing 1",
					"Doing this thing 2",
					"Doing this thing 1",
				},
			},
			{
				title: "DONE",
				items: []string{
					"This is done 1!",
					"This is done 2!",
					"This is done 3!",
				},
			},
		},
		selectedList: 0,
		selectedItem: 0,
		input:        textinput.New(),
	}

	model.input.Prompt = ""

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Something is wrong %v", err)
		os.Exit(1)
	}
}
