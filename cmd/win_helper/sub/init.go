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

var (
	//  --------------- <<< newInitReadmeCmd --------------------
	initProjectName   string
	initPythonVersion string
	initDjangoVersion string

	//go:embed ..\..\..\templates\Project\README.md.tpl
	initReadmeTemplate []byte
	//  --------------- newInitReadmeCmd >>> --------------------

	//  --------------- <<< newInitProjectCmd --------------------
	initProjectDir    string
	isGenLanguageBool bool
	//  --------------- newInitProjectCmd >>> --------------------
)

func newInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init something",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.UsageString())
			return nil
		},
	}

	initCmd.AddCommand(newInitReadmeCmd())
	initCmd.AddCommand(newInitProjectCmd())
	initCmd.AddCommand(newInitLanguageCmd())

	return initCmd
}

func newInitProjectCmd() *cobra.Command {
	initProjectCmd := &cobra.Command{
		Use:   "project",
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

func newInitLanguageCmd() *cobra.Command {
	initLanguageCmd := &cobra.Command{
		Use:   "language",
		Short: "init language directory",
		Long:  `init language directory`,
		Run: func(cmd *cobra.Command, args []string) {
			languagePaths := project.GenLanguagePaths(baseDir)
			project.CreateProjectDirs(languagePaths)
		},
	}
	return initLanguageCmd
}

func newInitReadmeCmd() *cobra.Command {
	initReadmeCmd := &cobra.Command{
		Use:   "readme",
		Short: "init readme",
		Args:  validateInitReadmeCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			tplExample := pongo2.Must(pongo2.FromBytes(initReadmeTemplate))
			out, err := tplExample.ExecuteBytes(
				pongo2.Context{
					"djangoVersion": initDjangoVersion,
					"pythonVersion": initPythonVersion,
					"projectName":   initProjectName,
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

	initReadmeCmd.Flags().StringVar(&initProjectName, "project-name", "", "project name (required)")
	initReadmeCmd.Flags().StringVar(&initPythonVersion, "python-version", "", "python version (required)")
	initReadmeCmd.Flags().StringVar(&initDjangoVersion, "django-version", "", "django version (required)")
	return initReadmeCmd
}

func validateInitReadmeCmd(cmd *cobra.Command, args []string) error {
	if initProjectName == "" {
		return fmt.Errorf("project-name is required")
	}
	if initDjangoVersion != "" {
		if initPythonVersion == "" {
			return fmt.Errorf("python-version is required")
		}
	}

	if initPythonVersion == "" {
		return fmt.Errorf("python-version is required")
	}

	return nil
}
