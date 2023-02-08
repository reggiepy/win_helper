package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/helper"
	"win_helper/pkg/util"
)

var (
	baseDir string
)

func init() {
	initProjectCmd.Flags().StringVarP(&baseDir, "dir", "d", "", "base directory")
	_ = initProjectCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(initProjectCmd)
}

var initProjectCmd = &cobra.Command{
	Use:   "initProject",
	Short: "init project directory",
	Long:  `init project directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if baseDir == "" {
			baseDir, err = util.GetCurrentDirectory()
			if err != nil {
				fmt.Println(fmt.Sprintf("base directory not set: %v", err))
				return
			}
		}

		project := helper.NewProject(
			helper.BaseDir(baseDir),
		)
		project.Create()
	},
}
