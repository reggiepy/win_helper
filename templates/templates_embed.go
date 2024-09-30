package templates

import (
	_ "embed"
)

var (
	//go:embed project/README.md.tpl
	ReadmeTemplate []byte

	//go:embed winsw/WinSW-x64.exe
	WinSW []byte
)
