package data

import (
	"fmt"
	"os/user"
	"path/filepath"
)

const saveDir = ".tuido"

func getSavePath(fileName string) (string, error) {
	user, err := user.Current()

	if err != nil {
		return "", fmt.Errorf("Error getting the current user %w", err)
	}

	return filepath.Join(user.HomeDir, saveDir, fileName), nil
}
