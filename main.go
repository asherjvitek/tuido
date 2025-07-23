package main

import (
	"fmt"
	"os"
	"slices"
	// "os"
	// "strings"
	//
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type list struct {
	title          string
	items          []string
	scrollposition int
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

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#23d18b"))

	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#2a3d41"))

	border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
)

func (m model) Init() tea.Cmd {
	return nil
}

func countLines(str string) int {
	lines := 1
	for _, r := range str {
		if r == '\n' {
			lines++
		}
	}

	return lines
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.input.Blur()
				m.lists[m.selectedList].items[m.selectedItem] = m.input.Value()
				m.editing = false
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)

				return m, cmd
			}
		} else {
			switch msg.String() {
			case tea.KeyCtrlC.String(), "q":
				return m, tea.Quit
			case "down", "j":
				if m.selectedItem < len(m.lists[m.selectedList].items)-1 {
					m.selectedItem++
				}
			case "up", "k":
				if m.selectedItem > 0 {
					m.selectedItem--
				}
			case "right", "l":
				if m.selectedList < len(m.lists)-1 {
					m.selectedList++
					if m.selectedItem > len(m.lists[m.selectedList].items) {
						m.selectedItem = len(m.lists[m.selectedList].items) - 1
					}
				}
			case "left", "h":
				if m.selectedList > 0 {
					m.selectedList--
				}
			case "ctrl+home":
				m.selectedList = 0
				if m.selectedItem > len(m.lists[m.selectedList].items) {
					m.selectedItem = len(m.lists[m.selectedList].items) - 1
				}
			case "ctrl+end":
				m.selectedList = len(m.lists) - 1
				if m.selectedItem > len(m.lists[m.selectedList].items) {
					m.selectedItem = len(m.lists[m.selectedList].items) - 1
				}
			case "home":
				m.selectedItem = 0
			case "end":
				m.selectedItem = len(m.lists[m.selectedList].items) - 1
			case "shift+down", "J":
				if m.selectedItem == len(m.lists[m.selectedList].items)-1 {
					break
				}

				a := m.lists[m.selectedList].items[m.selectedItem]
				b := m.lists[m.selectedList].items[m.selectedItem+1]

				m.lists[m.selectedList].items[m.selectedItem] = b
				m.lists[m.selectedList].items[m.selectedItem+1] = a

				m.selectedItem++
			case "shift+up", "K":
				if m.selectedItem == 0 {
					break
				}

				a := m.lists[m.selectedList].items[m.selectedItem]
				b := m.lists[m.selectedList].items[m.selectedItem-1]

				m.lists[m.selectedList].items[m.selectedItem] = b
				m.lists[m.selectedList].items[m.selectedItem-1] = a

				m.selectedItem--
			case "shift+right", "L":
				if m.selectedList == len(m.lists)-1 {
					break
				}

				a := m.lists[m.selectedList].items[m.selectedItem]
				m.lists[m.selectedList].items = slices.Delete(m.lists[m.selectedList].items, m.selectedItem, m.selectedItem+1)

				m.selectedList++

				if m.selectedItem > len(m.lists[m.selectedList].items) {
					m.selectedItem = len(m.lists[m.selectedList].items)
				}

				m.lists[m.selectedList].items = slices.Insert(m.lists[m.selectedList].items, m.selectedItem, a)
			case "shift+left", "H":
				if m.selectedList == 0 {
					break
				}

				a := m.lists[m.selectedList].items[m.selectedItem]
				m.lists[m.selectedList].items = slices.Delete(m.lists[m.selectedList].items, m.selectedItem, m.selectedItem+1)

				m.selectedList--

				if m.selectedItem > len(m.lists[m.selectedList].items) {
					m.selectedItem = len(m.lists[m.selectedList].items)
				}

				m.lists[m.selectedList].items = slices.Insert(m.lists[m.selectedList].items, m.selectedItem, a)
			case "o":
				m.selectedItem++
				m.lists[m.selectedList].items = slices.Insert(m.lists[m.selectedList].items, m.selectedItem, "")
				m.input.SetValue(m.lists[m.selectedList].items[m.selectedItem])
				m.input.Focus()
				m.editing = true
			case "O":
				m.lists[m.selectedList].items = slices.Insert(m.lists[m.selectedList].items, m.selectedItem, "")
				m.input.SetValue(m.lists[m.selectedList].items[m.selectedItem])
				m.input.Focus()
				m.editing = true
			case "D":
				m.lists[m.selectedList].items = slices.Delete(m.lists[m.selectedList].items, m.selectedItem, m.selectedItem+1)
				if m.selectedItem > 0 {
					m.selectedItem--
				}
			case "i", "I":
				m.input.SetValue(m.lists[m.selectedList].items[m.selectedItem])
				m.input.Focus()
				m.input.CursorStart()
				m.editing = true
			case "a", "A":
				m.input.SetValue(m.lists[m.selectedList].items[m.selectedItem])
				m.input.Focus()
				m.input.CursorEnd()
				m.editing = true
			case "s", "S":
				m.input.SetValue("")
				m.input.Focus()
				m.editing = true
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	// contentWidth := m.width - border.GetHorizontalBorderSize() - border.GetHorizontalPadding()
	contentHeight := m.height - border.GetVerticalBorderSize() - border.GetVerticalPadding()

	listWidth := (m.width-border.GetHorizontalBorderSize())/len(m.lists) - 1
	todoWidth := listWidth - border.GetHorizontalBorderSize() - border.GetHorizontalPadding()

	lists := make([]string, len(m.lists))
	for li, v := range m.lists {
		title := border.
			BorderBottom(true).
			BorderTop(false).
			BorderRight(false).
			BorderLeft(false).
			Width(listWidth).
			Align(lipgloss.Center).
			Render(titleStyle.Render(v.title))

		titleHeight := countLines(title)

		pages := make([][]string, 1)
		pages[0] = make([]string, 0)
		pageIndex := 0
		pageLen := 0
		selectedPage := 0

		for ii, v := range v.items {
			var content string
			if m.selectedList == li && m.selectedItem == ii {
				if m.editing {
					content = border.Width(todoWidth).Render(m.input.View())
				} else {
					content = selectedItemStyle.Render(border.Width(todoWidth).Render(v))
				}
			} else {
				content = border.Width(todoWidth).Render(v)
			}

			itemHeight := countLines(content)
			if pageLen+itemHeight+titleHeight > contentHeight {
				pages = append(pages, make([]string, 1))
				pageIndex++
				pageLen = 0
			}

			if m.selectedItem == ii {
				selectedPage = pageIndex
			}

			pages[pageIndex] = append(pages[pageIndex], content)
			pageLen += itemHeight
		}

		todos := lipgloss.JoinVertical(lipgloss.Left, pages[selectedPage]...)
		list := lipgloss.JoinVertical(lipgloss.Left, title, todos)

		text := lipgloss.NewStyle().
			Width(listWidth).
			Height(contentHeight).
			Align(lipgloss.Left).
			AlignVertical(lipgloss.Top).
			Render(list)

		lists[li] = border.
			Width(listWidth).
			Render(text)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, lists...)
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
