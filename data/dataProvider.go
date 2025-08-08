package data

type Provider interface {
	Boards() ([]Board, error)
	Lists(boardId int) ([]List, error)
	InsertBoard(board *Board) error
	UpdateBoard(board Board) error
	DeleteBoard(boardId int) error
	InsertList(list *List) error
	DeleteList(list List) error
	UpdateList(list List) error
	InsertItem(item *Item) error
	DeleteItem(item Item) error
	UpdateItem(item Item) error
}
