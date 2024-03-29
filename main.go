package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"win_helper/pkg/obr/app"
	"win_helper/pkg/project"
	"win_helper/pkg/util/fileUtils"
)

func T1() {
	fmt.Println(fmt.Sprintf("%02d", 100))
	fmt.Println(path.Join("a", "b"))
	p := project.NewProject(
		project.WithBaseDir("./test"),
	)
	for _, p := range p.DirPaths {
		fmt.Println(p)
	}
	fmt.Println("**********************************************************************************")
	p.GenerateDirs()
	var file string
	file, _ = fileUtils.GetExeDirectory()
	fmt.Println(file)
	fmt.Println(filepath.Abs(filepath.Join(file, "../../test")))
	fmt.Println(project.GenLanguagePaths("./"))
}

func T2() {
	var err error
	versionDir, _ := os.Getwd()
	v := app.NewVersion(
		app.WithVersion("+"),
		app.WithVersionDir(versionDir),
	)
	fmt.Println(v.GetVersion())
	err = v.AddTag()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = v.PushTags()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = v.SaveVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func T3() {
	a := "1.0.0"
	switch a {
	case "+++":
		fmt.Printf("+++: %s\n", a)
	case "++":
		fmt.Printf("++: %s\n", a)
	case "+":
		fmt.Printf("+: %s\n", a)
	default:
		fmt.Printf("default: %s\n", a)
	}
}

func T4() {
	originalPath := "test"
	linkPath := "test2"
	err := os.Symlink(originalPath, linkPath)
	if err != nil {
		fmt.Printf("Error creating symlink: %v", err)
		return
	}
	fmt.Println("Symlink created successfully.")
}

func main() {
	T4()
}
