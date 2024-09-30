package templates

import (
	_ "embed"
)

var (
	//go:embed project/README.md.tpl
	ReadmeTemplate []byte
)
