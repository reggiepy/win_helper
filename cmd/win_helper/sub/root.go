package sub

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"win_helper/pkg/util/version"
)

var showVersion bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")

	rootCmd.PersistentFlags().Bool("verbose", false, "show verbose output")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

var rootCmd = &cobra.Command{
	Use:   "win_helper",
	Short: "A generator for windows helper",
	Long:  `win_helper is a CLI generator for windows service script`,
	PreRun: func(c *cobra.Command, args []string) {
		return
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		fmt.Println(cmd.UsageString())
		return nil
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func showMessage(format string, a ...any) {
	if viper.GetBool("verbose") {
		fmt.Printf(format, a...)
	}
}

func initConfig() {
	err := godotenv.Load()
	if err != nil {
		showMessage("Error loading .env file: %v\n", err)
	}
	viper.SetEnvPrefix("WH")
	// baseDir, _ := os.Getwd()
	// viper.Set("BASE_DIR", baseDir)
	viper.AutomaticEnv()
}
