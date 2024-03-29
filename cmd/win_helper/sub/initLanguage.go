package sub

import (
	"github.com/spf13/cobra"
	"win_helper/pkg/project"
)

func newInitLanguageCmd() *cobra.Command {
	var initLanguageCmd = &cobra.Command{
		Use:   "initLanguage",
		Short: "init language directory",
		Long:  `init language directory`,
		Run: func(cmd *cobra.Command, args []string) {
			languagePaths := project.GenLanguagePaths(baseDir)
			project.GenerateDirs(languagePaths)
		},
	}
	return initLanguageCmd
}
