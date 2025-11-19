package gitgo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/yyle88/erero"
)

// HasStagingChanges checks if staged changes are waiting to commit
// Returns true upon finding staged changes, false otherwise
// Use case: avoid commit operations without staged changes to prevent issues and blank commits
//
// HasStagingChanges 检查是否有暂存的更改等待提交
// 在找到暂存更改时返回 true，否则返回 false
// 使用场景：避免在没有暂存更改时执行 commit 操作，防止问题和产生空提交
func (G *Gcm) HasStagingChanges() (bool, error) {
	const (
		nonChanges = "non-changes"
		hasChanges = "has-changes"
	)

	output, err := G.execConfig.NewConfig().WithBash().
		Exec(`git diff-index --cached --quiet HEAD
case $? in
    0) echo "` + nonChanges + `" ;;
    1) echo "` + hasChanges + `" ;;
    *) exit $? ;;
esac`)
	if err != nil {
		return false, erero.Wro(err)
	}
	return strings.TrimSpace(string(output)) == hasChanges, nil
}

// HasUnstagedChanges checks if working tree has unstaged modifications
// Returns true upon detecting unstaged changes, false if working tree matches staging area
// Use case: check staging needs at commit operation time
//
// HasUnstagedChanges 检查工作树是否有未暂存的修改
// 在检测到未暂存更改时返回 true，如果工作树与暂存区匹配则返回 false
// 使用场景：在提交操作时检查暂存需求
func (G *Gcm) HasUnstagedChanges() (bool, error) {
	const (
		nonChanges = "non-changes"
		hasChanges = "has-changes"
	)

	output, err := G.execConfig.NewConfig().WithBash().
		Exec(`git diff --quiet
case $? in
    0) echo "` + nonChanges + `" ;;
    1) echo "` + hasChanges + `" ;;
    *) exit $? ;;
esac`)
	if err != nil {
		return false, erero.Wro(err)
	}
	return strings.TrimSpace(string(output)) == hasChanges, nil
}

// HasChanges checks if changes exist (staged and unstaged)
// Returns true upon detecting modifications anywhere, false if repo is clean
// Use case: quick check on work in progress when switching context
//
// HasChanges 检查是否存在更改（已暂存和未暂存）
// 在检测到任何修改时返回 true，如果仓库干净则返回 false
// 使用场景：在上下文切换时快速检查进行中的工作
func (G *Gcm) HasChanges() (bool, error) {
	const (
		nonChanges = "non-changes"
		hasChanges = "has-changes"
	)

	output, err := G.execConfig.NewConfig().WithBash().
		Exec(`git diff-index --quiet HEAD
case $? in
    0) echo "` + nonChanges + `" ;;
    1) echo "` + hasChanges + `" ;;
    *) exit $? ;;
esac`)
	if err != nil {
		return false, erero.Wro(err)
	}
	return strings.TrimSpace(string(output)) == hasChanges, nil
}

