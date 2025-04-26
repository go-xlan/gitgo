package gitgo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/yyle88/erero"
)

// HasStagingChanges 查看是否有变动，使用场景是：没有变动就不要执行 commit 否则会报错或者产生空提交
func (G *Gcm) HasStagingChanges() (bool, error) {
	if output, err := G.execConfig.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(output) != 0 {
			return false, errors.New(string(output))
		}
		return true, nil
	}
	return false, nil
}

// CheckStagedChanges 查看是否有变动，使用场景是：没有变动就不要执行 commit 否则会报错或者产生空提交
func (G *Gcm) CheckStagedChanges() *Gcm {
	if output, err := G.execConfig.Exec("git", "diff", "--cached", "--quiet"); err == nil {
		if len(output) != 0 {
			return newWaGcm(G.execConfig, output, err, G.debugMode)
		}
		return newWaGcm(G.execConfig, []byte{}, errors.New("no-staged-changes"), G.debugMode)
	}
	return G
}

// LatestGitTag 获得项目的最新标签的名称
func (G *Gcm) LatestGitTag() (string, error) {
	const commandBash = `
res=$(git describe --tags --abbrev=0 2>/dev/null)
if [ -z "$res" ]; then
  echo ""
else
  echo "$res"
fi
`
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagHasPrefix 获得指定前缀的子项目或主项目的最新标签名称
func (G *Gcm) LatestGitTagHasPrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", erero.New("param prefix is none")
	}

	// Bash 命令：获取匹配前缀的最新标签。这里去掉 if [ -z "$res" ] 分支，因为 echo "$res" 在 res 为空时会自动输出空字符串
	const commandBashTemplate = `
res=$(git for-each-ref --sort=creatordate --format '%%(refname:short)' refs/tags | grep "^%s" | tail -n 1)
echo "$res"
`
	// Bash 命令：使用转义后的 prefix 格式化命令，获得完整的命令字符串
	commandBash := fmt.Sprintf(commandBashTemplate, regexp.QuoteMeta(prefix))

	// 存在问题，就是标签是打在commit提交上面的，因此在相同的提交上打多个标签时，他们的时间是相同的，这时候取到的，有可能不是最新的那个标签
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagMatchRegexp 获得匹配正则(shell的glob模式正则)的子项目或主项目的最新标签名称
func (G *Gcm) LatestGitTagMatchRegexp(regexpPattern string) (string, error) {
	if regexpPattern == "" {
		return "", erero.New("param regexp_pattern is none")
	}

	// Bash 命令：获取匹配前缀的最新标签。这里去掉 if [ -z "$res" ] 分支，因为 echo "$res" 在 res 为空时会自动输出空字符串
	const commandBashTemplate = `
res=$(git for-each-ref --sort=creatordate --format '%%(refname:short)' refs/tags/%s | tail -n 1)
echo "$res"
`
	// Bash 命令：使用转义后的 prefix 格式化命令，获得完整的命令字符串
	commandBash := fmt.Sprintf(commandBashTemplate, regexp.QuoteMeta(regexpPattern))

	// 存在问题，就是标签是打在commit提交上面的，因此在相同的提交上打多个标签时，他们的时间是相同的，这时候取到的，有可能不是最新的那个标签
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GitCommitHash 获得某个分支或标签的哈希
func (G *Gcm) GitCommitHash(refName string) (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", refName)
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// SortedGitTags 获取项目的标签列表，按时间从前到后排列。使用场景是仅仅观察标签，让用户自己想出下一个标签的序号，因此这里只返回个字符串就行
func (G *Gcm) SortedGitTags() (string, error) {
	const commandBash = "git for-each-ref --sort=creatordate --format '%(refname) %(creatordate)' refs/tags"
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetTopPath 获取 git 项目的根目录
func (G *Gcm) GetTopPath() (string, error) {
	const commandBash = "git rev-parse --show-toplevel"
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitDirAbsPath 获取 git 项目的 .git 目录的绝对路径（如 "/home/user/project/.git"）
func (G *Gcm) GetGitDirAbsPath() (string, error) {
	const commandBash = "git rev-parse --absolute-git-dir" // 这里不要使用 git rev-parse --git-dir 稍微有点不好用
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPathToRoot 获取从当前目录到 git 项目根目录的相对路径（如 "../"，如果在根目录则为空字符串）
func (G *Gcm) GetSubPathToRoot() (string, error) {
	const commandBash = "git rev-parse --show-cdup"
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPath 获取从 git 项目根目录到当前目录的相对路径（如 "subdir"，如果在根目录则为空字符串）
func (G *Gcm) GetSubPath() (string, error) {
	const commandBash = "git rev-parse --show-prefix"
	output, err := G.execConfig.ShallowClone().WithBash().Exec(strings.TrimSpace(commandBash))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}
