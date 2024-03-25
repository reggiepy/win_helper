package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"win_helper/pkg/util/fileUtils"
	versionUtil "win_helper/pkg/util/version"
)

// DefaultOption 返回默认选项。
func DefaultVersion() *Version {
	return &Version{
		Version: "+",
		Message: "New version",
		IssPath: "C:\\dist\\chemical_server.iss",
	}
}

// VersionFunc 用于封装配置选项的函数。
type VersionFunc func(v *Version) error

// WithVersion 设置版本号的选项。
func WithVersion(version string) VersionFunc {
	return func(v *Version) error {
		v.Version = version
		return nil
	}
}

// WithMessage 设置消息的选项。
func WithMessage(message string) VersionFunc {
	return func(v *Version) error {
		v.Message = message
		return nil
	}
}

// WithVersionDir 设置版本目录的选项。
func WithVersionDir(versionDir string) VersionFunc {
	return func(v *Version) error {
		v.VersionDir = versionDir
		return nil
	}
}

// WithIssPath
func WithIssPath(issPath string) VersionFunc {
	return func(v *Version) error {
		v.IssPath = issPath
		return nil
	}
}

// NewVersion 创建一个新的 Version 实例，并根据提供的选项进行配置。
func NewVersion(vfs ...VersionFunc) *Version {
	v := DefaultVersion()
	for _, f := range vfs {
		err := f(v)
		if err != nil {
			panic(err)
		}
	}
	return v
}

// Version 结构体封装了版本管理的功能，包括获取当前版本、保存版本、添加标签和推送标签。
type Version struct {
	Version    string `json:"version"`
	Message    string `json:"message"`
	VersionDir string `json:"version_dir"`
	IssPath    string `json:"iss_path"`
}

// GetVersion 获取下一个版本。
func (v *Version) GetVersion() (string, error) {
	versionStr, err := v.CurrentVersion()
	if err != nil {
		return "", err
	}
	versionStr = strings.TrimSpace(versionStr)
	versionProto := versionUtil.Proto(versionStr)
	versionMajor := versionUtil.Major(versionStr)
	versionMinor := versionUtil.Minor(versionStr)

	// 根据配置的版本号增加或减少版本号的各个部分
	switch v.Version {
	case "":
		return "", fmt.Errorf("version %s is not supported", v.Version)
	case "+++":
		versionProto++
	case "---":
		versionProto--
	case "++":
		versionMajor++
	case "--":
		versionMajor--
	case "+":
		versionMinor++
	case "-":
		versionMinor--
	default:
		// 如果版本号类型未知，则直接返回配置的版本号
		return v.Version, nil
	}

	// 根据修改后的版本号构建新的版本字符串
	version := fmt.Sprintf("%d.%d.%d", versionProto, versionMajor, versionMinor)
	return version, nil
}

// CurrentVersion 获取当前版本号。
func (v *Version) CurrentVersion() (string, error) {
	var err error
	fileName := filepath.Join(v.VersionDir, "VERSION")
	if !fileUtils.FileExist(fileName) {
		return "", fmt.Errorf("%v does not exist", fileName)
	}
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("open: %v error", err)
	}
	defer file.Close()

	buf := make([]byte, 4096)
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("read: %v error", err)
	}
	versionStr := string(buf[:n])
	return versionStr, nil
}

// SaveVersion 保存版本号。
func (v *Version) SaveVersion() error {
	version, err := v.GetVersion()
	if err != nil {
		return err
	}
	if !fileUtils.FileExist(v.VersionDir) {
		return fmt.Errorf("%v does not exist", v.VersionDir)
	}
	fileName := filepath.Join(v.VersionDir, "VERSION")
	// 创建或打开文件，以只写模式打开，文件权限为 0644
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer file.Close()
	// 创建一个带缓冲的写入器，并指定编码为 UTF-8
	writer := bufio.NewWriter(file)

	// 将 UTF-8 编码的字符串写入文件
	_, err = writer.WriteString(version)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// 刷新写入器的缓冲区，确保所有数据被写入文件
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

// AddTag 添加一个新的标签。
func (v *Version) AddTag() error {
	version, err := v.GetVersion()
	tagName := fmt.Sprintf("v%s", version)
	var cmd = []string{"git", "tag", tagName}
	if v.Message != "" {
		cmd = append(cmd, "-m", v.Message)
	}
	_, err = runCommand(cmd...)
	if err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}
	return nil
}

// PushTags 推送所有标签。
func (v *Version) PushTags() error {
	_, err := runCommand("git", "push", "origin", "--tags")
	if err != nil {
		return fmt.Errorf("推送标签失败: %v", err)
	}
	return nil
}

func (v *Version) HandleTag() error {
	var err error
	fmt.Printf("git message: %s\n", v.Message)
	err = v.AddTag()
	if err != nil {
		return err
	}
	err = v.PushTags()
	if err != nil {
		return err
	}
	return nil
}

// ReplaceIssVersion 用于替换文件中指定行的版本号。
// 参数：
//   - filePath: 文件路径
//   - newVersion: 新的版本号
//
// 返回值：
//   - error: 如果发生错误，返回非空错误；否则返回nil。
func (v *Version) ReplaceIssVersion() error {
	// 获取新的版本号
	newVersion, err := v.GetVersion()
	if err != nil {
		return err
	}
	filePath := v.IssPath

	// 读取原始文件内容到内存中
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// 将文件内容转换为字符串进行处理
	fileContent := string(fileBytes)

	// 以换行符分割文件内容，保留换行符
	lines := strings.SplitAfterN(fileContent, "\n", -1)

	// 替换版本号
	for i, line := range lines {
		if strings.Contains(line, "#define MyAppVersion") {
			// 如果是版本号行，则替换版本号
			lines[i] = fmt.Sprintf("#define MyAppVersion \"%s\"\n", newVersion)
			break // 替换完第一个版本号后退出循环
		}
	}

	// 将修改后的内容拼接为字符串
	newFileContent := strings.Join(lines, "")

	// 创建一个临时文件用于保存修改后的内容
	tempFile, err := os.CreateTemp("", "temp_*.iss")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tempFile.Close()

	// 将修改后的内容写入临时文件
	_, err = tempFile.WriteString(newFileContent)
	if err != nil {
		return fmt.Errorf("error writing to temporary file: %v", err)
	}

	// 关闭临时文件
	tempFile.Close()

	// 删除原始文件
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error removing original file: %v", err)
	}

	// 重命名临时文件为原始文件名
	err = os.Rename(tempFile.Name(), filePath)
	if err != nil {
		return fmt.Errorf("error renaming temporary file: %v", err)
	}

	return nil
}

// runCommand 执行系统命令，并返回标准输出和标准错误输出。
func runCommand(args ...string) (string, error) {
	var cmd *exec.Cmd
	var err error
	cmdStr := strings.Join(args, " ")
	fmt.Printf("执行命令: %s\n", cmdStr)
	cmd = exec.Command(args[0], args[1:]...) // 将命令和参数分开传递
	out, err := cmd.CombinedOutput()
	fmt.Printf("执行命令输出: %s\n", string(out))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