// GetPorcelainStatus checks if uncommitted changes exist in repo
// Returns clean status if repo has no staged and unstaged changes
// Use case: check clean state during key operations like branch switching and releases
//
// GetPorcelainStatus 检查仓库中是否存在未提交的更改
// 如果仓库没有已暂存和未暂存更改则返回干净状态
// 使用场景：在分支切换和发布等关键操作期间检查干净状态
func (G *Gcm) GetPorcelainStatus() (string, error) {
	output, err := G.execConfig.Exec("git", "status", "--porcelain")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CheckStagedChanges validates staged changes and returns issues if none exist
// Returns Gcm instance with problem state upon finding no staged changes
// Use case: prevent commit operations if working path is clean
//
// CheckStagedChanges 验证暂存更改，如果没有更改则返回问题
// 在没有找到暂存更改时返回带有问题状态的 Gcm 实例
// 使用场景：如果工作路径干净则防止 commit 操作
func (G *Gcm) CheckStagedChanges() *Gcm {
	const (
		nonStagedChanges = "non-staged-changes"
		hasStagedChanges = "has-staged-changes"
	)

	output, err := G.execConfig.NewConfig().WithBash().
		Exec(`git diff-index --cached --quiet HEAD
case $? in
    0) echo "` + nonStagedChanges + `"; exit 0 ;;
    1) echo "` + hasStagedChanges + `"; exit 0 ;;
    *) exit $? ;;
esac`)
	if err != nil {
		return newWaGcm(G.execConfig, output, err, G.debugMode)
	}
	// Check staged changes based on output content // 根据输出内容判断是否有暂存更改
	state := strings.TrimSpace(string(output))
	switch state {
	case hasStagedChanges:
		return G
	case nonStagedChanges:
		return newWaGcm(G.execConfig, []byte{}, errors.New(nonStagedChanges), G.debugMode)
	default:
		return newWaGcm(G.execConfig, output, errors.Errorf("UNEXPECTED OUTPUT: %s", state), G.debugMode)
	}
}

// LatestGitTag retrieves the name of the most recent tag in the project
// Returns the latest tag name and issues if no tags exist with command issues
// Uses git describe command to find the most recent annotated and lightweight tag
//
// LatestGitTag 获取项目中最新的标签名称
// 如果没有标签存在和命令问题则返回问题和最新标签名称
// 使用 git describe 命令查找最新的注释标签和轻量级标签
func (G *Gcm) LatestGitTag() (string, error) {
	output, err := G.execConfig.Exec("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", erero.Wro(err) // Wrap and return error // 包装并返回错误
	}
	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagHasPrefix retrieves the latest tag name with the specified prefix
// Returns the most recent tag that starts with the given prefix
// Helps find version tags on specific subprojects and components
//
// LatestGitTagHasPrefix 获取带有指定前缀的最新标签名称
// 返回以给定前缀开头的最新标签
// 帮助在特定子项目和组件上查找版本标签
func (G *Gcm) LatestGitTagHasPrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", erero.New("prefix is required")
	}
	// Validate prefix security to prevent command injection // 验证prefix安全性，避免命令注入
	if strings.Contains(prefix, "'") || strings.Contains(prefix, "`") || strings.Contains(prefix, "$") {
		return "", erero.New("prefix contains unsafe characters")
	}

	// Use git tag -l to list tags matching prefix, sorted by version desc, take first // 使用 git tag -l 列出匹配前缀的标签，按版本号降序排序，取第一个
	output, err := G.execConfig.NewConfig().WithBash().Exec(fmt.Sprintf("git tag -l --sort=-version:refname '%s*' | head -n 1", prefix))
	if err != nil {
		return "", erero.Wro(err)
	}

	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagMatchRegexp retrieves the latest tag matching shell glob pattern
// Returns most recent tag that matches the given glob pattern and empty string
// Supports wildcard patterns during subproject and component versioning
//
// LatestGitTagMatchRegexp 获取匹配 shell glob 模式的最新标签
// 返回匹配给定 glob 模式的最新标签和空字符串
// 在子项目和组件的版本管理期间支持通配符模式
func (G *Gcm) LatestGitTagMatchRegexp(regexpPattern string) (string, error) {
	if regexpPattern == "" {
		return "", erero.New("regexpPattern is required")
	}
	// Validate pattern security to prevent command injection // 验证regexpPattern安全性，避免命令注入
	if strings.Contains(regexpPattern, "'") || strings.Contains(regexpPattern, "`") || strings.Contains(regexpPattern, "$") {
		return "", erero.New("regexpPattern contains unsafe characters")
	}

	// Use git tag -l with glob pattern, sorted by version desc, take first // 使用 git tag -l 匹配 glob 模式，按版本号降序排序，取第一个
	output, err := G.execConfig.NewConfig().WithBash().Exec(fmt.Sprintf("git tag -l --sort=-version:refname '%s' | head -n 1", regexpPattern))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GitCommitHash retrieves the commit hash on a specified branch and tag reference
// Returns the complete commit hash string with the given reference name
// Supports branch names, tag names, and Git references
//
// GitCommitHash 获取指定分支和标签引用的提交哈希
// 返回给定引用名称的完整提交哈希字符串
// 支持分支名、标签名和 Git 引用
func (G *Gcm) GitCommitHash(refName string) (string, error) {
	if refName == "" {
		return "", erero.New("refName is required")
	}
	output, err := G.execConfig.Exec("git", "rev-parse", refName)
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// SortedGitTags retrieves sorted list of project tags with dates
// Returns tags with creation dates in ascending order as formatted string
// Use case: examine tag content to choose next version numbering
//
// SortedGitTags 获取项目标签的排序列表及日期
// 返回按创建日期升序排列的标签和日期格式化字符串
// 使用场景：检查标签内容以选择下一个版本编号
func (G *Gcm) SortedGitTags() (string, error) {
	output, err := G.execConfig.Exec("git", "for-each-ref", "--sort=creatordate", "--format=%(refname) %(creatordate)", "refs/tags")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetTopPath retrieves the absolute path of Git repo at base location
// Returns the top-level path and error if not in a Git repo
// Use case: navigate to project base and resolve paths
//
// GetTopPath 获取 Git 仓库在基础位置的绝对路径
// 如果不在 Git 仓库则返回顶层路径和错误
// 使用场景：导航到项目基础和解析路径
func (G *Gcm) GetTopPath() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitDIRAbsPath retrieves absolute path of the .git location
// Returns complete path to .git location (e.g., "/home/admin/project/.git")
// Use case: access Git metadata, hooks, and configuration files
//
// GetGitDIRAbsPath 获取 .git 位置的绝对路径
// 返回 .git 位置的完整路径（如 "/home/admin/project/.git"）
// 使用场景：访问 Git 元数据、钩子和配置文件
func (G *Gcm) GetGitDIRAbsPath() (string, error) {
	// Use --absolute-git-dir instead of --git-dir for better usability // 使用 --absolute-git-dir 替代 --git-dir 以获得更好的可用性
	output, err := G.execConfig.Exec("git", "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPathToRoot retrieves path from current location to base
// Returns path like "../" if in subdirs and empty string if at base
// Use case: build paths to base-level assets
//
// GetSubPathToRoot 获取从当前位置到基础的路径
// 如果在子目录则返回如 "../" 的路径，如果在基础则返回空字符串
// 使用场景：构建到基础级别资源的路径
func (G *Gcm) GetSubPathToRoot() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--show-cdup")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPath retrieves path from base to current location
// Returns path like "subpath/" if in subdirs and empty string if at base
// Use case: find current location within project structure
//
// GetSubPath 获取从基础到当前位置的路径
// 如果在子目录则返回如 "subpath/" 的路径，如果在基础则返回空字符串
// 使用场景：查找在项目结构中的当前位置
func (G *Gcm) GetSubPath() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--show-prefix")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// IsInsideWorkTree checks if the current path is inside a Git work tree
// Returns true if inside a Git project, false otherwise
//
// IsInsideWorkTree 检查当前路径是否在 Git 工作树中
// 如果在 Git 项目中则返回 true，否则返回 false
func (G *Gcm) IsInsideWorkTree() (bool, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, erero.Wro(err)
	}
	return strings.TrimSpace(string(output)) == "true", nil
}

// GetCurrentBranch gets the name of the current branch
// Returns the current branch name and error if not in a Git repo
//
// GetCurrentBranch 获取当前分支的名称
// 如果不在 Git 仓库中则返回当前分支名称和错误
func (G *Gcm) GetCurrentBranch() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRemoteURL gets the URL of the specified remote repo
// Returns the remote URL and error if remote does not exist
//
// GetRemoteURL 获取指定远程仓库的URL
// 如果远程仓库不存在则返回远程URL和错误
func (G *Gcm) GetRemoteURL(remoteName string) (string, error) {
	if remoteName == "" {
		remoteName = "origin"
	}
	output, err := G.execConfig.Exec("git", "remote", "get-url", remoteName)
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCommitCount gets the total count of commits in the current branch
// Returns the commit count and error if not in a Git repo or no commits exist
//
// GetCommitCount 获取当前分支的提交总数
// 如果不在 Git 仓库中或没有提交则返回提交数量和错误
func (G *Gcm) GetCommitCount() (int, error) {
	output, err := G.execConfig.Exec("git", "rev-list", "--count", "HEAD")
	if err != nil {
		return 0, erero.Wro(err)
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, erero.Wro(err)
	}
	return count, nil
}

// ListBranches gets a list of all local branch names
// Returns slice of branch names and error if not in a Git repo
//
// ListBranches 获取所有本地分支名称列表
// 如果不在 Git 仓库中则返回分支名称切片和错误
func (G *Gcm) ListBranches() ([]string, error) {
	output, err := G.execConfig.Exec("git", "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, erero.Wro(err)
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var results []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			results = append(results, branch)
		}
	}
	return results, nil
}

// ListRemoteBranches gets a list of all remote branch names
// Returns slice of remote branch names and error if not in a Git repo
//
// ListRemoteBranches 获取所有远程分支名称列表
// 如果不在 Git 仓库中则返回远程分支名称切片和错误
func (G *Gcm) ListRemoteBranches() ([]string, error) {
	output, err := G.execConfig.Exec("git", "branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return nil, erero.Wro(err)
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var results []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" && !strings.Contains(branch, "HEAD") {
			results = append(results, branch)
		}
	}
	return results, nil
}

// GetLogOneLine gets a concise commit log with specified limit
// Returns slice of one-line commit entries and error if not in a Git repo
//
// GetLogOneLine 获取指定数量的简洁提交日志
// 如果不在 Git 仓库中则返回单行提交条目切片和错误
func (G *Gcm) GetLogOneLine(limit int) ([]string, error) {
	if limit <= 0 {
		return nil, erero.New("limit must > 0")
	}
	if limit >= 10000 {
		return nil, erero.New("limit must < 10000")
	}
	output, err := G.execConfig.Exec("git", "log", "--oneline", fmt.Sprintf("-n%d", limit))
	if err != nil {
		return nil, erero.Wro(err)
	}
	var results []string
	for _, outLine := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		outLine = strings.TrimSpace(outLine)
		if outLine != "" {
			results = append(results, outLine)
		}
	}
	return results, nil
}
