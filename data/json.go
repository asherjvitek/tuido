package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// This is the first pass at something comepletely simple and naive but it should get persistence and loading.
const saveDir = ".tuido"
const fileName = "boards.json"

func getPath() string {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Unable to loac the current user, err: %s", err.Error())
	}

	return filepath.Join(user.HomeDir, saveDir, fileName)
}

func SaveData(boards any) {
	jsonData, err := json.Marshal(boards)

	if err != nil {
		panic("Something broke with the jsonData")
	}

	path := getPath()
	dir := filepath.Dir(path)

	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("Unable to create dir %s, err %s", dir, err))
		}

		err = os.Mkdir(dir, os.ModePerm)

		if err != nil {
			panic(fmt.Sprintf("Unable to loac the current user, err: %s", err.Error()))
		}
	}

	err = os.WriteFile(path, jsonData, 0644)

	if err != nil {
		panic(fmt.Sprintf("Failed to Write file %s, err %s", path, err))
	}
}

func LoadData(boards *any) error {
	path := getPath()
	data, err := os.ReadFile(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("failed to read file %s, err: %w", path, err)
	}

	err = json.Unmarshal(data, &boards)

	if err != nil {
		return fmt.Errorf("Failed to Deserialize data from %s", path)
	}

	//We have to do this to make sure that everything is setup right
	// for i := range boards.Boards {
	// 	boards.Boards[i].Setup()
	// }

	return nil
}
