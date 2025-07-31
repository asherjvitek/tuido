package data

type Board struct {
	Id       int
	Name     string
	Position int
	Lists    []List
}

type List struct {
	Id       int
	BoardId  int
	Title    string
	Position int
	Items    []Item
}

type Item struct {
	Id       int
	ListId   int
	Text     string
	Position int
}
