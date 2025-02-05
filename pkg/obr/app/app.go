package app

import (
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"log"
	"os"
	"path/filepath"
)

func GetCurrentVersion(dir string) (string, error) {
	var err error
	fileName := filepath.Join(dir, "VERSION")
	if !fsutil.FileExist(fileName) {
		return "", fmt.Errorf("%v does not exist", fileName)
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	return string(data), nil
}
