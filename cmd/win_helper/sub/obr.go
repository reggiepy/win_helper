package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"win_helper/pkg/obr/app"
	"win_helper/pkg/obr/git"
	"win_helper/pkg/obr/iss"
	"win_helper/pkg/util/versionUtils/manager"
)

type OBRUpdateAppCmdConfig struct {
	IssPath string
	Version string

	PushGit    bool
	GitMessage string
}

type OBRUpdateISSCmdConfig struct {
	IssPath string
	Version string

	PushGit    bool
	GitMessage string
}

var (
	obrUpdateISSCmdConfig = &OBRUpdateISSCmdConfig{}
	obrUpdateAppCmdConfig = &OBRUpdateAppCmdConfig{}
)

func init() {
	rootCmd.AddCommand(obrCmd)
	obrCmd.AddCommand(updateAppCmd)
	obrCmd.AddCommand(updateISSCmd)

	updateISSCmd.Flags().StringVarP(&obrUpdateISSCmdConfig.Version, "version", "v", "+", "git version message")
	updateISSCmd.Flags().StringVarP(&obrUpdateISSCmdConfig.IssPath, "iss-path", "", "C:\\dist\\chemical_server.iss", "iss_path")

	updateISSCmd.Flags().BoolVar(&obrUpdateISSCmdConfig.PushGit, "push-git", false, "push git tag")
	updateISSCmd.Flags().StringVar(&obrUpdateISSCmdConfig.GitMessage, "git-message", "", "git version message")

	updateAppCmd.Flags().StringVarP(&obrUpdateAppCmdConfig.Version, "version", "v", "+", "git version message")
	updateAppCmd.Flags().StringVarP(&obrUpdateAppCmdConfig.IssPath, "iss-path", "", "C:\\dist\\chemical_server.iss", "iss_path")

	updateAppCmd.Flags().BoolVarP(&obrUpdateAppCmdConfig.PushGit, "push-git", "", false, "push git tag")
	updateAppCmd.Flags().StringVarP(&obrUpdateAppCmdConfig.GitMessage, "git-message", "m", "", "git version message")
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

// 更新 ISS 版本号子命令
var updateISSCmd = &cobra.Command{
	Use:   "update-iss",
	Short: "Update ISS version",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取当前目录
		originalDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %v", err)
		}
		fmt.Println("Current directory:", originalDir)

		// 获取绝对路径
		issDir := filepath.Dir(obrUpdateISSCmdConfig.IssPath)
		fmt.Println("Changing to ISS directory:", issDir)

		// 切换到 ISS 文件所在目录
		if err := os.Chdir(issDir); err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}

		// 你的更新 ISS 版本的逻辑
		fmt.Println("ISS version updated.")

		version, err := iss.GetCurrentVersion(obrUpdateISSCmdConfig.IssPath)
		if err != nil {
			return err
		}
		versionManager := manager.NewVersionManager(manager.WithVersion(version))
		currentVersion := versionManager.GetVersion()

		if err := versionManager.SetVersion(obrUpdateISSCmdConfig.Version); err != nil {
			return fmt.Errorf("error updating version: %v", err)
		}
		newVersion := versionManager.GetVersion()
		fmt.Printf("update version %s ---> %s", currentVersion, newVersion)
		if err := iss.SaveVersion(newVersion, obrUpdateISSCmdConfig.IssPath); err != nil {
			return fmt.Errorf("error saving version: %v", err)
		}

		if obrUpdateISSCmdConfig.PushGit {
			commitMessage := fmt.Sprintf("update version %s --> %s", currentVersion, newVersion)
			if obrUpdateISSCmdConfig.GitMessage != "" {
				commitMessage = obrUpdateISSCmdConfig.GitMessage
			}
			if err := git.CommitChanges(commitMessage); err != nil {
				return err
			}
			if err := git.TagAndPush(newVersion, obrUpdateISSCmdConfig.GitMessage); err != nil {
				return err
			}
		}

		// 切换回原始目录
		if err := os.Chdir(originalDir); err != nil {
			return fmt.Errorf("error changing back to original directory: %v", err)
		}
		return nil
	},
}

// 更新 App 版本号子命令
var updateAppCmd = &cobra.Command{
	Use:   "update-app",
	Short: "Update App version",
	RunE: func(cmd *cobra.Command, args []string) error {
		originalDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %v", err)
		}
		version, err := app.GetCurrentVersion(originalDir)
		if err != nil {
			return err
		}
		versionManager := manager.NewVersionManager(manager.WithVersion(version))
		currentVersion := versionManager.GetVersion()

		if err := versionManager.SetVersion(obrUpdateAppCmdConfig.Version); err != nil {
			return fmt.Errorf("error updating version: %v", err)
		}
		newVersion := versionManager.GetVersion()
		fmt.Printf("update version %s ---> %s", currentVersion, newVersion)
		versionFile := filepath.Join(originalDir, "VERSION")
		if err := versionManager.SaveAs(versionFile, true); err != nil {
			return fmt.Errorf("error saving version: %v", err)
		}

		if obrUpdateAppCmdConfig.PushGit {
			commitMessage := fmt.Sprintf("update version %s --> %s", currentVersion, newVersion)
			if obrUpdateAppCmdConfig.GitMessage != "" {
				commitMessage = obrUpdateAppCmdConfig.GitMessage
			}
			if err := git.CommitChanges(commitMessage); err != nil {
				return err
			}
			if err := git.TagAndPush(newVersion, obrUpdateISSCmdConfig.GitMessage); err != nil {
				return err
			}
		}
		return nil
	},
}
