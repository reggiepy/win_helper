package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	oldname string
	newname string
)

func makeLink1(oldName string, newName string) error {
	_, err := os.Stat(oldname)
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

func newMakeLinkCmd() *cobra.Command {
	makeLinkCmd := &cobra.Command{
		Use:   "mklink",
		Short: "windows make link",
		Long:  `windows make link`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if oldname == "" || newname == "" {
				return fmt.Errorf("newname or oldname can't be empty'")
			}
			err = makeLink2(oldname, newname)
			if err != nil {
				return err
			}
			return err
		},
	}
	makeLinkCmd.Flags().StringVarP(&oldname, "oldname", "o", "", "gen language directory")
	makeLinkCmd.Flags().StringVarP(&newname, "newname", "n", "", "gen language directory")
	return makeLinkCmd
}
