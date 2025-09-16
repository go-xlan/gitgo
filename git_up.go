package gitgo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/yyle88/erero"
)

// HasStagingChanges checks if there are staged changes ready for commit
// Returns true if there are staged changes, false otherwise
// Use case: avoid commit operations when no changes are staged to prevent errors or empty commits
//
// HasStagingChanges 检查是否有暂存的更改等待提交
// 如果有暂存更改返回 true，否则返回 false
// 使用场景：没有变动时避免执行 commit 操作，防止报错或产生空提交
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

// HasUnstagedChanges checks if working tree has any unstaged modifications
// Returns true if there are unstaged changes, false if working tree matches staging area
// Use case: determine if files need to be staged before commit
//
// HasUnstagedChanges 检查工作树是否有未暂存的修改
// 如果有未暂存更改返回 true，如果工作树与暂存区匹配返回 false
// 使用场景：确定提交前是否需要暂存文件
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

// HasChanges checks if repository has any changes at all (staged or unstaged)
// Returns true if any modifications exist anywhere, false if repository is completely clean
// Use case: quick check if there's any work in progress before switching contexts
//
// HasChanges 检查仓库是否有任何更改（已暂存或未暂存）
// 如果任何地方存在修改返回 true，如果仓库完全干净返回 false
// 使用场景：切换上下文前快速检查是否有进行中的工作
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

