package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"win_helper/pkg/util/fileUtils"
)

func GetCurrentVersion(dir string) (string, error) {
	var err error
	fileName := filepath.Join(dir, "VERSION")
	if !fileUtils.FileExist(fileName) {
		return "", fmt.Errorf("%v does not exist", fileName)
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	return string(data), nil
}
