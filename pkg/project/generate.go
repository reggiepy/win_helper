package project

import (
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"os"
	"path"
)

func GenLanguagePaths(parent string) []string {
	var dirPaths []string
	for index, childDir := range ChildDirs {
		childPath := path.Join(parent, fmt.Sprintf("%02d %sProject", index+1, childDir))
		dirPaths = append(dirPaths, childPath)
	}
	return dirPaths
}

func CreateProjectDirs(dirs []string) {
	for _, dir := range dirs {
		if fsutil.FileExist(dir) {
			continue
		}
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			fmt.Printf("Error creating: %s\n", err)
		}
	}
}
