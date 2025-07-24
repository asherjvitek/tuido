package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"os"
	"slices"
)

type list struct {
	title          string
	items          []string
	scrollposition int
}

type board struct {
	lists        []list
	width        int
	height       int
	selectedList int
	selectedItem int
	editing      bool
	editingTitle bool
	input        textinput.Model
}


type EditType int

const (
	EditTypeStart EditType = iota
	EditTypeEnd
	EditTypeClear
)

func (m board) WorkingList() *list {
	return &m.lists[m.selectedList]
}

func (m board) WorkingItems() *[]string {
	return &m.lists[m.selectedList].items
}

func (m board) WorkingItemsLen() int {
	return len(m.lists[m.selectedList].items)
}

func (m *board) MoveItemToList(dest int) {
	if dest > len(m.lists)-1 || dest < 0 {
		return
	}

	a := (*m.WorkingItems())[m.selectedItem]
	*m.WorkingItems() = slices.Delete(*m.WorkingItems(), m.selectedItem, m.selectedItem+1)

	m.selectedList = dest

	if m.selectedItem > m.WorkingItemsLen() {
		m.selectedItem = m.WorkingItemsLen()
	}

	*m.WorkingItems() = slices.Insert(*m.WorkingItems(), m.selectedItem, a)
}

func (m *board) MoveItem(dest int) {
	if dest < 0 || dest > m.WorkingItemsLen() {
		return
	}

	a := (*m.WorkingItems())[m.selectedItem]
	b := (*m.WorkingItems())[dest]

	(*m.WorkingItems())[m.selectedItem] = b
	(*m.WorkingItems())[dest] = a

	m.selectedItem = dest
}

func (m *board) AddItem(dest int) {
	if dest > m.selectedItem && m.WorkingItemsLen() > 0 {
		m.selectedItem = dest
	}

	*m.WorkingItems() = slices.Insert(*m.WorkingItems(), m.selectedItem, "")
	m.input.SetValue((*m.WorkingItems())[m.selectedItem])
	m.input.Focus()
	m.editing = true
}

func (m *board) DeleteItem() {
	if m.WorkingItemsLen() == 0 {
		return
	}

	*m.WorkingItems() = slices.Delete(*m.WorkingItems(), m.selectedItem, m.selectedItem+1)
	if m.selectedItem > 0 {
		m.selectedItem--
	}
}

func (m *board) EditItem(cursorLocation EditType) {
	if m.WorkingItemsLen() == 0 {
		return
	}

	m.input.Focus()

	switch cursorLocation {
	case EditTypeStart:
		m.input.SetValue((*m.WorkingItems())[m.selectedItem])
		m.input.CursorStart()
	case EditTypeEnd:
		m.input.SetValue((*m.WorkingItems())[m.selectedItem])
		m.input.CursorEnd()
	case EditTypeClear:
		m.input.SetValue("")
	}

	m.editing = true
}

func (m *board) EditTitle() {
	m.input.SetValue(m.WorkingList().title)
	m.input.Focus()
	m.editingTitle = true
}

func (m *board) AddList() {
	m.lists = append(m.lists, list{
		title:          "",
		items:          make([]string, 0),
		scrollposition: 0,
	})
	m.selectedList = len(m.lists) - 1
	m.selectedItem = 0
	m.editingTitle = true
	m.input.SetValue(m.WorkingList().title)
	m.input.Focus()
}

func (m *board) MoveList(dest int) {
	if dest < 0 || dest > len(m.lists) {
		return
	}

	a := m.lists[m.selectedList]
	b := m.lists[dest]

	m.lists[m.selectedList] = b
	m.lists[dest] = a

	m.selectedList = dest
}

func (m *board) Navigate(itemDest int, listDest int) {

	if itemDest != m.selectedItem {
		if itemDest > m.WorkingItemsLen()-1 || itemDest < 0 {
			return
		}

		m.selectedItem = itemDest
	}

	if listDest != m.selectedList {
		if listDest > len(m.lists)-1 || listDest < 0 {
			return
		}

		m.selectedList = listDest
		if m.selectedItem > m.WorkingItemsLen()-1 {
			m.selectedItem = m.WorkingItemsLen() - 1
		}
	}
}

