package sub

import (
	"fmt"

	"github.com/spf13/cobra"
	"win_helper/pkg/server/win"
)

var (
	sForce bool

	sId                     string
	sExecutable             string
	sName                   string
	sDescription            string
	sStartMode              string
	sDepends                string
	sLogPath                string
	sArguments              string
	sStartArguments         string
	sStopExecutable         string
	sStopArguments          string
	sEnv                    string
	sFailure                string
	sWorkingDirectory       string
	sLogMode                string
	sLogPattern             string
	sLogAutoRollAtTime      string
	sLogSizeThreshold       string
	sLogZipOlderThanNumDays string
	sLogZipDateFormat       string
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&sId, "id", "", "Id(default=name)")
	serverCmd.Flags().StringVar(&sName, "name", "", "name")
	serverCmd.Flags().StringVar(&sExecutable, "executable", "", "executable")
	serverCmd.Flags().StringVar(&sDescription, "description", "", "description")
	serverCmd.Flags().StringVar(&sStartMode, "start-mode", "", "start-mode(Boot|System|Automatic|Manual|Disabled) (default: Automatic)")
	serverCmd.Flags().StringVar(&sDepends, "depends", "", "depends")
	serverCmd.Flags().StringVar(&sLogPath, "log-path", "logs", "log path")
	serverCmd.Flags().StringVar(&sArguments, "arguments", "", "arguments")
	serverCmd.Flags().StringVar(&sStartArguments, "start-arguments", "", "start arguments")
	serverCmd.Flags().StringVar(&sStopExecutable, "stop-executable", "", "stop executable")
	serverCmd.Flags().StringVar(&sStopArguments, "stop-arguments", "", "stop arguments")
	serverCmd.Flags().StringVar(&sEnv, "env", "", "environment variables")
	serverCmd.Flags().StringVar(&sFailure, "failure", "", "failure")
	serverCmd.Flags().StringVar(&sWorkingDirectory, "working-directory", "", "working directory")
	serverCmd.Flags().StringVar(&sLogMode, "log-mode", "roll", "log mode")
	serverCmd.Flags().StringVar(&sLogPattern, "log-pattern", "", "log pattern")
	serverCmd.Flags().StringVar(&sLogAutoRollAtTime, "log-auto-roll-at-time", "", "log auto roll at time")
	serverCmd.Flags().StringVar(&sLogSizeThreshold, "log-size-threshold", "", "log size threshold")
	serverCmd.Flags().StringVar(&sLogZipOlderThanNumDays, "log-zip-older-than-num-days", "", "log zip older than num days")
	serverCmd.Flags().StringVar(&sLogZipDateFormat, "log-zip-date-format", "", "log zip date format")
	serverCmd.Flags().BoolVar(&sForce, "force", true, "force write")

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
		if sName == "" {
			return fmt.Errorf("missing name")
		}
		if sId == "" {
			sId = sName
		}
		if sExecutable == "" {
			return fmt.Errorf("missing executable")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s := win.NewServer(
			win.WithSId(sId),
			win.WithSName(sName),
			win.WithSExecutable(sExecutable),
			win.WithSDescription(sDescription),
			win.WithSStartMode(sStartMode),
			win.WithSDepends(sDepends),
			win.WithSLogPath(sLogPath),
			win.WithSArguments(sArguments),
			win.WithSStartArguments(sStartArguments),
			win.WithSStopExecutable(sStopExecutable),
			win.WithSStopArguments(sStopArguments),
			win.WithSEnv(sEnv),
			win.WithSFailure(sFailure),
			win.WithSWorkingDirectory(sWorkingDirectory),
			win.WithSLogMode(sLogMode),
			win.WithSLogPattern(sLogPattern),
			win.WithSLogAutoRollAtTime(sLogAutoRollAtTime),
			win.WithSLogSizeThreshold(sLogSizeThreshold),
			win.WithSLogZipOlderThanNumDays(sLogZipOlderThanNumDays),
			win.WithSLogZipDateFormat(sLogZipDateFormat),
			win.WithSForce(sForce),
		)
		err := s.Run()
		if err != nil {
			return err
		}
		return nil
	},
}
