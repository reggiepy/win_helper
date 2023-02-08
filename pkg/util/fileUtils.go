package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
func GetExeDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", fmt.Errorf("get current directory error: %v", err)
	}
	return dir, nil
	//return strings.Replace(dir, "\\", "/", -1), nil
}

func GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get current directory error: %v", err)
	}
	return dir, nil
}
