package sub

import (
	"fmt"
	"win_helper/pkg/server"

	"github.com/spf13/cobra"
)


type WinServiceConfig struct {
	Force                  bool
	ID                     string
	Executable             string
	Name                   string
	Description            string
	StartMode              string
	Depends                string
	LogPath                string
	Arguments              string
	StartArguments         string
	StopExecutable         string
	StopArguments          string
	Env                    string
	Failure                string
	WorkingDirectory       string
	LogMode                string
	LogPattern             string
	LogAutoRollAtTime      string
	LogSizeThreshold       string
	LogZipOlderThanNumDays string
	LogZipDateFormat       string
}

var serverConfig = WinServiceConfig{}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&serverConfig.ID, "id", "", "Id(default=name)")
	serverCmd.Flags().StringVar(&serverConfig.Name, "name", "", "name")
	serverCmd.Flags().StringVar(&serverConfig.Executable, "executable", "", "executable")
	serverCmd.Flags().StringVar(&serverConfig.Description, "description", "", "description")
	serverCmd.Flags().StringVar(&serverConfig.StartMode, "start-mode", "", "start-mode(Boot|System|Automatic|Manual|Disabled) (default: Automatic)")
	serverCmd.Flags().StringVar(&serverConfig.Depends, "depends", "", "depends")
	serverCmd.Flags().StringVar(&serverConfig.LogPath, "log-path", "logs", "log path")
	serverCmd.Flags().StringVar(&serverConfig.Arguments, "arguments", "", "arguments")
	serverCmd.Flags().StringVar(&serverConfig.StartArguments, "start-arguments", "", "start arguments")
	serverCmd.Flags().StringVar(&serverConfig.StopExecutable, "stop-executable", "", "stop executable")
	serverCmd.Flags().StringVar(&serverConfig.StopArguments, "stop-arguments", "", "stop arguments")
	serverCmd.Flags().StringVar(&serverConfig.Env, "env", "", "environment variables")
	serverCmd.Flags().StringVar(&serverConfig.Failure, "failure", "", "failure")
	serverCmd.Flags().StringVar(&serverConfig.WorkingDirectory, "working-directory", "", "working directory")
	serverCmd.Flags().StringVar(&serverConfig.LogMode, "log-mode", "roll", "log mode")
	serverCmd.Flags().StringVar(&serverConfig.LogPattern, "log-pattern", "", "log pattern")
	serverCmd.Flags().StringVar(&serverConfig.LogAutoRollAtTime, "log-auto-roll-at-time", "", "log auto roll at time")
	serverCmd.Flags().StringVar(&serverConfig.LogSizeThreshold, "log-size-threshold", "", "log size threshold")
	serverCmd.Flags().StringVar(&serverConfig.LogZipOlderThanNumDays, "log-zip-older-than-num-days", "", "log zip older than num days")
	serverCmd.Flags().StringVar(&serverConfig.LogZipDateFormat, "log-zip-date-format", "", "log zip date format")
	serverCmd.Flags().BoolVar(&serverConfig.Force, "force", true, "force write")

	// Boot Start ("Boot")
	// Device driver started by the operating system loader. This value is valid only for driver services.
	// System ("System")
	// Device driver started by the operating system initialization process. This value is valid only for driver services.
	// Auto Start ("Automatic")
	// Service to be started automatically by the service control manager during system startup.
	// Demand Start ("Manual")
	// Service to be started by the service control manager when a process calls the StartService method.
	// Disabled ("Disabled")
	// Service that can no longer be started.
}

var serverCmd = &cobra.Command{
	Use:   "winserver-gen",
	Short: "generate exe file's windows server",
	Long:  `generate exe file's windows server`,
	Args: func(cmd *cobra.Command, args []string) error {
		if serverConfig.Name == "" {
			return fmt.Errorf("missing name")
		}
		if serverConfig.ID == "" {
			serverConfig.ID = serverConfig.Name
		}
		if serverConfig.Executable == "" {
			return fmt.Errorf("missing executable")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := server.NewServer(
			server.WithSId(serverConfig.ID),
			server.WithSName(serverConfig.Name),
			server.WithSExecutable(serverConfig.Executable),
			server.WithSDescription(serverConfig.Description),
			server.WithSStartMode(serverConfig.StartMode),
			server.WithSDepends(serverConfig.Depends),
			server.WithSLogPath(serverConfig.LogPath),
			server.WithSArguments(serverConfig.Arguments),
			server.WithSStartArguments(serverConfig.StartArguments),
			server.WithSStopExecutable(serverConfig.StopExecutable),
			server.WithSStopArguments(serverConfig.StopArguments),
			server.WithSEnv(serverConfig.Env),
			server.WithSFailure(serverConfig.Failure),
			server.WithSWorkingDirectory(serverConfig.WorkingDirectory),
			server.WithSLogMode(serverConfig.LogMode),
			server.WithSLogPattern(serverConfig.LogPattern),
			server.WithSLogAutoRollAtTime(serverConfig.LogAutoRollAtTime),
			server.WithSLogSizeThreshold(serverConfig.LogSizeThreshold),
			server.WithSLogZipOlderThanNumDays(serverConfig.LogZipOlderThanNumDays),
			server.WithSLogZipDateFormat(serverConfig.LogZipDateFormat),
			server.WithSForce(serverConfig.Force),
		)
		if err != nil {
			return err
		}
		err = s.Generate()
		if err != nil {
			return err
		}
		return nil
	},
}
