package sub

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"win_helper/pkg/util/version"
)

var (
	cfgFile     string
	showVersion bool
	verbose     bool

	// 命令执行路径
	baseDir string
)

func init() {
	cobra.OnInitialize(initConfig)
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
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
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.win_helper.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show verbose output")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	_ = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	rootCmd.AddCommand(newMakeLinkCmd())
	rootCmd.AddCommand(newObrCmd())
	rootCmd.AddCommand(newServerCmd())
	rootCmd.AddCommand(newInitCmd())
	return rootCmd
}

// Execute executes the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func showMessage(msg interface{}) {
	if verbose {
		fmt.Println(msg)
	}
}

func initConfig() {
	baseDir, _ = os.Getwd()
	showMessage(fmt.Sprintf("baseDir : %s", baseDir))
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".win_helper")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
