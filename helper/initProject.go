package helper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"win_helper/pkg/util"
)

var Dirs = []string{
	"CompanyProject",
	"PersonProject",
	"PublicProject",
	"Scripts",
	"Picture",
	"Document",
	"Software",
}
var ChildDirs = []string{
	"Python",
	"Front",
	"Golang",
	"c51",
	"Docker",
	"Nginx",
	"C",
	"Inno Setup",
	"windows",
}

type Project struct {
	BaseDir  string
	DirPaths []string
	// 配置
	IsGenLanguage bool
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
func IsGenLanguageDir(isGenLanguage bool) ProjectOption {
	return func(p *Project) {
		p.IsGenLanguage = isGenLanguage
	}
}

func GenLanguagePaths(parent string) []string {
	var dirPaths []string
	for index, childDir := range ChildDirs {
		childPath := path.Join(parent, fmt.Sprintf("%02d %sProject", index+1, childDir))
		dirPaths = append(dirPaths, childPath)
	}
	return dirPaths
}

func CreateDirs(dirs []string) {
	for _, dir := range dirs {
		if util.FileExist(dir) {
			continue
		}
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating: %s\n", err)
		}
	}
}

func NewProject(options ...ProjectOption) *Project {
	project := &Project{
		DirPaths: []string{},
	}
	for _, o := range options {
		o(project)
	}
	if project.BaseDir == "" {
		panic((fmt.Errorf("project.BaseDir cannot be empty")).(any))
	}
	for i, d := range Dirs {
		dPath := fmt.Sprintf("%02d %s", i+1, d)
		project.DirPaths = append(project.DirPaths, dPath)
		if project.IsGenLanguage {
			childPaths := GenLanguagePaths(dPath)
			project.DirPaths = append(project.DirPaths, childPaths...)
		}
	}
	project.updateDirPaths()
	return project
}

func (p *Project) Create() {
	CreateDirs(p.DirPaths)
}

func (p *Project) updateDirPaths() {
	for i, dir := range p.DirPaths {
		p.DirPaths[i] = filepath.Join(p.BaseDir, dir)
	}
}
