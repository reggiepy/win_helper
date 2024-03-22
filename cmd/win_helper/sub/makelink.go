package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"win_helper/pkg/util/fileUtils"
)

var (
	oldname string
	newname string
)

func init() {
	mklink.Flags().StringVarP(&baseDir, "dir", "d", "", "base directory")
	mklink.Flags().StringVarP(&oldname, "oldname", "o", "", "gen language directory")
	mklink.Flags().StringVarP(&newname, "newname", "n", "", "gen language directory")
	rootCmd.AddCommand(mklink)
}

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
	err := os.Link(oldName, newName)
	if err != nil {
		return fmt.Errorf("创建软链接失败: %v", err)
	}
	return nil
}

var mklink = &cobra.Command{
	Use:   "mklink",
	Short: "windows make link",
	Long:  `windows make link`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if baseDir == "" {
			baseDir, err = fileUtils.GetCurrentDirectory()
			if err != nil {
				fmt.Println(fmt.Sprintf("base directory not set: %v", err))
				return
			}
		}
		if oldname == "" || newname == "" {
			fmt.Println(fmt.Sprintf("newname or oldname can't be empty. oldname: %v newname: %v", oldname, newname))
			return
		}
		err = makeLink2(oldname, newname)
		if err != nil {
			fmt.Println(err)
		}
	},
}
