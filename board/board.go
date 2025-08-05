package board

import (
	"tuido/commands"
	"tuido/config"
	"tuido/data"
	"tuido/style"
	"tuido/util"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type List struct {
	Title string
	Items []string
}

type Model struct {
	Board        data.Board
	Lists        []data.List
	width        int
	height       int
	selectedList int
	selectedItem int
	// I do not like this.... not sure if we should remove..
	Editing   bool
	editField editField
	input     textinput.Model
	Config    config.Config
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

type initMsg []data.List

func (m Model) Init() tea.Cmd {
	lists, err := data.Lists(m.Board.BoardId)

	if err != nil {
		return commands.ErrorCmd(err)
	}

	return func() tea.Msg {
		return initMsg(lists)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Editing {
			return m.handleEditing(msg)
		}

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
		case "alt+D":
			return m.deleteList()

		//Return to boards
		case "b":
			return m, commands.ChangeScreenBoardsCmd(m.Board.BoardId)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case initMsg:
		m.Lists = msg
		m.input = util.GetTextInput()
	}

	return m, nil
}

func (m Model) View() string {
	listLen := len(m.Lists)

	boardNameStyle = boardNameStyle.Width(m.width - boardNameStyle.GetHorizontalFrameSize())

	if listLen == 0 {
		if m.Editing && m.editField == editBoard {
			return boardNameStyle.Render(m.input.View())
		} else {
			return boardNameStyle.Render(m.Board.Name)
		}
	}

	contentHeight := m.height - listStyle.GetHorizontalFrameSize() - boardNameStyle.GetHeight() - boardNameStyle.GetVerticalFrameSize()

	listMaxWidth := listStyle.GetMaxWidth() - listStyle.GetHorizontalFrameSize()
	listWidth := min((m.width-listStyle.GetHorizontalFrameSize()*listLen)/listLen, listMaxWidth)

	todoWidth := listWidth - border.GetHorizontalFrameSize()

	lists := make([]string, listLen)
	for li, v := range m.Lists {
		styledTitle := titleStyle.Render(v.Name)

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
		board = boardNameStyle.Render(m.Board.Name)
	}

	return lg.JoinVertical(lg.Left, board, lg.JoinHorizontal(lg.Left, lists...))
}
