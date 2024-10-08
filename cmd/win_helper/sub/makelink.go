package sub

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type MakeLinkConfig struct {
	OldName string
	NewName string
}

var makelinkConfig = MakeLinkConfig{}

func init() {
	rootCmd.AddCommand(makeLinkCmd)

	makeLinkCmd.Flags().StringVarP(&makelinkConfig.OldName, "old-name", "o", "", "gen language directory")
	makeLinkCmd.Flags().StringVarP(&makelinkConfig.NewName, "new-name", "n", "", "gen language directory")
}

var makeLinkCmd = &cobra.Command{
	Use:   "mklink",
	Short: "windows make link",
	Long:  `windows make link`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if makelinkConfig.OldName == "" || makelinkConfig.NewName == "" {
			return fmt.Errorf("newname or oldname can't be empty'")
		}
		err = makeLink2(makelinkConfig.OldName, makelinkConfig.NewName)
		if err != nil {
			return err
		}
		return err
	},
}

func makeLink1(oldName string, newName string) error {
	_, err := os.Stat(oldName)
	if err != nil {
		return fmt.Errorf("获取文件夹信息失败: %v", err)
	}
	cmd := exec.Command("cmd", "/c", "mklink", "/D", oldName, newName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建软连接失败: %v", err)
	}
	fmt.Print(string(out))
	return nil
}

func makeLink2(oldName string, newName string) error {
	err := os.Symlink(oldName, newName)
	if err != nil {
		return fmt.Errorf("创建软链接失败: %v", err)
	}
	return nil
}
