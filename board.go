package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"slices"
)

type list struct {
	Title          string
	Items          []string
	scrollposition int
}

type board struct {
	Id           int
	Name         string
	Lists        []list
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
	return &m.Lists[m.selectedList]
}

func (m board) workingItems() *[]string {
	return &m.Lists[m.selectedList].Items
}

func (m board) workingItemsLen() int {
	return len(m.Lists[m.selectedList].Items)
}

func (m *board) moveItemToList(dest int) (tea.Model, tea.Cmd) {
	if dest > len(m.Lists)-1 || dest < 0 {
		return m, nil
	}

	a := (*m.workingItems())[m.selectedItem]
	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)

	m.selectedList = dest

	if m.selectedItem > m.workingItemsLen() {
		m.selectedItem = m.workingItemsLen()
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, a)

	return m, saveBoard
}

func (m *board) moveItem(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > m.workingItemsLen() {
		return m, nil
	}

	a := (*m.workingItems())[m.selectedItem]
	b := (*m.workingItems())[dest]

	(*m.workingItems())[m.selectedItem] = b
	(*m.workingItems())[dest] = a

	m.selectedItem = dest

	return m, saveBoard
}

func (m *board) addItem(dest int) (tea.Model, tea.Cmd) {
	if dest > m.selectedItem && m.workingItemsLen() > 0 {
		m.selectedItem = dest
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, "")
	m.input.SetValue((*m.workingItems())[m.selectedItem])
	m.input.Focus()
	m.editing = true
	m.editField = editItem

	return m, saveBoard
}

func (m *board) deleteItem() (tea.Model, tea.Cmd) {
	if m.workingItemsLen() == 0 {
		return m, nil
	}

	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)
	if m.selectedItem > 0 {
		m.selectedItem--
	}

	return m, saveBoard
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
	m.input.SetValue(m.workingList().Title)
	m.input.Focus()
	m.editing = true
	m.editField = editTitle
}

func (m *board) editBoard() {
	m.input.SetValue(m.Name)
	m.input.Focus()
	m.editing = true
	m.editField = editBoard
}

func (m *board) addList() {
	m.Lists = append(m.Lists, list{
		Title:          "",
		Items:          make([]string, 0),
		scrollposition: 0,
	})
	m.selectedList = len(m.Lists) - 1
	m.selectedItem = 0
	m.editing = true
	m.editField = editTitle
	m.input.SetValue(m.workingList().Title)
	m.input.Focus()
}

func (m *board) moveList(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > len(m.Lists) {
		return m, nil
	}

	a := m.Lists[m.selectedList]
	b := m.Lists[dest]

	m.Lists[m.selectedList] = b
	m.Lists[dest] = a

	m.selectedList = dest

	return m, saveBoard
}

func (m *board) navigate(itemDest int, listDest int) {

	if itemDest != m.selectedItem {
		if itemDest > m.workingItemsLen()-1 || itemDest < 0 {
			return
		}

		m.selectedItem = itemDest
	}

	if listDest != m.selectedList {
		if listDest > len(m.Lists)-1 || listDest < 0 {
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
			Foreground(lg.Color(Blue)).
			BorderForeground(lg.Color(Blue)).
			Height(1).
			Margin(1, 1, 0, 1).
			Border(lg.ThickBorder(), false, false, true, false)

	border = lg.NewStyle().
		Border(lg.RoundedBorder())

	titleStyle = lg.NewStyle().
			Foreground(lg.Color(Purple))

	titleContainerStyle = lg.NewStyle().
				Border(lg.NormalBorder(), false, false, true, false).
				Align(lg.Center)

	listStyle = border.
			MaxWidth(40).
			Margin(0, 1, 1, 1)

	selectedListStyle = listStyle.
				BorderForeground(lg.Color(Green))

	selectedItemStyle = border.
				BorderForeground(lg.Color(Yellow))
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
					m.workingList().Title = m.input.Value()
				case editBoard:
					m.Name = m.input.Value()
				}
				m.editing = false

				return m, saveBoard
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
				m.navigate(m.selectedItem, len(m.Lists)-1)
			case "home":
				m.navigate(0, m.selectedList)
			case "end":
				m.navigate(m.workingItemsLen()-1, m.selectedList)

			// Moving things
			case "shift+down", "J":
				return m.moveItem(m.selectedItem + 1)
			case "shift+up", "K":
				return m.moveItem(m.selectedItem - 1)
			case "shift+right", "L":
				return m.moveItemToList(m.selectedList + 1)
			case "shift+left", "H":
				return m.moveItemToList(m.selectedList - 1)
			case "alt+right", "alt+l", "alt+L":
				return m.moveList(m.selectedList + 1)
			case "alt+left", "alt+h", "alt+H":
				return m.moveList(m.selectedList - 1)

			//Editing
			case "o":
				m.addItem(m.selectedItem + 1)
			case "O":
				m.addItem(m.selectedItem)
			case "D":
				return m.deleteItem()
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
	listLen := len(m.Lists)

	boardNameStyle = boardNameStyle.Width(m.width - boardNameStyle.GetHorizontalFrameSize())

	if listLen == 0 {
		if m.editing && m.editField == editBoard {
			return boardNameStyle.Render(m.input.View())
		} else {
			return boardNameStyle.Render(m.Name)
		}
	}

	contentHeight := m.height - listStyle.GetHorizontalFrameSize() - boardNameStyle.GetHeight() - boardNameStyle.GetVerticalFrameSize()

	listMaxWidth := listStyle.GetMaxWidth() - listStyle.GetHorizontalFrameSize()
	listWidth := min((m.width-listStyle.GetHorizontalFrameSize()*listLen)/listLen, listMaxWidth)

	todoWidth := listWidth - border.GetHorizontalFrameSize()

	lists := make([]string, listLen)
	for li, v := range m.Lists {
		styledTitle := titleStyle.Render(v.Title)

		if m.editing && m.editField == editTitle && li == m.selectedList {
			styledTitle = titleStyle.Render(m.input.View())
		}

		title := titleContainerStyle.
			Width(listWidth).
			Render(styledTitle)

		titleHeight := countLines(title)

		pages := make([][]string, 1)
		pages[0] = make([]string, 0)
		pageIndex := 0
		pageLen := 0
		selectedPage := 0

		for ii, v := range v.Items {
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
		board = boardNameStyle.Render(m.Name)
	}

	return lg.JoinVertical(lg.Left, board, lg.JoinHorizontal(lg.Left, lists...))
}
