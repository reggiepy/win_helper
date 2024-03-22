package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/obr/app"
	"win_helper/pkg/util/fileUtils"
)

var (
	// 控制命令是否执行
	pushTag           bool
	replaceIssVersion bool

	gitMessage string
	gitVersion string
	issPath    string
)

func init() {
	rootCmd.AddCommand(obrCmd)
	obrCmd.AddCommand(obrVersionCmd)
	obrVersionCmd.Flags().StringVarP(&gitVersion, "version", "v", "+", "git version message")

	obrVersionCmd.Flags().BoolVarP(&pushTag, "push-tag", "", false, "push tag")
	obrVersionCmd.Flags().StringVarP(&gitMessage, "message", "m", "new version", "git version message")

	obrVersionCmd.Flags().BoolVarP(&replaceIssVersion, "replace-iss-version", "", false, "replace iss version")
	obrVersionCmd.Flags().StringVarP(&issPath, "iss-path", "", "C:\\dist\\chemical_server.iss", "iss_path")
}

var obrCmd = &cobra.Command{
	Use:   "obr",
	Short: "obr tools",
	Long:  `obr tools.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.UsageString())
		return nil
	},
}

var obrVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "handle app version",
	Long:  `handle app version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		versionDir, err := fileUtils.GetCurrentDirectory()
		if err != nil {
			return err
		}
		v := app.NewVersion(
			app.WithMessage(gitMessage),
			app.WithVersion(gitVersion),
			app.WithVersionDir(versionDir),
			app.WithIssPath(issPath),
		)
		currentVersion, err := v.CurrentVersion()
		if err != nil {
			return err
		}
		nowVersion, err := v.GetVersion()
		if err != nil {
			return err
		}
		fmt.Printf("%s ---> %s\n", currentVersion, nowVersion)
		if pushTag {
			fmt.Printf("git message: %s\n", gitMessage)
			err = v.AddTag()
			if err != nil {
				return err
			}
			err = v.PushTags()
			if err != nil {
				return err
			}
		}
		// 替换 iss version
		if replaceIssVersion {
			err = v.ReplaceIssVersion()
			if err != nil {
				return err
			}
		}
		err = v.SaveVersion()
		if err != nil {
			return err
		}
		return nil
	},
}
