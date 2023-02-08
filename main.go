package main

import (
	"fmt"
	"path"
	"path/filepath"
	"win_helper/pkg/helper"
	"win_helper/pkg/util"
)

func main() {
	fmt.Println(fmt.Sprintf("%02d", 100))
	fmt.Println(path.Join("a", "b"))
	project := helper.NewProject(
		helper.BaseDir("./test"),
	)
	for _, p := range project.DirPaths {
		fmt.Println(p)
	}
	project.Create()
	var file string
	file, _ = util.GetExeDirectory()
	fmt.Println(file)
	file, _ = util.GetCurrentDirectory()
	fmt.Println(file)
	fmt.Println(filepath.Abs(filepath.Join(file, "../../test")))
}
