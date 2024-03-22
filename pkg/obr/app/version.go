package app

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"win_helper/pkg/util/fileUtils"
	versionUtil "win_helper/pkg/util/version"
)

// Options 结构体定义了一些配置项，比如版本号、消息和版本目录。
type Options struct {
	Version    string `json:"version"`
	Message    string `json:"message"`
	VersionDir string `json:"version_dir"`
	IssPath    string `json:"iss_path"`
}

// DefaultOption 返回默认选项。
func DefaultOption() *Options {
	return &Options{
		Version: "+",
		Message: "New version",
		IssPath: "C:\\dist\\chemical_server.iss",
	}
}

// Option 接口定义了配置选项的通用行为。
type Option interface {
	apply(opt *Options) error
}

// OptionFunc 用于封装配置选项的函数。
type OptionFunc func(opt *Options) error

// apply 方法应用配置选项到 Options 结构体。
func (f OptionFunc) apply(opt *Options) error {
	return f(opt)
}

// WithVersion 设置版本号的选项。
func WithVersion(version string) Option {
	return OptionFunc(func(opt *Options) error {
		opt.Version = version
		return nil
	})
}

// WithMessage 设置消息的选项。
func WithMessage(message string) Option {
	return OptionFunc(func(opt *Options) error {
		opt.Message = message
		return nil
	})
}

// WithVersionDir 设置版本目录的选项。
func WithVersionDir(versionDir string) Option {
	return OptionFunc(func(opt *Options) error {
		opt.VersionDir = versionDir
		return nil
	})
}

// WithIssPath
func WithIssPath(issPath string) Option {
	return OptionFunc(func(opt *Options) error {
		opt.IssPath = issPath
		return nil
	})
}

// NewVersion 创建一个新的 Version 实例，并根据提供的选项进行配置。
func NewVersion(opts ...Option) *Version {
	opt := DefaultOption()
	for _, o := range opts {
		err := o.apply(opt)
		if err != nil {
			panic(err)
		}
	}
	return &Version{
		Options: opt,
	}
}

// Version 结构体封装了版本管理的功能，包括获取当前版本、保存版本、添加标签和推送标签。
type Version struct {
	Options *Options `json:"options"`
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
	switch v.Options.Version {
	case "":
		return "", fmt.Errorf("version %s is not supported", v.Options.Version)
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
		return v.Options.Version, nil
	}

	// 根据修改后的版本号构建新的版本字符串
	version := fmt.Sprintf("%d.%d.%d", versionProto, versionMajor, versionMinor)
	return version, nil
}

// CurrentVersion 获取当前版本号。
func (v *Version) CurrentVersion() (string, error) {
	var err error
	fileName := filepath.Join(v.Options.VersionDir, "VERSION")
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
	if !fileUtils.FileExist(v.Options.VersionDir) {
		return fmt.Errorf("%v does not exist", v.Options.VersionDir)
	}
	fileName := filepath.Join(v.Options.VersionDir, "VERSION")
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
	if v.Options.Message != "" {
		cmd = append(cmd, "-m", v.Options.Message)
	}
	stdout, stderr, err := runCommand(cmd...)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	fmt.Printf("stderr: %s\n", stderr)
	fmt.Printf("stdout: %s\n", stdout)
	return nil
}

// PushTags 推送所有标签。
func (v *Version) PushTags() error {
	stdout, stderr, err := runCommand("git", "push", "origin", "--tags")
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	fmt.Printf("stderr: %s\n", stderr)
	fmt.Printf("stdout: %s\n", stdout)
	return nil
}

// replaceVersionInFile 用于替换文件中指定行的版本号。
// 参数：
//   - filePath: 文件路径
//   - newVersion: 新的版本号
//
// 返回值：
//   - error: 如果发生错误，返回非空错误；否则返回nil。
func (v *Version) ReplaceIssVersion() error {
	newVersion, err := v.GetVersion()
	if err != nil {
		return err
	}
	filePath := v.Options.IssPath
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// 创建一个临时文件用于保存修改后的内容
	tempFile, err := os.CreateTemp("", "temp_*.iss")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tempFile.Close()

	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(tempFile)
	defer writer.Flush()

	// 使用带缓冲的读取器逐行读取文件内容并替换指定行的版本号
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#define MyAppVersion") {
			// 如果是要替换的行，则将版本号替换为新的版本号
			line = fmt.Sprintf("#define MyAppVersion \"%s\"", newVersion)
		}
		// 将每一行写入临时文件
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return fmt.Errorf("error writing to temporary file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %v", err)
	}

	// 关闭文件
	file.Close()
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
func runCommand(args ...string) (string, string, error) {
	var cmd *exec.Cmd
	var err error
	fmt.Println(strings.Join(args, " "))
	cmd = exec.Command(args[0], args[1:]...) // 将命令和参数分开传递
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}
	return outbuf.String(), errbuf.String(), nil
}