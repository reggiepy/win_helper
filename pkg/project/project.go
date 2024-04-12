package project

import (
	"fmt"
	"path/filepath"
)

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
	p := &Project{
		DirPaths: []string{},
	}
	for _, o := range options {
		o(p)
	}
	if p.BaseDir == "" {
		panic((fmt.Errorf("project.BaseDir cannot be empty")).(any))
	}
	for i, d := range Dirs {
		dName := fmt.Sprintf("%02d %s", i+1, d)
		dPath := filepath.Join(p.BaseDir, dName)
		p.DirPaths = append(p.DirPaths, dPath)
		if p.IsGenLanguage {
			childPaths := GenLanguagePaths(dPath)
			p.DirPaths = append(p.DirPaths, childPaths...)
		}
	}
	return p
}

func (p *Project) CreateProjectDirs() {
	CreateProjectDirs(p.DirPaths)
}
