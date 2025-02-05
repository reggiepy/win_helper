package templates

import (
	"embed"
	"fmt"
	"io"
)

//go:embed templates/*
var Templates embed.FS

func GetTemplateByName(name string) ([]byte, error) {
	file, err := Templates.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open template file %s: %w", name, err)
	}

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", name, err)
	}
	return data, nil
}
