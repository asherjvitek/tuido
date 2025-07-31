package data

import (
	"database/sql"
	"fmt"
	"time"
	"tuido/util"

	_ "modernc.org/sqlite"
)

const sqliteFileName = "tuido.db"

func connect() (*sql.DB, error) {
	var db *sql.DB

	dbUrl, err := getSavePath(sqliteFileName)
	if err != nil {
		return nil, err
	}

	db, err = sql.Open("sqlite", dbUrl)
	if err != nil {
		return db, fmt.Errorf("error opening %s: %w", dbUrl, err)
	}

	db.SetConnMaxIdleTime(9 * time.Second)

	return db, nil
}

func InitDatabase() error {

	// Get database URL and auth token from environment variables

	db, err := connect()

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS board (
    board_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    position INTEGER
);

CREATE TABLE IF NOT EXISTS list (
    list_id INTEGER PRIMARY KEY AUTOINCREMENT,
    board_id INTEGER,
    name TEXT,
    position INTEGER,
    FOREIGN KEY(board_id) REFERENCES board(board_id)
);

CREATE TABLE IF NOT EXISTS item (
    item_id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_id INTEGER,
    text TEXT,
    position INTEGER,
    FOREIGN KEY(list_id) REFERENCES list(list_id)
);`)

	if err != nil {
		return fmt.Errorf("error creating the database schema: %w", err)
	}

	return nil
}

func GetBoards() ([]Board, error) {
	boards := make([]Board, 0)
	db, err := connect()

	if err != nil {
		return boards, err
	}

	rows, err := db.Query("SELECT board_id, name, position FROM board;")

	if err != nil {
		return boards, err
	}

	for rows.Next() {
		var board Board
		if err := rows.Scan(&board.Id, &board.Name, &board.Position); err != nil {
			return boards, err
		}

		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return boards, err
	}

	return boards, nil
}

func UpdateBoard(board Board) error {
	db, err := connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
UPDATE board
SET 
	name = ?,
	position = ?
WHERE board_id = ?;`, board.Name, board.Position, board.Id)

	return err
}

func NewBoard() (Board, error) {
	var board Board
	db, err := connect()

	if err != nil {
		return board, err
	}

	rows, err := db.Query(`
INSERT INTO board (
    name,
    position
) VALUES ( 'New Board', (SELECT MAX(position) + 1 FROM board) ) RETURNING *;`)

	if err != nil {
		return board, err
	}

	for rows.Next() {
		if err := rows.Scan(&board.Id, &board.Name, &board.Position); err != nil {
			return board, err
		}
	}

	_, err = db.Exec(`
INSERT INTO list (
	board_id,
    name,
    position
) VALUES ( ?, 'New List', 0 ) RETURNING list_id;`, board.Id)

	if err != nil {
		return board, err
	}

	var listId int
	for rows.Next() {
		if err := rows.Scan(&listId); err != nil {
			return board, err
		}
	}

	_, err = db.Exec(`
INSERT INTO item (
	list_id,
    text,
    position
) VALUES ( ?, 'Something Here', 0 );`, listId)

	if err != nil {
		return board, err
	}

	return board, err
}

func DeleteBoard(board Board) error {
	db, err := connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
DELETE 
FROM 
    item
WHERE
	list_id IN (SELECT list_id FROM list WHERE board_id = :board_id);

DELETE 
FROM 
    list 
WHERE
	board_id = :board_id;

DELETE 
FROM 
    board 
WHERE
	board_id = :board_id;
	`, sql.Named("board_id", board.Id))

	return err
}

func DeleteList(listId int) error {
	db, err := connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(`
DELETE 
FROM 
    item
WHERE
	list_id = :list_id

DELETE 
FROM 
    list 
WHERE
	list_id = :list_id;

	`, sql.Named("list_id", listId))

	return err
}

func GetLists(boardId int) []List {
	lists := make([]List, 0)
	db, err := connect()

	if err != nil {
		util.Error("Error connecting to database", err)
	}

	rows, err := db.Query("SELECT list_id, board_id, name, position, (SELECT COUNT('x') FROM item WHERE list_id = list_id) item_count FROM list WHERE board_id = ? ORDER BY position", boardId)

	if err != nil {
		util.Error("Error getting lists", err)
	}

	for rows.Next() {
		var list List
		itemCount := 0
		if err := rows.Scan(&list.Id, &list.BoardId, &list.Title, &list.Position, &itemCount); err != nil {
			util.Error("Error scanning lists", err)
		}

		list.Items = make([]Item, itemCount)
		loadItems(db, &list.Items, list.Id)

		lists = append(lists, list)
	}

	return lists
}

func loadItems(db *sql.DB, items *[]Item, listId int) {

	rows, err := db.Query("SELECT item_id, list_id, text, position FROM item WHERE list_id = ? ORDER BY position", listId)

	if err != nil {
		util.Error("Error loading items", err)
	}

	i := 0
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Id, &item.ListId, &item.Text, &item.Position); err != nil {
			util.Error("Error scanning items", err)
		}

		(*items)[i] = item
	}
}
