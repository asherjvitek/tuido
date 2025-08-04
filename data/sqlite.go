package data

import (
	"database/sql"
	"math"
	_ "modernc.org/sqlite"
	"os"
	"os/user"
	"path/filepath"
)

func getDbPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".tuido", "tuido.db"), nil
}

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Init() error {
	path, err := getDbPath()

	if err != nil {
		return err
	}

	initDb := false

	if _, err := os.Stat(path); os.IsNotExist(err) {
		initDb = true
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
	}

	if !initDb {
		return nil
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	initDbSql := `
CREATE TABLE IF NOT EXISTS Board (
    BoardId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT,
    Position REAL
);

CREATE TABLE IF NOT EXISTS List (
    ListId INTEGER PRIMARY KEY AUTOINCREMENT,
    BoardId INTEGER,
    Name TEXT,
    Position REAL,
    FOREIGN KEY (BoardId) REFERENCES Board (BoardId)
);

CREATE TABLE IF NOT EXISTS Item (
    ItemId INTEGER PRIMARY KEY AUTOINCREMENT,
    ListId INTEGER,
    Text TEXT,
    Position REAL,
    FOREIGN KEY (ListId) REFERENCES List (ListId)
);

INSERT INTO
    Board (Name, Position)
VALUES
	('YOUR FIRST BOARD', :startingPosition);

INSERT INTO
    List (BoardId, Name, Position)
VALUES
	(1, 'YOUR FIRST BOARD', :startingPosition);

INSERT INTO
    Item (ListId, Text, Position)
VALUES
	(1, 'LETS GO TODO!', :startingPosition);
	`

	startingPosition := math.MaxFloat64 / 2.0

	_, err = db.Exec(initDbSql, sql.Named("startingPosition", startingPosition))

	if err != nil {
		return err
	}

	return nil
}

func Boards() ([]Board, error) {

	path, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT BoardId, Name, Position FROM Board ORDER BY Position")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]Board, 0)
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.BoardId, &b.Name, &b.Position); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func Lists(boardId int) ([]List, error) {

	path, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT ListId, BoardId, Name, Position FROM List WHERE BoardId = ? ORDER BY Position", boardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := make([]List, 0)
	for rows.Next() {
		var l List
		if err := rows.Scan(&l.ListId, &l.BoardId, &l.Name, &l.Position); err != nil {
			return nil, err
		}

		items, err := Items(l.ListId)

		if err != nil {
			return nil, err
		}

		l.Items = items

		lists = append(lists, l)
	}

	return lists, nil
}

func Items(listId int) ([]Item, error) {
	path, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT ItemId, ListId, Text, Position FROM Item WHERE ListId = ? ORDER BY Position", listId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var l Item
		if err := rows.Scan(&l.ItemId, &l.ListId, &l.Text, &l.Position); err != nil {
			return nil, err
		}
		items = append(items, l)
	}

	return items, nil
}

func InsertBoard(board *Board) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("INSERT INTO Board (Name, Position) VALUES (?, ?) RETURNING BoardId", board.Name, board.Position)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&board.BoardId); err != nil {
			return err
		}
	}

	return nil
}

func UpdateBoard(board Board) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Board SET Name = ?, Position = ? WHERE BoardId = ?", board.Name, board.Position, board.BoardId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBoard(board Board) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM Item WHERE ListId IN (SELECT ListId FROM List WHERE BoardId = :BoardId);
		DELETE FROM List WHERE BoardId = :BoardId;
		DELETE FROM Board WHERE BoardId = :BoardId;`, sql.Named("BoardId", board.BoardId))
	if err != nil {
		return err
	}

	return nil
}

func InsertList(list *List) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("INSERT INTO List (BoardId, Name, Position) VALUES (?, ?, ?) RETURNING ListId", list.BoardId, list.Name, list.Position)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&list.ListId); err != nil {
			return err
		}
	}

	return nil
}

func DeleteList(list List) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Item WHERE ListId = ?", list.ListId)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM List WHERE ListId = ?", list.ListId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateList(list List) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE List SET Name = ?, Position = ? WHERE ListId = ?", list.Name, list.Position, list.ListId)
	if err != nil {
		return err
	}

	return nil
}

func InsertItem(item *Item) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("INSERT INTO Item (ListId, Text, Position) VALUES (?, ?, ?) RETURNING ItemId", item.ListId, item.Text, item.Position)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&item.ItemId); err != nil {
			return err
		}
	}

	return nil
}

func DeleteItem(item Item) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Item WHERE ItemId = ?", item.ItemId)

	if err != nil {
		return err
	}

	return nil
}

func UpdateItem(item Item) error {
	path, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := Open(path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Item SET ListId = ?, Text = ?, Position = ? WHERE ItemId = ?", item.ListId, item.Text, item.Position, item.ItemId)
	if err != nil {
		return err
	}

	return nil
}
