package board

import (
	tea "github.com/charmbracelet/bubbletea"
	"slices"
	"strings"
	"tuido/commands"
	"tuido/data"
)

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

	pos, err := data.GetPosition(m.workingList().Items, m.selectedItem)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	a.Position = pos
	a.ListId = m.workingList().ListId

	err = m.Provider.UpdateItem(a)

	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, a)

	return m, nil
}

func (m *Model) moveItem(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > m.workingItemsLen()-1 {
		return m, nil
	}

	a := (*m.workingItems())[m.selectedItem]
	aPos := a.Position
	b := (*m.workingItems())[dest]

	a.Position = b.Position
	b.Position = aPos

	err := m.Provider.UpdateItem(a)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	err = m.Provider.UpdateItem(b)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	(*m.workingItems())[m.selectedItem] = b
	(*m.workingItems())[dest] = a

	m.selectedItem = dest

	return m, nil
}

func (m *Model) addItem(dest int) (tea.Model, tea.Cmd) {
	if dest > m.selectedItem && m.workingItemsLen() > 0 {
		m.selectedItem = dest
	}

	if m.selectedItem < 0 {
		m.selectedItem = 0
	}

	pos, err := data.GetPosition(*m.workingItems(), m.selectedItem)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	item := data.Item{
		ListId:   m.workingList().ListId,
		Text:     "",
		Position: pos,
	}

	err = m.Provider.InsertItem(&item)

	// TODO: need to set the position
	*m.workingItems() = slices.Insert(*m.workingItems(), m.selectedItem, item)
	m.input.SetValue(item.Text)
	m.input.Focus()
	m.Editing = true
	m.editField = editItem

	return m, nil
}

func (m *Model) deleteItem() (tea.Model, tea.Cmd) {
	if m.workingItemsLen() == 0 {
		return m, nil
	}

	err := m.Provider.DeleteItem((*m.workingItems())[m.selectedItem])

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	*m.workingItems() = slices.Delete(*m.workingItems(), m.selectedItem, m.selectedItem+1)
	if m.selectedItem > 0 {
		m.selectedItem--
	}

	return m, nil
}

func (m *Model) deleteList() (tea.Model, tea.Cmd) {
	if len(m.Lists) == 0 {
		return m, nil
	}

	err := m.Provider.DeleteList(*m.workingList())

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	m.Lists = slices.Delete(m.Lists, m.selectedList, m.selectedList+1)
	if m.selectedList > 0 {
		m.selectedList--
	}

	return m, nil
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
	m.input.SetValue(m.workingList().Name)
	m.input.Focus()
	m.Editing = true
	m.editField = editTitle
}

func (m *Model) editBoard() {
	m.input.SetValue(m.Board.Name)
	m.input.Focus()
	m.Editing = true
	m.editField = editBoard
}

func (m *Model) addList() (tea.Model, tea.Cmd) {
	pos, err := data.GetPosition(m.Lists, len(m.Lists))

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	list := data.List{
		BoardId:  m.Board.BoardId,
		Name:     "",
		Items:    make([]data.Item, 0),
		Position: pos,
	}

	m.Provider.InsertList(&list)

	m.Lists = append(m.Lists, list)
	m.selectedList = len(m.Lists) - 1
	m.selectedItem = 0
	m.Editing = true
	m.editField = editTitle
	m.input.SetValue(m.workingList().Name)
	m.input.Focus()

	return m, nil
}

func (m *Model) moveList(dest int) (tea.Model, tea.Cmd) {
	if dest < 0 || dest > len(m.Lists)-1 {
		return m, nil
	}

	a := m.Lists[m.selectedList]
	aPos := a.Position
	b := m.Lists[dest]

	a.Position = b.Position
	b.Position = aPos

	err := m.Provider.UpdateList(a)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	err = m.Provider.UpdateList(b)

	if err != nil {
		return m, commands.ErrorCmd(err)
	}

	m.Lists[m.selectedList] = b
	m.Lists[dest] = a

	m.selectedList = dest

	return m, nil
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

			err := m.Provider.UpdateItem((*m.workingItems())[m.selectedItem])

			if err != nil {
				return m, commands.ErrorCmd(err)
			}

			if string == "enter" {
				return m.addItem(m.selectedItem + 1)
			}
		case editTitle:
			m.workingList().Name = m.input.Value()

			err := m.Provider.UpdateList(*m.workingList())

			if err != nil {
				return m, commands.ErrorCmd(err)
			}
		case editBoard:
			m.Board.Name = m.input.Value()

			err := m.Provider.UpdateBoard(m.Board)

			if err != nil {
				return m, commands.ErrorCmd(err)
			}
		}

		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		return m, cmd
	}
}