var (
	border = lg.NewStyle().
		Border(lg.RoundedBorder())

	titleStyle = lg.NewStyle().
			Foreground(lg.Color("5"))

	listStyle = border.
			Margin(1, 1, 1, 1)

	selectedListStyle = listStyle.
				BorderForeground(lg.Color("2"))

	selectedItemStyle = border.
				BorderForeground(lg.Color("11"))
)

func (m board) Init() tea.Cmd {
	return nil
}

func (m board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.input.Blur()
				(*m.WorkingItems())[m.selectedItem] = m.input.Value()
				m.editing = false
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)

				return m, cmd
			}
		} else if m.editingTitle {
			switch msg.String() {
			case "esc":
				m.input.Blur()
				m.WorkingList().title = m.input.Value()
				m.editingTitle = false
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)

				return m, cmd
			}
		} else {
			switch msg.String() {
			case tea.KeyCtrlC.String(), "q":
				return m, tea.Quit

			//Navigation
			case "down", "j":
				m.Navigate(m.selectedItem+1, m.selectedList)
			case "up", "k":
				m.Navigate(m.selectedItem-1, m.selectedList)
			case "right", "l":
				m.Navigate(m.selectedItem, m.selectedList+1)
			case "left", "h":
				m.Navigate(m.selectedItem, m.selectedList-1)
			case "ctrl+home":
				m.Navigate(m.selectedItem, 0)
			case "ctrl+end":
				m.Navigate(m.selectedItem, len(m.lists)-1)
			case "home":
				m.Navigate(0, m.selectedList)
			case "end":
				m.Navigate(m.WorkingItemsLen()-1, m.selectedList)

			// Moving things
			case "shift+down", "J":
				m.MoveItem(m.selectedItem + 1)
			case "shift+up", "K":
				m.MoveItem(m.selectedItem - 1)
			case "shift+right", "L":
				m.MoveItemToList(m.selectedList + 1)
			case "shift+left", "H":
				m.MoveItemToList(m.selectedList - 1)
			case "alt+right", "alt+l", "alt+L":
				m.MoveList(m.selectedList + 1)
			case "alt+left", "alt+h", "alt+H":
				m.MoveList(m.selectedList - 1)

			//Editing
			case "o":
				m.AddItem(m.selectedItem + 1)
			case "O":
				m.AddItem(m.selectedItem)
			case "D":
				m.DeleteItem()
			case "i", "I":
				m.EditItem(EditTypeStart)
			case "a", "A":
				m.EditItem(EditTypeEnd)
			case "s", "S":
				m.EditItem(EditTypeClear)
			case "t":
				m.EditTitle()
			case "N":
				m.AddList()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m board) View() string {
	listLen := len(m.lists)

	contentHeight := m.height - listStyle.GetVerticalBorderSize() - listStyle.GetVerticalMargins()
	listWidth := (m.width - (listStyle.GetHorizontalBorderSize()+listStyle.GetHorizontalMargins())*listLen) / listLen
	todoWidth := listWidth - border.GetHorizontalBorderSize() - border.GetHorizontalPadding()

	lists := make([]string, listLen)
	for li, v := range m.lists {
		styledTitle := titleStyle.Render(v.title)

		if m.editingTitle && li == m.selectedList {
			styledTitle = titleStyle.Render(m.input.View())
		}

		title := border.
			BorderBottom(true).
			BorderTop(false).
			BorderRight(false).
			BorderLeft(false).
			Width(listWidth).
			Align(lg.Center).
			Render(styledTitle)

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
					content = selectedItemStyle.Width(todoWidth).Render(m.input.View())
				} else {
					content = selectedItemStyle.Width(todoWidth).Render(v)
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

		todos := lg.JoinVertical(lg.Left, pages[selectedPage]...)
		list := lg.JoinVertical(lg.Left, title, todos)

		text := lg.NewStyle().
			Width(listWidth).
			Height(contentHeight).
			Align(lg.Left).
			AlignVertical(lg.Top).
			Render(list)

		ls := listStyle
		if m.selectedList == li {
			ls = selectedListStyle
		}

		lists[li] = ls.
			Width(listWidth).
			Render(text)
	}

	return lg.JoinHorizontal(lg.Left, lists...)
}

func main() {
	model := board{
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
