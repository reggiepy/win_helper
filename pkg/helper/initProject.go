package helper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"win_helper/pkg/util"
)

var Dirs = []string{"CompanyProject", "PersonProject", "PersonProject"}
var ChildDirs = []string{"PythonProject", "FrontProject", "GolangProject"}

type Project struct {
	BaseDir  string
	DirPaths []string
}

// 定义一个函数签名
type ProjectOption func(*Project)

func BaseDir(baseDir string) ProjectOption {
	return func(p *Project) {
		var err error
		if baseDir == "" {
			baseDir, err = util.GetCurrentDirectory()
			if err != nil {
				panic((fmt.Errorf("get current dir error: %v\n", err)).(any))
			}
		}
		if !path.IsAbs(baseDir) {
			baseDir, err = filepath.Abs(baseDir)
			if err != nil {
				panic((fmt.Errorf("get file abs path error: %v\n", err)).(any))
			}
		}
		p.BaseDir = baseDir
	}
}

func NewProject(options ...ProjectOption) *Project {
	var dirPaths []string
	for i, d := range Dirs {
		dPath := fmt.Sprintf("%02d %s", i+1, d)
		filepath.Join()
		dirPaths = append(dirPaths, dPath)
		for ii, dd := range ChildDirs {
			ddPath := path.Join(dPath, fmt.Sprintf("%02d %s", ii+1, dd))
			dirPaths = append(dirPaths, ddPath)
		}
	}
	project := &Project{
		DirPaths: dirPaths,
	}
	for _, o := range options {
		o(project)
	}
	if project.BaseDir == "" {
		panic((fmt.Errorf("project.BaseDir cannot be empty")).(any))
	}
	project.updateDirPaths()
	return project
}

func (p *Project) Create() {
	for _, dir := range p.DirPaths {
		if util.FileExist(dir) {
			continue
		}
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating: %s\n", err)
		}
	}
}

func (p *Project) updateDirPaths() {
	for i, dir := range p.DirPaths {
		p.DirPaths[i] = filepath.Join(p.BaseDir, dir)
	}
}
