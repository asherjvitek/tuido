package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"slices"
)

type list struct {
	title          string
	items          []string
	scrollposition int
}

type board struct {
	id           int
	name         string
	lists        []list
	width        int
	height       int
	selectedList int
	selectedItem int
	editing      bool
	editField    editField
	input        textinput.Model
}

type editField int

const (
	editItem editField = iota
	editTitle
	editBoard
)

type EditType int

const (
	EditTypeStart EditType = iota
	EditTypeEnd
	EditTypeClear
)

func (m board) workingList() *list {
	return &m.lists[m.selectedList]
}

func (m board) workingItems() *[]string {
	return &m.lists[m.selectedList].items
}

func (m board) workingItemsLen() int {
	return len(m.lists[m.selectedList].items)
}

func (m *board) moveItemToList(dest int) {
	if dest > len(m.lists)-1 || dest < 0 {
		return
	}

	a := (*m.workingItems())[m.selectedItem]
	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)

	m.selectedList = dest

	if m.selectedItem > m.workingItemsLen() {
		m.selectedItem = m.workingItemsLen()
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, a)
}

func (m *board) moveItem(dest int) {
	if dest < 0 || dest > m.workingItemsLen() {
		return
	}

	a := (*m.workingItems())[m.selectedItem]
	b := (*m.workingItems())[dest]

	(*m.workingItems())[m.selectedItem] = b
	(*m.workingItems())[dest] = a

	m.selectedItem = dest
}

func (m *board) addItem(dest int) {
	if dest > m.selectedItem && m.workingItemsLen() > 0 {
		m.selectedItem = dest
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, "")
	m.input.SetValue((*m.workingItems())[m.selectedItem])
	m.input.Focus()
	m.editing = true
	m.editField = editItem
}

func (m *board) deleteItem() {
	if m.workingItemsLen() == 0 {
		return
	}

	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)
	if m.selectedItem > 0 {
		m.selectedItem--
	}
}

func (m *board) editItem(cursorLocation EditType) {
	if m.workingItemsLen() == 0 {
		return
	}

	m.input.Focus()

	switch cursorLocation {
	case EditTypeStart:
		m.input.SetValue((*m.workingItems())[m.selectedItem])
		m.input.CursorStart()
	case EditTypeEnd:
		m.input.SetValue((*m.workingItems())[m.selectedItem])
		m.input.CursorEnd()
	case EditTypeClear:
		m.input.SetValue("")
	}

	m.editing = true
	m.editField = editItem
}

func (m *board) editTitle() {
	m.input.SetValue(m.workingList().title)
	m.input.Focus()
	m.editing = true
	m.editField = editTitle
}

func (m *board) editBoard() {
	m.input.SetValue(m.name)
	m.input.Focus()
	m.editing = true
	m.editField = editBoard
}

func (m *board) addList() {
	m.lists = append(m.lists, list{
		title:          "",
		items:          make([]string, 0),
		scrollposition: 0,
	})
	m.selectedList = len(m.lists) - 1
	m.selectedItem = 0
	m.editing = true
	m.editField = editTitle
	m.input.SetValue(m.workingList().title)
	m.input.Focus()
}

func (m *board) moveList(dest int) {
	if dest < 0 || dest > len(m.lists) {
		return
	}

	a := m.lists[m.selectedList]
	b := m.lists[dest]

	m.lists[m.selectedList] = b
	m.lists[dest] = a

	m.selectedList = dest
}

func (m *board) navigate(itemDest int, listDest int) {

	if itemDest != m.selectedItem {
		if itemDest > m.workingItemsLen()-1 || itemDest < 0 {
			return
		}

		m.selectedItem = itemDest
	}

	if listDest != m.selectedList {
		if listDest > len(m.lists)-1 || listDest < 0 {
			return
		}

		m.selectedList = listDest
		if m.selectedItem > m.workingItemsLen()-1 {
			m.selectedItem = m.workingItemsLen() - 1
		}
	}
}

var (
	boardNameStyle = lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("5")).
			Height(1).
			Margin(1, 1, 1, 1)

	border = lg.NewStyle().
		Border(lg.RoundedBorder())

	titleStyle = lg.NewStyle().
			Foreground(lg.Color("5"))

	listStyle = border.
			Margin(0, 1, 1, 1)

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
			case "esc", "enter":
				m.input.Blur()
				switch m.editField {
				case editItem:
					(*m.workingItems())[m.selectedItem] = m.input.Value()
				case editTitle:
					m.workingList().title = m.input.Value()
				case editBoard:
					m.name = m.input.Value()
				}
				m.editing = false

				return m, nil
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)

				return m, cmd
			}
		} else {
			switch msg.String() {

			//Navigation
			case "down", "j":
				m.navigate(m.selectedItem+1, m.selectedList)
			case "up", "k":
				m.navigate(m.selectedItem-1, m.selectedList)
			case "right", "l":
				m.navigate(m.selectedItem, m.selectedList+1)
			case "left", "h":
				m.navigate(m.selectedItem, m.selectedList-1)
			case "ctrl+home":
				m.navigate(m.selectedItem, 0)
			case "ctrl+end":
				m.navigate(m.selectedItem, len(m.lists)-1)
			case "home":
				m.navigate(0, m.selectedList)
			case "end":
				m.navigate(m.workingItemsLen()-1, m.selectedList)

			// Moving things
			case "shift+down", "J":
				m.moveItem(m.selectedItem + 1)
			case "shift+up", "K":
				m.moveItem(m.selectedItem - 1)
			case "shift+right", "L":
				m.moveItemToList(m.selectedList + 1)
			case "shift+left", "H":
				m.moveItemToList(m.selectedList - 1)
			case "alt+right", "alt+l", "alt+L":
				m.moveList(m.selectedList + 1)
			case "alt+left", "alt+h", "alt+H":
				m.moveList(m.selectedList - 1)

			//Editing
			case "o":
				m.addItem(m.selectedItem + 1)
			case "O":
				m.addItem(m.selectedItem)
			case "D":
				m.deleteItem()
			case "i", "I":
				m.editItem(EditTypeStart)
			case "a", "A":
				m.editItem(EditTypeEnd)
			case "s", "S":
				m.editItem(EditTypeClear)
			case "t", "T":
				m.editTitle()
			case "B":
				m.editBoard()
			case "N":
				m.addList()

			//Return to boards
			case "b":
				return m, func() tea.Msg {
					return changeScreenBoards{}
				}
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

	if listLen == 0 {
		if m.editing && m.editField == editBoard {
			return boardNameStyle.Render(m.input.View())
		} else {
			return boardNameStyle.Render(m.name)
		}
	}

	contentHeight := m.height - listStyle.GetVerticalBorderSize() - listStyle.GetVerticalMargins() - boardNameStyle.GetHeight() - boardNameStyle.GetVerticalMargins()
	listWidth := (m.width - (listStyle.GetHorizontalBorderSize()+listStyle.GetHorizontalMargins())*listLen) / listLen
	todoWidth := listWidth - border.GetHorizontalBorderSize() - border.GetHorizontalPadding()

	lists := make([]string, listLen)
	for li, v := range m.lists {
		styledTitle := titleStyle.Render(v.title)

		if m.editing && m.editField == editTitle && li == m.selectedList {
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
				if m.editing && m.editField == editItem {
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

	board := ""
	if m.editing && m.editField == editBoard {
		board = boardNameStyle.Render(m.input.View())
	} else {
		board = boardNameStyle.Render(m.name)
	}

	return lg.JoinVertical(lg.Left, board, lg.JoinHorizontal(lg.Left, lists...))
}
