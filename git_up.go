package gitgo

import (
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
