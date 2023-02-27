package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/helper"
	"win_helper/pkg/util"
)

var (
	baseDir           string
	isGenLanguageStr  string
	isGenLanguageBool bool
)

func init() {
	initProjectCmd.Flags().StringVarP(&baseDir, "dir", "d", "", "base directory")
	initProjectCmd.Flags().StringVarP(&isGenLanguageStr, "language", "l", "false", "gen language directory")
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

		switch isGenLanguageStr {
		case "true":
			isGenLanguageBool = true
		case "false":
			isGenLanguageBool = false
		case "True":
			isGenLanguageBool = true
		case "False":
			isGenLanguageBool = false
		case "1":
			isGenLanguageBool = true
		case "0":
			isGenLanguageBool = false
		}
		option := helper.IsGenLanguageDir(isGenLanguageBool)
		project := helper.NewProject(
			helper.BaseDir(baseDir),
			option,
		)
		project.Create()
	},
}
