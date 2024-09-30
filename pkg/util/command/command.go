package command

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand 执行系统命令，并返回标准输出和标准错误输出。
func RunCommand(args ...string) (string, error) {
	var cmd *exec.Cmd
	var err error
	cmdStr := strings.Join(args, " ")
	fmt.Printf("执行命令: %s\n", cmdStr)
	cmd = exec.Command(args[0], args[1:]...) // 将命令和参数分开传递
	out, err := cmd.CombinedOutput()
	fmt.Printf("执行命令输出: %s\n", string(out))
	if err != nil {
		return "", err
	}
	return string(out), nil
}
