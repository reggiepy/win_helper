package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"win_helper/pkg/util/versionUtils"
)

// NewVersionManager 创建一个新的 VersionManager 实例，并根据提供的选项进行配置。
func NewVersionManager(opts ...Options) *VersionManager {
	v := &VersionManager{}
	return v.WithOptions(opts...)
}

// VersionManager 结构体封装了版本管理的功能，包括获取当前版本、设置版本、保存版本及其他功能。
type VersionManager struct {
	version string // 当前版本
}

// clone 返回 VersionManager 的一个副本，确保对原始实例的修改不会影响到副本。
func (v *VersionManager) clone() *VersionManager {
	clone := *v
	return &clone
}

// WithOptions 根据提供的选项配置 VersionManager。
func (v *VersionManager) WithOptions(opts ...Options) *VersionManager {
	c := v.clone() // 克隆当前实例
	for _, o := range opts {
		o.apply(c) // 应用每个选项到克隆实例
	}
	return c // 返回配置后的新实例
}

// adjustVersion 根据输入的操作符调整版本号
func adjustVersion(proto, major, minor int, op string) (int, int, int, error) {
	switch op {
	case "+++": // 增加补丁版本号
		proto++
	case "---": // 减少补丁版本号
		proto--
	case "++": // 增加主版本号
		major++
	case "--": // 减少主版本号
		major--
	case "+": // 增加次版本号
		minor++
	case "-": // 减少次版本号
		minor--
	case "=", "==", "===": // 减少次版本号

	default:
		return proto, major, minor, fmt.Errorf("invalid version operation: %s", op)
	}

	// 边界检查（可以根据需求进行更多限制，比如不能低于 0）
	if proto < 0 || major < 0 || minor < 0 {
		return proto, major, minor, fmt.Errorf("version number cannot be negative")
	}

	return proto, major, minor, nil
}

// SetVersion 设置下一个版本，根据输入的字符串进行版本号的增减。
func (v *VersionManager) SetVersion(versionStr string) error {
	if versionStr == "" {
		return fmt.Errorf("version string cannot be empty")
	}

	// 当当前版本为空时，直接设置为输入的版本
	if v.version == "" {
		v.version = versionStr
		return nil
	}

	// 提取当前版本的 proto, major 和 minor
	versionProto := versionUtils.Proto(v.version)
	versionMajor := versionUtils.Major(v.version)
	versionMinor := versionUtils.Minor(v.version)

	// 调整版本号
	newProto, newMajor, newMinor, err := adjustVersion(int(versionProto), int(versionMajor), int(versionMinor), versionStr)
	if err != nil {
		return err
	}

	// 设置新的版本字符串
	v.version = fmt.Sprintf("%d.%d.%d", newProto, newMajor, newMinor)
	return nil
}

// GetVersion 获取当前设置的版本。
func (v *VersionManager) GetVersion() string {
	return v.version
}

// Save 将当前版本保存到指定的版本文件中。
func (v *VersionManager) Save(versionFile string, force bool) error {
	fileName := filepath.Join(versionFile, "VERSION") // 组合文件路径
	return v.SaveAs(fileName, force)                  // 调用 SaveAs 方法保存
}

// SaveAs 将当前版本保存到指定文件，支持强制覆盖。
func (v *VersionManager) SaveAs(fileName string, force bool) error {
	if v.version == "" {
		return fmt.Errorf("version cannot be empty") // 版本为空的错误
	}

	// 创建或打开文件，以只写模式打开，文件权限为 0644
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	if !force {
		flags |= os.O_EXCL // 如果不强制，则添加文件存在时的排他性
	}

	// 打开文件
	file, err := os.OpenFile(fileName, flags, 0o644)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err) // 打开文件失败的错误
	}
	defer file.Close() // 确保文件在函数退出时被关闭

	// 写入当前版本到文件
	if _, err := file.WriteString(v.GetVersion()); err != nil {
		return fmt.Errorf("写入文件失败: %v", err) // 写入失败的错误
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("同步文件失败: %v", err) // 同步失败的错误
	}
	return nil
}
