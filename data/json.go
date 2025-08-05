package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"tuido/util"
)

// This is the first pass at something comepletely simple and naive but it should get persistence and loading.
const fileName = "boards.json"

func getPath() (string, error) {
	appDir, err := util.GetAppDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, fileName), nil
}

func SaveData(boards any) error {
	jsonData, err := json.Marshal(boards)

	if err != nil {
		return err
	}

	path, err := getPath()

	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.Mkdir(dir, os.ModePerm)

		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, jsonData, 0644)

	if err != nil {
		return err
	}

	return nil
}

func LoadData(boards *any) error {
	path, err := getPath()

	if err != nil {
		return err
	}
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

	return nil
}
