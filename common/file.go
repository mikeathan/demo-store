package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func DirExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func CreateDirIfNotExists(path string) {
	if !DirExists(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to create log directory %s Error: %v", path, err))
		}
	}
}

func ToJson(users any) (string, error) {
	data, err := json.Marshal(users)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
