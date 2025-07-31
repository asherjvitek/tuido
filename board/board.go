package board

import (
	"slices"
	"strings"
	"tuido/commands"
	"tuido/data"
	"tuido/style"
	"tuido/util"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type List struct {
	Id             int
	Title          string
	Position       int
	Items          []string
}

type Model struct {
	Id           int
	Name         string
	Lists        []data.List
	Position     int
	Width        int
	Height       int
	selectedList int
	selectedItem int
	Editing      bool
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

func New(id int) Model {
	return Model{
		Id:   id,
		Name: "New Board",
		Lists: []data.List{
			{
				Title: "New List",
				Items: []data.Item{},
			},
		},

		selectedList: 0,
		selectedItem: 0,
		input:        util.GetTextInput(),
	}
}

// func (m *Model) Setup() {
// 	m.input = util.GetTextInput()
// }

func (m Model) workingList() *data.List {
	return &m.Lists[m.selectedList]
}

func (m Model) workingItems() *[]data.Item {
	return &m.Lists[m.selectedList].Items
}

func (m Model) workingItemsLen() int {
	return len(m.Lists[m.selectedList].Items)
}

func (m *Model) moveItemToList(dest int) (tea.Model, tea.Cmd) {
	if dest > len(m.Lists)-1 || dest < 0 || m.workingItemsLen() == 0 {
		return m, nil
	}

	a := (*m.workingItems())[m.selectedItem]
	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)

	m.selectedList = dest

	if m.selectedItem > m.workingItemsLen() {
		m.selectedItem = m.workingItemsLen()
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, a)

	return m, commands.SaveDataMsg
}

func (m *Model) moveItem(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > m.workingItemsLen()-1 {
		return m, nil
	}

	a := (*m.workingItems())[m.selectedItem]
	b := (*m.workingItems())[dest]

	(*m.workingItems())[m.selectedItem] = b
	(*m.workingItems())[dest] = a

	m.selectedItem = dest

	return m, commands.SaveDataMsg
}

func (m *Model) addItem(dest int) (tea.Model, tea.Cmd) {
	if dest > m.selectedItem && m.workingItemsLen() > 0 {
		m.selectedItem = dest
	}

	if m.selectedItem < 0 {
		m.selectedItem = 0
	}

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, data.Item{Id: 0, Text: "", Position: dest, ListId: m.workingList().Id})
	m.input.SetValue((*m.workingItems())[m.selectedItem].Text)
	m.input.Focus()
	m.Editing = true
	m.editField = editItem

	return m, nil
}

func (m *Model) deleteItem() (tea.Model, tea.Cmd) {
	if m.workingItemsLen() == 0 {
		return m, nil
	}

	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)
	if m.selectedItem > 0 {
		m.selectedItem--
	}

	return m, commands.SaveDataMsg
}

func (m *Model) deleteList() (tea.Model, tea.Cmd) {
	if len(m.Lists) == 0 {
		return m, nil
	}

	list := m.workingList()

	data.DeleteList(list.Id);

	m.Lists = slices.Delete(m.Lists, m.selectedList, m.selectedList+1)
	if m.selectedList > 0 {
		m.selectedList--
	}

	return m, commands.SaveDataMsg
}

func (m *Model) editItem(cursorLocation EditType) {
	if m.workingItemsLen() == 0 {
		return
	}

	m.input.Focus()

	switch cursorLocation {
	case EditTypeStart:
		m.input.SetValue((*m.workingItems())[m.selectedItem].Text)
		m.input.CursorStart()
	case EditTypeEnd:
		m.input.SetValue((*m.workingItems())[m.selectedItem].Text)
		m.input.CursorEnd()
	case EditTypeClear:
		m.input.SetValue("")
	}

	m.Editing = true
	m.editField = editItem
}

func (m *Model) editTitle() {
	m.input.SetValue(m.workingList().Title)
	m.input.Focus()
	m.Editing = true
	m.editField = editTitle
}

func (m *Model) editBoard() {
	m.input.SetValue(m.Name)
	m.input.Focus()
	m.Editing = true
	m.editField = editBoard
}

func (m *Model) addList() {
	m.Lists = append(m.Lists, data.List{
		Title:          "",
		Items:          make([]data.Item, 0),
	})
	m.selectedList = len(m.Lists) - 1
	m.selectedItem = 0
	m.Editing = true
	m.editField = editTitle
	m.input.SetValue(m.workingList().Title)
	m.input.Focus()
}

