package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/project"
	"win_helper/pkg/util/fileUtils"
)

var (
	baseDir           string
	isGenLanguageBool bool
)

func newInitProjectCmd() *cobra.Command {
	var initProjectCmd = &cobra.Command{
		Use:   "initProject",
		Short: "init project directory",
		Long:  `init project directory`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if baseDir == "" {
				baseDir, err = fileUtils.GetCurrentDirectory()
				if err != nil {
					fmt.Println(fmt.Sprintf("base directory not set: %v", err))
					return
				}
			}
			option := project.IsGenLanguageDir(isGenLanguageBool)
			project := project.NewProject(
				project.BaseDir(baseDir),
				option,
			)
			project.Create()
		},
	}
	initProjectCmd.Flags().StringVarP(&baseDir, "dir", "d", "", "base directory")
	initProjectCmd.Flags().BoolVarP(&isGenLanguageBool, "language", "l", false, "gen language directory")
	return initProjectCmd
}
