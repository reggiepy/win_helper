package main

import (
	"fmt"
	"path"
	"path/filepath"
	"win_helper/helper"
	"win_helper/pkg/obr/app"
	"win_helper/pkg/util/fileUtils"
)

func T1() {
	fmt.Println(fmt.Sprintf("%02d", 100))
	fmt.Println(path.Join("a", "b"))
	project := helper.NewProject(
		helper.BaseDir("./test"),
	)
	for _, p := range project.DirPaths {
		fmt.Println(p)
	}
	fmt.Println("**********************************************************************************")
	project.Create()
	var file string
	file, _ = fileUtils.GetExeDirectory()
	fmt.Println(file)
	file, _ = fileUtils.GetCurrentDirectory()
	fmt.Println(file)
	fmt.Println(filepath.Abs(filepath.Join(file, "../../test")))
	fmt.Println(helper.GenLanguagePaths("./"))
}

func T2() {
	versionDir, err := fileUtils.GetCurrentDirectory()
	if err != nil {
		return
	}
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

func main() {
	T3()
}
