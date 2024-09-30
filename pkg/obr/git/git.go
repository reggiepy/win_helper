package git

import (
	"fmt"
	"os"
	"os/exec"
)

// TagAndPush 封装了 Git 打标签和推送标签的逻辑
//
// 参数:
//   - tagName: 要创建的 Git 标签名称 (如 "v2.2.0")。
//   - tagMessage: 标签的附加说明信息 (如 "Release version 2.2.0")。
//
// 返回:
//   - error: 如果执行过程中出现错误，返回错误信息；如果执行成功，返回 nil。
//
// 流程:
//   1. 使用 `git tag -a <tagName> -m <tagMessage>` 创建带注释的 Git 标签。
//   2. 使用 `git push origin --tags` 将标签推送到远程仓库。
func TagAndPush(tagName, tagMessage string) error {
	// 执行 `git tag -a <tagName> -m <tagMessage>`
	cmdTag := exec.Command("git", "tag", "-a", tagName, "-m", tagMessage)
	cmdTag.Stdout = os.Stdout
	cmdTag.Stderr = os.Stderr
	if err := cmdTag.Run(); err != nil {
		return fmt.Errorf("error running git tag: %v", err)
	}

	// 执行 `git push origin --tags`
	cmdPush := exec.Command("git", "push", "origin", "--tags")
	cmdPush.Stdout = os.Stdout
	cmdPush.Stderr = os.Stderr
	if err := cmdPush.Run(); err != nil {
		return fmt.Errorf("error running git push: %v", err)
	}

	fmt.Println("Git tag and push successful.")
	return nil
}

// CommitChanges 提交更改到 Git 仓库
//
// 参数:
//   - commitMessage: Git 提交时的说明信息。
//   - commitFiles: 可选的文件列表，表示需要提交的文件。如果为空，默认添加全部文件 (`git add .`)。
//
// 返回:
//   - error: 如果执行过程中出现错误，返回错误信息；如果执行成功，返回 nil。
//
// 流程:
//   1. 使用 `git add` 添加指定文件到暂存区。如果未提供文件，则添加全部更改。
//   2. 使用 `git commit -m <commitMessage>` 提交更改到 Git 仓库。
func CommitChanges(commitMessage string, commitFiles ...string) error {
	// 判断是否有需要提交的指定文件
	if len(commitFiles) > 0 {
		// 逐个添加指定文件
		for _, commitFile := range commitFiles {
			cmdAdd := exec.Command("git", "add", commitFile)
			cmdAdd.Stdout = os.Stdout
			cmdAdd.Stderr = os.Stderr
			if err := cmdAdd.Run(); err != nil {
				return fmt.Errorf("error running git add: %v", err)
			}
		}
	} else {
		// 没有指定文件，默认添加全部文件
		cmdAdd := exec.Command("git", "add", ".")
		cmdAdd.Stdout = os.Stdout
		cmdAdd.Stderr = os.Stderr
		if err := cmdAdd.Run(); err != nil {
			return fmt.Errorf("error running git add: %v", err)
		}
	}

	// 提交更改
	cmdCommit := exec.Command("git", "commit", "-m", commitMessage)
	cmdCommit.Stdout = os.Stdout
	cmdCommit.Stderr = os.Stderr
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("error running git commit: %v", err)
	}

	fmt.Println("Git commit successful.")
	return nil
}
