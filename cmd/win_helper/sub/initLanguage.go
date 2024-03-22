package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/helper"
	"win_helper/pkg/util/fileUtils"
)

func init() {
	rootCmd.AddCommand(initLanguageCmd)
}

var initLanguageCmd = &cobra.Command{
	Use:   "initLanguage",
	Short: "init language directory",
	Long:  `init language directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		baseDir, err := fileUtils.GetCurrentDirectory()
		if err != nil {
			fmt.Println(fmt.Sprintf("base directory not set: %v", err))
			return
		}
		languagePaths := helper.GenLanguagePaths(baseDir)
		helper.CreateDirs(languagePaths)
	},
}
