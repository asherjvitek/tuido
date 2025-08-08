package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

type HttpProvider struct{
	Url string
}

func (dp HttpProvider) getUrl(path string) string {
	return fmt.Sprintf("%s%s", dp.Url, path)
}

func (dp HttpProvider) Boards() ([]Board, error) {
	resp, err := http.Get(dp.getUrl("/boards"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boards: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch boards, status code: %d", resp.StatusCode)
	}

	var boards []Board
	json.NewDecoder(resp.Body).Decode(&boards)

	return boards, nil
}

func (dp HttpProvider) Lists(boardId int) ([]List, error) {
	resp, err := http.Get(dp.getUrl("/Lists/boardId="+strconv.Itoa(boardId)))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch lists: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch lists, status code: %d", resp.StatusCode)
	}

	var lists []List
	json.NewDecoder(resp.Body).Decode(&lists)

	return lists, nil
}

func (dp HttpProvider) InsertBoard(board *Board) error {
	content, err := json.Marshal(board)

	if err != nil {
		return fmt.Errorf("failed to marshal board: %w", err)
	}

	resp, err := http.Post(dp.getUrl("boards"), "application/json", bytes.NewBuffer(content))

	if err != nil {
		return fmt.Errorf("failed to fetch lists: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch lists, status code: %d", resp.StatusCode)
	}

	var lists []List
	json.NewDecoder(resp.Body).Decode(&lists)

	return nil
}

func (dp HttpProvider) UpdateBoard(board Board) error {
	panic("unimplemented")
}

func (dp HttpProvider) DeleteBoard(boardId int) error {
	panic("unimplemented")
}

func (dp HttpProvider) InsertList(list *List) error {
	panic("unimplemented")
}

func (dp HttpProvider) DeleteList(list List) error {
	panic("unimplemented")
}

func (dp HttpProvider) UpdateList(list List) error {
	panic("unimplemented")
}

func (dp HttpProvider) InsertItem(item *Item) error {
	panic("unimplemented")
}

func (dp HttpProvider) DeleteItem(item Item) error {
	panic("unimplemented")
}

func (dp HttpProvider) UpdateItem(item Item) error {
	panic("unimplemented")
}
