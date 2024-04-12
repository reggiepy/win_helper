package sub

import (
	"fmt"

	"github.com/spf13/cobra"
	"win_helper/pkg/obr/app"
)

var (
	// 控制命令是否执行
	pushTag           bool
	replaceIssVersion bool

	gitMessage string
	gitVersion string
	issPath    string
)

func newObrCmd() *cobra.Command {
	obrCmd := &cobra.Command{
		Use:   "obr",
		Short: "obr tools",
		Long:  `obr tools.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.UsageString())
			return nil
		},
	}
	obrCmd.AddCommand(newObrVersionCmd())
	return obrCmd
}

func newObrVersionCmd() *cobra.Command {
	obrVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "handle app version",
		Long:  `handle app version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			v := app.NewVersion(
				app.WithMessage(gitMessage),
				app.WithVersion(gitVersion),
				app.WithVersionDir(baseDir),
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
				err = v.HandleTag()
				if err != nil {
					fmt.Printf("处理tag异常: %v\n", err)
				}
			}
			// 替换 iss version
			if replaceIssVersion {
				err = v.ReplaceIssVersion()
				if err != nil {
					fmt.Printf("替换iss版本异常: %v\n", err)
				}
			}
			err = v.SaveVersion()
			if err != nil {
				fmt.Printf("保存版本文件异常: %v\n", err)
				return err
			}
			return nil
		},
	}
	obrVersionCmd.Flags().StringVarP(&gitVersion, "version", "v", "+", "git version message")

	obrVersionCmd.Flags().BoolVarP(&pushTag, "push-tag", "", false, "push tag")
	obrVersionCmd.Flags().StringVarP(&gitMessage, "message", "m", "new version", "git version message")

	obrVersionCmd.Flags().BoolVarP(&replaceIssVersion, "replace-iss-version", "", false, "replace iss version")
	obrVersionCmd.Flags().StringVarP(&issPath, "iss-path", "", "C:\\dist\\chemical_server.iss", "iss_path")
	return obrVersionCmd
}
