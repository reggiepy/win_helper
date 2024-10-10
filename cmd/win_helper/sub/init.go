package sub

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/cobra"

	"win_helper/templates"

	"win_helper/pkg/project"
	"win_helper/pkg/util/fileUtils"
)

type InitProjectConfig struct {
	Directory string
	Language  bool
}

type InitReadmeConfig struct {
	ProjectName   string
	PythonVersion string
	DjangoVersion string
	Force         bool
	Verbose       bool
	Shields       []string
	OutFile       string
}

var (
	initReadmeConfig  = InitReadmeConfig{}
	initProjectConfig = InitProjectConfig{}
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initProjectCmd)
	initCmd.AddCommand(initLanguageCmd)
	initCmd.AddCommand(initReadmeCmd)

	initProjectCmd.Flags().StringVarP(&initProjectConfig.Directory, "dir", "d", "./", "base directory")
	initProjectCmd.Flags().BoolVarP(&initProjectConfig.Language, "language", "l", false, "gen language directory")

	initReadmeCmd.Flags().StringArrayVar(&initReadmeConfig.Shields, "shields", []string{}, "name|value|description")
	initReadmeCmd.Flags().StringVarP(&initReadmeConfig.ProjectName, "name", "n", "", "project name (*)")
	initReadmeCmd.Flags().StringVarP(&initReadmeConfig.OutFile, "out", "o", "README.md", "output file (default: README.md)")
	initReadmeCmd.Flags().StringVar(&initReadmeConfig.PythonVersion, "python-version", "", "python version")
	initReadmeCmd.Flags().StringVar(&initReadmeConfig.DjangoVersion, "django-version", "", "django version")
	initReadmeCmd.Flags().BoolVarP(&initReadmeConfig.Force, "force", "f", false, "force write file (default: false)")
	initReadmeCmd.Flags().BoolVar(&initReadmeConfig.Verbose, "verbose", false, "verbose")
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
		// 将相对路径转换为绝对路径
		absPath, err := filepath.Abs(initProjectConfig.Directory)
		if err != nil {
			return err
		}
		p := project.NewProject(
			project.WithBaseDir(absPath),
			project.WithIsGenLanguageDir(initProjectConfig.Language),
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

type Shield struct {
	Name        string
	Value       string
	Description string
}

var initReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "init readme",
	Args:  validateInitReadmeCmd,
	RunE: func(cmd *cobra.Command, args []string) error {
		shields := make([]Shield, 0, 4)
		shields = append(shields, ParseShieldString(initReadmeConfig.Shields)...)
		if initReadmeConfig.DjangoVersion != "" {
			shields = append(shields, Shield{
				Name:        "Django",
				Value:       initReadmeConfig.DjangoVersion,
				Description: "django version",
			})
		}
		if initReadmeConfig.PythonVersion != "" {
			shields = append(shields, Shield{
				Name:        "Python",
				Value:       initReadmeConfig.PythonVersion,
				Description: "python version",
			})
		}
		shields = append(shields, Shield{
			Name:        "build",
			Value:       "pass",
			Description: "build status",
		})
		tplExample := pongo2.Must(pongo2.FromBytes(templates.ReadmeTemplate))

		out, err := tplExample.ExecuteBytes(pongo2.Context{
			"shields":     shields,
			"projectName": initReadmeConfig.ProjectName,
		})
		if err != nil {
			return err
		}
		if initReadmeConfig.Verbose {
			fmt.Println(string(out))
		} else {
			filename := initReadmeConfig.OutFile
			if fileUtils.FileExist(filename) && !initReadmeConfig.Force {
				if !initReadmeConfig.Force {
				outerLoop:
					for {
						fmt.Print("File exists. Do you want to continue? (yes/no)(default:no): ")
						var input string
						_, _ = fmt.Scanln(&input)
						switch input {
						case "yes", "1", "true", "True":
							break outerLoop
						case "no", "0", "false", "False":
							return nil
						default:
							fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
						}
					}
				} else {
					return fmt.Errorf("文件已经存在")
				}
			}
			err = os.WriteFile(filename, out, 0o644)
			if err != nil {
				return fmt.Errorf("写入文件失败。%v", err)
			}
		}
		return nil
	},
}

func ParseShieldString(shieldString []string) []Shield {
	var shields []Shield
	for _, shield := range shieldString {
		if shield == "" {
			continue
		}
		shieldsSplit := strings.Split(shield, "|")
		if len(shieldsSplit) != 3 {
			fmt.Printf("Invalid shield: %s", shield)
			continue
		}
		name, value, description := shieldsSplit[0], shieldsSplit[1], shieldsSplit[2]
		escapedValue := url.PathEscape(value)
		escapedName := url.PathEscape(name)
		shields = append(shields, Shield{Name: escapedName, Value: escapedValue, Description: description})
	}
	return shields
}

func validateInitReadmeCmd(cmd *cobra.Command, args []string) error {
	if initReadmeConfig.ProjectName == "" {
		return fmt.Errorf("project-name is required")
	}
	return nil
}