func (m *Model) moveList(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > len(m.Lists)-1 {
		return m, nil
	}

	a := m.Lists[m.selectedList]
	b := m.Lists[dest]

	m.Lists[m.selectedList] = b
	m.Lists[dest] = a

	m.selectedList = dest

	return m, commands.SaveDataMsg
}

func (m *Model) navigate(itemDest int, listDest int) {

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

	m.selectedItem = max(m.selectedItem, 0)
}

func (m Model) handleEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	string := msg.String()
	switch string {
	case "esc", "enter":
		m.input.Blur()
		m.Editing = false

		switch m.editField {
		case editItem:
			value := m.input.Value()
			if len(strings.Trim(value, " ")) == 0 {
				return m.deleteItem()
			}

			(*m.workingItems())[m.selectedItem].Text = m.input.Value()

			if string == "enter" {
				return m.addItem(m.selectedItem + 1)
			}
		case editTitle:
			m.workingList().Title = m.input.Value()
		case editBoard:
			m.Name = m.input.Value()
		}

		return m, commands.SaveDataMsg
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		return m, cmd
	}
}

var (
	boardNameStyle = lg.NewStyle().
			Foreground(lg.Color(style.Blue)).
			BorderForeground(lg.Color(style.Blue)).
			Height(1).
			Margin(1, 1, 0, 1).
			Border(lg.ThickBorder(), false, false, true, false)

	border = lg.NewStyle().
		Border(lg.RoundedBorder())

	titleStyle = lg.NewStyle().
			Foreground(lg.Color(style.Purple))

	titleContainerStyle = lg.NewStyle().
				Border(lg.NormalBorder(), false, false, true, false).
				Align(lg.Center)

	listStyle = border.
			MaxWidth(40).
			Margin(0, 1, 1, 1)

	selectedListStyle = listStyle.
				BorderForeground(lg.Color(style.Green))

	selectedItemStyle = border.
				BorderForeground(lg.Color(style.Yellow))
)

type boardInit struct {
	lists []data.List
}

func (m Model) Init() tea.Cmd {
	lists := data.GetLists(m.Id)
	return func () tea.Msg {
		return boardInit { lists: lists }
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Editing {
			return m.handleEditing(msg)
		}

		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit

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
		case "d":
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
		case "D":
			return m.deleteList()

		//Return to boards
		case "b":
			return m, commands.ChangeScreenBoardsCmd(m.Id)
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case boardInit:
		m.Lists = msg.lists
		m.input = util.GetTextInput()
	}

	return m, nil
}

func (m Model) View() string {
	listLen := len(m.Lists)

	boardNameStyle = boardNameStyle.Width(m.Width - boardNameStyle.GetHorizontalFrameSize())

	if listLen == 0 {
		if m.Editing && m.editField == editBoard {
			return boardNameStyle.Render(m.input.View())
		} else {
			return boardNameStyle.Render(m.Name)
		}
	}

	contentHeight := m.Height - listStyle.GetHorizontalFrameSize() - boardNameStyle.GetHeight() - boardNameStyle.GetVerticalFrameSize()

	listMaxWidth := listStyle.GetMaxWidth() - listStyle.GetHorizontalFrameSize()
	listWidth := min((m.Width-listStyle.GetHorizontalFrameSize()*listLen)/listLen, listMaxWidth)

	todoWidth := listWidth - border.GetHorizontalFrameSize()

	lists := make([]string, listLen)
	for li, v := range m.Lists {
		styledTitle := titleStyle.Render(v.Title)

		if m.Editing && m.editField == editTitle && li == m.selectedList {
			styledTitle = titleStyle.Render(m.input.View())
		}

		title := titleContainerStyle.
			Width(listWidth).
			Render(styledTitle)

		titleHeight := util.CountLines(title)

		pages := make([][]string, 1)
		pages[0] = make([]string, 0)
		pageIndex := 0
		pageLen := 0
		selectedPage := 0

		for ii, v := range v.Items {
			var content string
			if m.selectedList == li && m.selectedItem == ii {
				if m.Editing && m.editField == editItem {
					content = selectedItemStyle.Width(todoWidth).Render(m.input.View())
				} else {
					content = selectedItemStyle.Width(todoWidth).Render(v.Text)
				}
			} else {
				content = border.Width(todoWidth).Render(v.Text)
			}

			itemHeight := util.CountLines(content)
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
	if m.Editing && m.editField == editBoard {
		board = boardNameStyle.Render(m.input.View())
	} else {
		board = boardNameStyle.Render(m.Name)
	}

	return lg.JoinVertical(lg.Left, board, lg.JoinHorizontal(lg.Left, lists...))
}