// GetPorcelainStatus checks if repository has any uncommitted changes
// Returns true if repository is completely clean (no staged or unstaged changes), false otherwise
// Use case: verify clean state before critical operations like branch switching or releases
//
// GetPorcelainStatus 检查仓库是否有任何未提交的更改
// 如果仓库完全干净（无已暂存或未暂存更改）返回 true，否则返回 false
// 使用场景：在分支切换或发布等关键操作前验证干净状态
func (G *Gcm) GetPorcelainStatus() (string, error) {
	output, err := G.execConfig.Exec("git", "status", "--porcelain")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CheckStagedChanges validates staged changes and returns error if none exist
// Returns Gcm instance with error state if no staged changes found
// Use case: prevent commit operations when working DIR is clean
//
// CheckStagedChanges 验证暂存更改，如果没有更改则返回错误
// 如果没有找到暂存更改，返回带有错误状态的 Gcm 实例
// 使用场景：工作目录干净时防止 commit 操作
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
	// 根据输出内容判断是否有暂存更改
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
// Returns the latest tag name or error if no tags exist or command fails
// Uses git describe command to find the most recent annotated or lightweight tag
//
// LatestGitTag 获取项目中最新的标签名称
// 返回最新标签名称，如果没有标签存在或命令失败则返回错误
// 使用 git describe 命令查找最新的注释标签或轻量级标签
func (G *Gcm) LatestGitTag() (string, error) {
	output, err := G.execConfig.Exec("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", erero.Wro(err) // Wrap and return error // 包装并返回错误
	}
	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagHasPrefix retrieves the latest tag name with the specified prefix
// Returns the most recent tag that starts with the given prefix
// Useful for finding version tags for specific subprojects or components
//
// LatestGitTagHasPrefix 获取带有指定前缀的最新标签名称
// 返回以给定前缀开头的最新标签
// 适用于查找特定子项目或组件的版本标签
func (G *Gcm) LatestGitTagHasPrefix(prefix string) (string, error) {
	if prefix == "" {
		return "", erero.New("prefix is required")
	}
	// 验证prefix安全性，避免命令注入
	if strings.Contains(prefix, "'") || strings.Contains(prefix, "`") || strings.Contains(prefix, "$") {
		return "", erero.New("prefix contains unsafe characters")
	}

	// 使用 git tag -l 列出匹配前缀的标签，--sort=-version:refname 按版本号降序排序，head -n 1 取第一个
	// Use git tag -l to list tags, sorted by version desc, head -n 1 to get the latest
	output, err := G.execConfig.NewConfig().WithBash().Exec(fmt.Sprintf("git tag -l --sort=-version:refname '%s*' | head -n 1", prefix))
	if err != nil {
		return "", erero.Wro(err)
	}

	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagMatchRegexp retrieves the latest tag matching shell glob pattern
// Returns most recent tag that matches the given glob pattern or empty string
// Supports wildcard patterns for versioning subprojects or components
//
// LatestGitTagMatchRegexp 获取匹配 shell glob 模式的最新标签
// 返回匹配给定 glob 模式的最新标签或空字符串
// 支持通配符模式用于子项目或组件的版本管理
func (G *Gcm) LatestGitTagMatchRegexp(regexpPattern string) (string, error) {
	if regexpPattern == "" {
		return "", erero.New("regexpPattern is required")
	}
	// 验证regexpPattern安全性，避免命令注入
	if strings.Contains(regexpPattern, "'") || strings.Contains(regexpPattern, "`") || strings.Contains(regexpPattern, "$") {
		return "", erero.New("regexpPattern contains unsafe characters")
	}

	// 使用 git tag -l 匹配 glob 模式，--sort=-version:refname 按版本号降序排序，head -n 1 取第一个
	// Use git tag -l with glob pattern, sorted by version desc, head -n 1 to get the latest
	output, err := G.execConfig.NewConfig().WithBash().Exec(fmt.Sprintf("git tag -l --sort=-version:refname '%s' | head -n 1", regexpPattern))
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GitCommitHash retrieves the commit hash for a specified branch or tag reference
// Returns the full commit hash string for the given reference name
// Supports branch names, tag names, and other Git references
//
// GitCommitHash 获取指定分支或标签引用的提交哈希
// 返回给定引用名称的完整提交哈希字符串
// 支持分支名、标签名和其他 Git 引用
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

// SortedGitTags retrieves chronologically sorted list of project tags
// Returns tags with creation dates in ascending order as formatted string
// Use case: review tag history to determine next version number manually
//
// SortedGitTags 获取项目标签的时间顺序列表
// 返回按创建日期升序排列的标签和日期格式化字符串
// 使用场景：查看标签历史以手动确定下一个版本号
func (G *Gcm) SortedGitTags() (string, error) {
	output, err := G.execConfig.Exec("git", "for-each-ref", "--sort=creatordate", "--format=%(refname) %(creatordate)", "refs/tags")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetTopPath retrieves the absolute path of Git repository root DIR
// Returns the top-level DIR path or error if not in a Git repository
// Use case: navigate to project root or resolve relative paths
//
// GetTopPath 获取 Git 仓库根目录的绝对路径
// 返回顶层目录路径，如果不在 Git 仓库中则返回错误
// 使用场景：导航到项目根目录或解析相对路径
func (G *Gcm) GetTopPath() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitDIRAbsPath retrieves absolute path of the .git DIR
// Returns full path to .git DIR (e.g., "/home/user/project/.git")
// Use case: access Git metadata, hooks, or configuration files
//
// GetGitDIRAbsPath 获取 .git 目录的绝对路径
// 返回 .git 目录的完整路径（如 "/home/user/project/.git"）
// 使用场景：访问 Git 元数据、钩子或配置文件
func (G *Gcm) GetGitDIRAbsPath() (string, error) {
	// 这里不要使用 git rev-parse --git-dir 稍微有点不好用
	output, err := G.execConfig.Exec("git", "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPathToRoot retrieves relative path from current DIR to repository root
// Returns path like "../" for subdirectories or empty string if at root
// Use case: construct relative paths to root-level resources
//
// GetSubPathToRoot 获取从当前目录到仓库根目录的相对路径
// 子目录返回如 "../" 的路径，如果在根目录则返回空字符串
// 使用场景：构造到根级别资源的相对路径
func (G *Gcm) GetSubPathToRoot() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--show-cdup")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPath retrieves relative path from repository root to current DIR
// Returns path like "subdir/" for subdirectories or empty string if at root
// Use case: identify current location within project structure
//
// GetSubPath 获取从仓库根目录到当前目录的相对路径
// 子目录返回如 "subdir/" 的路径，如果在根目录则返回空字符串
// 使用场景：识别在项目结构中的当前位置
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
// 如果在 Git 项目中返回 true，否则返回 false
func (G *Gcm) IsInsideWorkTree() (bool, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false, erero.Wro(err)
	}
	return strings.TrimSpace(string(output)) == "true", nil
}

// GetCurrentBranch gets the name of the current branch
// Returns the current branch name or error if not in a Git repository
//
// GetCurrentBranch 获取当前分支的名称
// 返回当前分支名称，如果不在 Git 仓库中则返回错误
func (G *Gcm) GetCurrentBranch() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRemoteURL gets the URL of the specified remote repository
// Returns the remote URL or error if remote does not exist
//
// GetRemoteURL 获取指定远程仓库的URL
// 返回远程URL，如果远程仓库不存在则返回错误
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

// GetCommitCount gets the total number of commits in the current branch
// Returns the commit count or error if not in a Git repository or no commits
//
// GetCommitCount 获取当前分支的提交总数
// 返回提交数量，如果不在 Git 仓库中或没有提交则返回错误
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
// Returns slice of branch names or error if not in a Git repository
//
// ListBranches 获取所有本地分支名称列表
// 返回分支名称切片，如果不在 Git 仓库中则返回错误
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
// Returns slice of remote branch names or error if not in a Git repository
//
// ListRemoteBranches 获取所有远程分支名称列表
// 返回远程分支名称切片，如果不在 Git 仓库中则返回错误
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
// Returns slice of one-line commit entries or error if not in a Git repository
//
// GetLogOneLine 获取指定数量的简洁提交日志
// 返回单行提交条目切片，如果不在 Git 仓库中则返回错误
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
