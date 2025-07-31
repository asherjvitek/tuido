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
);
