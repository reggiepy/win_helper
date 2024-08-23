package sub

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/cobra"

	"win_helper/pkg/project"
	"win_helper/pkg/util/fileUtils"
)

//go:embed ..\..\..\templates\Project\README.md.tpl
var initReadmeTemplate []byte

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initProjectCmd)
	initCmd.AddCommand(initLanguageCmd)
	initCmd.AddCommand(initReadmeCmd)

	initProjectCmd.Flags().StringP("dir", "d", "./", "base directory")
	initProjectCmd.Flags().BoolP("language", "l", false, "gen language directory")

	initReadmeCmd.Flags().String("project-name", "", "project name (required)")
	initReadmeCmd.Flags().String("python-version", "", "python version (required)")
	initReadmeCmd.Flags().String("django-version", "", "django version (required)")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init something",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.UsageString())
		return nil
	},
}

var initProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "init project directory。 生成",
	Long:  `init project directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		isGenLanguageBool, _ := cmd.Flags().GetBool("language")
		p := project.NewProject(
			project.WithBaseDir(dir),
			project.WithIsGenLanguageDir(isGenLanguageBool),
		)
		p.CreateProjectDirs()
		return nil
	},
}

var initLanguageCmd = &cobra.Command{
	Use:   "language",
	Short: "init language directory",
	Long:  `init language directory`,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, _ := os.Getwd()
		languagePaths := project.GenLanguagePaths(baseDir)
		project.CreateProjectDirs(languagePaths)
	},
}

var initReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "init readme",
	Args:  validateInitReadmeCmd,
	RunE: func(cmd *cobra.Command, args []string) error {
		tplExample := pongo2.Must(pongo2.FromBytes(initReadmeTemplate))
		projectName, _ := cmd.Flags().GetString("project-name")
		pythonVersion, _ := cmd.Flags().GetString("python-version")
		djangoVersion, _ := cmd.Flags().GetString("django-version")

		out, err := tplExample.ExecuteBytes(
			pongo2.Context{
				"djangoVersion": djangoVersion,
				"pythonVersion": pythonVersion,
				"projectName":   projectName,
			},
		)
		if err != nil {
			return err
		}
		save := func(filename string, data []byte) error {
			writeFile := func() error {
				err := os.WriteFile(filename, data, 0o644)
				if err != nil {
					return fmt.Errorf("写入服务失败。%v", err)
				}
				return nil
			}

			if fileUtils.FileExist(filename) {
				for {
					fmt.Print("File exists. Do you want to continue? (yes/no): ")
					var input string
					fmt.Scanln(&input)
					switch input {
					case "yes", "1":
						return writeFile()
					default:
						fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
					}
				}
			}
			return writeFile()
		}
		fmt.Println(string(out))

		err = save("README.md", out)
		if err != nil {
			return err
		}
		return nil
	},
}

func validateInitReadmeCmd(cmd *cobra.Command, args []string) error {
	projectName, _ := cmd.Flags().GetString("project-name")
	if projectName == "" {
		return fmt.Errorf("project-name is required")
	}

	return nil
}
