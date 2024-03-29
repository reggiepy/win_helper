package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/project"
	"win_helper/pkg/util/fileUtils"
)

func newInitLanguageCmd() *cobra.Command {
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
			languagePaths := project.GenLanguagePaths(baseDir)
			project.CreateDirs(languagePaths)
		},
	}
	return initLanguageCmd
}
