package iss

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func GetCurrentVersion(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 使用正则表达式匹配版本号
	versionPattern := `#define\s+MyAppVersion\s+"([\d\.]+)"`

	// 编译正则表达式
	re := regexp.MustCompile(versionPattern)

	// 创建一个 bufio Scanner 来逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 匹配正则表达式
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			version := matches[1]
			fmt.Println("Found version:", version)
			return version, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return "", fmt.Errorf("Error reading file:", err)
	} else {
		return "", fmt.Errorf("Version not found.")
	}
}

func SaveVersion(version string, filePath string) error {
	// 打开原始文件进行读取
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer inputFile.Close()

	// 创建一个临时文件，用于保存修改后的内容
	tempFile, err := os.CreateTemp("", "temp_*.iss")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tempFile.Close()

	// 使用 bufio.Scanner 逐行读取并处理文件
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(tempFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		// 查找并替换版本号
		if strings.Contains(line, "#define MyAppVersion") {
			line = fmt.Sprintf("#define MyAppVersion \"%s\"", version)
		}

		// 将处理后的行写入临时文件
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to temporary file: %v", err)
		}
	}

	// 检查读取文件时的错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// 删除原始文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error removing original file: %v", err)
	}

	// 重命名临时文件为原始文件名
	if err := os.Rename(tempFile.Name(), filePath); err != nil {
		return fmt.Errorf("error renaming temporary file: %v", err)
	}

	return nil
}
