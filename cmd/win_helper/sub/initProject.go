package sub

import (
	"github.com/spf13/cobra"
	"win_helper/pkg/project"
)

var (
	initProjectDir    string
	isGenLanguageBool bool
)

func newInitProjectCmd() *cobra.Command {
	var initProjectCmd = &cobra.Command{
		Use:   "initProject",
		Short: "init project directory。 生成",
		Long:  `init project directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := project.NewProject(
				project.WithBaseDir(initProjectDir),
				project.WithIsGenLanguageDir(isGenLanguageBool),
			)
			p.CreateProjectDirs()
			return nil
		},
	}
	initProjectCmd.Flags().StringVarP(&initProjectDir, "dir", "d", "./", "base directory")
	initProjectCmd.Flags().BoolVarP(&isGenLanguageBool, "language", "l", false, "gen language directory")
	return initProjectCmd
}
