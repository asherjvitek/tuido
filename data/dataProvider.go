package data

type Provider interface {
	Boards() ([]Board, error)
	Lists(boardId int) ([]List, error)

	InsertBoard(board *Board) error
	UpdateBoard(board Board) error
	DeleteBoard(boardId int) error

	InsertList(list *List) error
	UpdateList(list List) error
	DeleteList(listId int) error

	InsertItem(item *Item) error
	UpdateItem(item Item) error
	DeleteItem(itemId int) error
}
