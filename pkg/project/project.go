package project

import (
	"fmt"
	"path/filepath"
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
	"Arduino",
}

type Project struct {
	BaseDir  string
	DirPaths []string
	// 配置
	IsGenLanguage bool
}

// 定义一个函数签名
type ProjectOption func(*Project)

func WithBaseDir(baseDir string) ProjectOption {
	return func(p *Project) {
		p.BaseDir = baseDir
	}
}
func WithIsGenLanguageDir(isGenLanguage bool) ProjectOption {
	return func(p *Project) {
		p.IsGenLanguage = isGenLanguage
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

func (p *Project) GenerateDirs() {
	GenerateDirs(p.DirPaths)
}

func (p *Project) updateDirPaths() {
	for i, dir := range p.DirPaths {
		p.DirPaths[i] = filepath.Join(p.BaseDir, dir)
	}
}
