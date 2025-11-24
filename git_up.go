package gitgo

import (
	"fmt"
	"regexp"
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
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "HAS-CHANGES").ExecTake("git", "diff-index", "--cached", "--quiet", "HEAD")
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 1:
		return true, nil // Has staged changes // 有暂存更改
	case 0:
		return false, nil // No staged changes // 无暂存更改
	default:
		return false, erero.Errorf("git diff-index failed with exit code %d", exc)
	}
}

// HasUnstagedChanges checks if working tree has unstaged modifications
// Returns true upon detecting unstaged changes, false if working tree matches staging area
// Use case: check staging needs at commit operation time
//
// HasUnstagedChanges 检查工作树是否有未暂存的修改
// 在检测到未暂存更改时返回 true，如果工作树与暂存区匹配则返回 false
// 使用场景：在提交操作时检查暂存需求
func (G *Gcm) HasUnstagedChanges() (bool, error) {
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "HAS-CHANGES").ExecTake("git", "diff", "--quiet")
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 1:
		return true, nil // Has unstaged changes // 有未暂存更改
	case 0:
		return false, nil // No unstaged changes // 无未暂存更改
	default:
		return false, erero.Errorf("git diff failed with exit code %d", exc)
	}
}

// HasChanges checks if changes exist (staged and unstaged)
// Returns true upon detecting modifications anywhere, false if repo is clean
// Use case: quick check on work in progress when switching context
//
// HasChanges 检查是否存在更改（已暂存和未暂存）
// 在检测到任何修改时返回 true，如果仓库干净则返回 false
// 使用场景：在上下文切换时快速检查进行中的工作
func (G *Gcm) HasChanges() (bool, error) {
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "HAS-CHANGES").ExecTake("git", "diff-index", "--quiet", "HEAD")
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 1:
		return true, nil // Has changes // 有更改
	case 0:
		return false, nil // No changes // 无更改
	default:
		return false, erero.Errorf("git diff-index failed with exit code %d", exc)
	}
}

// GetPorcelainStatus checks if uncommitted changes exist in repo
// Returns clean status if repo has no staged and unstaged changes
// Use case: check clean state in important operations like branch switching and releases
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
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "HAS-STAGED-CHANGES").ExecTake("git", "diff-index", "--cached", "--quiet", "HEAD")
	if err != nil {
		return newWaGcm(G.execConfig, []byte{}, err, G.debugMode)
	}
	switch exc {
	case 1:
		return G // Has staged changes // 有暂存的更改
	case 0:
		return newWaGcm(G.execConfig, []byte{}, errors.New("NON-STAGED-CHANGES"), G.debugMode)
	default:
		return newWaGcm(G.execConfig, []byte{}, errors.Errorf("git diff-index failed with exit code %d", exc), G.debugMode)
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
		return "", erero.Wro(err) // Wrap and return failure // 包装并返回错误
	}
	return strings.TrimSpace(string(output)), nil
}

// GetLatestTag retrieves the latest tag name if it exists
// Returns tag name, exists flag, and possible errors
// When no tags exist, returns ("", false, nil)
//
// GetLatestTag 获取最新标签名称（如果存在）
// 返回标签名称、存在标志和可能的错误
// 当没有标签时，返回 ("", false, nil)
func (G *Gcm) GetLatestTag() (string, bool, error) {
	output, exc, err := G.execConfig.NewConfig().WithExpectExit(128, "NO-TAGS").ExecTake("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", false, erero.Wro(err)
	}
	switch exc {
	case 0:
		return strings.TrimSpace(string(output)), true, nil // Tag exists // 标签存在
	case 128:
		return "", false, nil // No tags exist // 没有标签
	default:
		return "", false, erero.Errorf("git describe failed with exit code %d", exc)
	}
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
	// Validate prefix to prevent command injection // 验证prefix安全性，避免命令注入
	if strings.Contains(prefix, "'") || strings.Contains(prefix, "`") || strings.Contains(prefix, "$") {
		return "", erero.New("prefix contains unsafe characters")
	}

	// Use git tag -l to list tags matching prefix, version desc sort, take first // 使用 git tag -l 列出匹配前缀的标签，按版本号降序排序，取第一个
	output, err := G.execConfig.NewConfig().WithBash().Exec(fmt.Sprintf("git tag -l --sort=-version:refname '%s*' | head -n 1", prefix))
	if err != nil {
		return "", erero.Wro(err)
	}

	return strings.TrimSpace(string(output)), nil
}

// LatestGitTagMatchRegexp retrieves the latest tag matching glob pattern
// Returns most recent tag that matches the given glob pattern, blank if none
// Supports wildcards in subproject and component versioning
//
// LatestGitTagMatchRegexp 获取匹配 shell glob 模式的最新标签
// 返回匹配给定 glob 模式的最新标签和空字符串
// 在子项目和组件的版本管理期间支持通配符模式
func (G *Gcm) LatestGitTagMatchRegexp(regexpPattern string) (string, error) {
	if regexpPattern == "" {
		return "", erero.New("regexpPattern is required")
	}
	// Validate pattern to prevent command injection // 验证regexpPattern安全性，避免命令注入
	if strings.Contains(regexpPattern, "'") || strings.Contains(regexpPattern, "`") || strings.Contains(regexpPattern, "$") {
		return "", erero.New("regexpPattern contains unsafe characters")
	}

	// Use git tag -l with glob pattern, version desc sort, take first // 使用 git tag -l 匹配 glob 模式，按版本号降序排序，取第一个
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
// Returns tags with creation dates sorted ascending as formatted string
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
// Returns the top path, fails if not in a Git repo
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
	// Use --absolute-git-dir instead of --git-dir to get absolute path // 使用 --absolute-git-dir 替代 --git-dir 以获得更好的可用性
	output, err := G.execConfig.Exec("git", "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetSubPathToRoot retrieves path from current location to base
// Returns path like "../" if in subdirs, blank if at base
// Use case: build paths to base assets
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
// Returns path like "subpath/" if in subdirs, blank if at base
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
// Returns the current branch name, fails if not in a Git repo
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

// GetRemoteURL gets the address of the specified remote repo
// Returns the remote address, fails if remote does not exist
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

// GetCommitCount gets the commit count in the current branch
// Returns the commit count, fails if not in a Git repo or no commits exist
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

// ListBranches gets branch names in the repo
// Returns slice of branch names, fails if not in a Git repo
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

// ListRemoteBranches gets remote branch names in the repo
// Returns slice of remote branch names, fails if not in a Git repo
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
// Returns slice of one-line commit entries, fails if not in a Git repo
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

// GetCurrentCommitHash gets the hash of HEAD commit
// Returns commit hash string from Git command output
// Use case: locate commit position when recording state and creating references
//
// GetCurrentCommitHash 获取 HEAD 提交的哈希值
// 从 Git 命令输出返回提交哈希字符串
// 使用场景：在记录状态和创建引用时识别提交位置
func (G *Gcm) GetCurrentCommitHash() (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "HEAD")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCommitMessage gets the commit message of specified reference
// Returns commit message text with subject and message
// Use case: examine commit contents when reviewing changes and generating notes
//
// GetCommitMessage 获取指定引用的提交消息
// 返回包含主题和消息的提交消息文本
// 使用场景：在审查更改和生成注释时检查提交内容
func (G *Gcm) GetCommitMessage(ref string) (string, error) {
	output, err := G.execConfig.Exec("git", "log", "-1", "--pretty=format:%B", ref)
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// BranchExists checks if branch exists in repo
// Returns true if branch exists, false otherwise
// Use case: validate branch names when switching and creating branches
//
// BranchExists 检查仓库中是否存在分支
// 如果分支存在则返回 true，否则返回 false
// 使用场景：在切换和创建分支时验证分支名称
func (G *Gcm) BranchExists(name string) (bool, error) {
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "NOT-EXIST").ExecTake("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name)
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 0:
		return true, nil // Branch exists // 分支存在
	case 1:
		return false, nil // Branch does not exist // 分支不存在
	default:
		return false, erero.Errorf("git show-ref failed with exit code %d", exc)
	}
}

// RemoteBranchExists checks if remote branch exists
// Returns true if remote branch exists, false otherwise
// Use case: validate remote branch references when fetching and tracking
//
// RemoteBranchExists 检查远程分支是否存在
// 如果远程分支存在则返回 true，否则返回 false
// 使用场景：在获取和跟踪时验证远程分支引用
func (G *Gcm) RemoteBranchExists(name string) (bool, error) {
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "NOT-EXIST").ExecTake("git", "show-ref", "--verify", "--quiet", "refs/remotes/"+name)
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 0:
		return true, nil // Remote branch exists // 远程分支存在
	case 1:
		return false, nil // Remote branch does not exist // 远程分支不存在
	default:
		return false, erero.Errorf("git show-ref failed with exit code %d", exc)
	}
}

// TagExists checks if tag exists in repo
// Returns true if tag exists, false otherwise
// Use case: prevent duplicate tag creation and validate tag references
//
// TagExists 检查仓库中是否存在标签
// 如果标签存在则返回 true，否则返回 false
// 使用场景：防止重复创建标签和验证标签引用
func (G *Gcm) TagExists(name string) (bool, error) {
	_, exc, err := G.execConfig.NewConfig().WithExpectExit(1, "NOT-EXIST").ExecTake("git", "show-ref", "--verify", "--quiet", "refs/tags/"+name)
	if err != nil {
		return false, erero.Wro(err)
	}
	switch exc {
	case 0:
		return true, nil // Tag exists // 标签存在
	case 1:
		return false, nil // Tag does not exist // 标签不存在
	default:
		return false, erero.Errorf("git show-ref failed with exit code %d", exc)
	}
}

// GetFileList gets tracked files in the repo
// Returns paths of files tracked with Git
// Use case: examine repo contents and validate file existence when processing assets
//
// GetFileList 获取仓库中的跟踪文件
// 返回 Git 跟踪的文件路径
// 使用场景：检查仓库内容并在处理资源时验证文件存在
func (G *Gcm) GetFileList() ([]string, error) {
	output, err := G.execConfig.Exec("git", "ls-files")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// GetUntrackedFiles gets files not tracked with Git
// Returns paths of files in working path but not in version management
// Use case: find new files when staging changes and cleaning workspace
//
// GetUntrackedFiles 获取未被 Git 跟踪的文件
// 返回工作路径中但不在版本管理中的文件路径
// 使用场景：在暂存更改和清理工作空间时识别新文件
func (G *Gcm) GetUntrackedFiles() ([]string, error) {
	output, err := G.execConfig.Exec("git", "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// GetModifiedFiles gets files with uncommitted modifications
// Returns paths of files changed in working path and staging area
// Use case: find affected files when reviewing changes and staging
//
// GetModifiedFiles 获取有未提交修改的文件
// 返回工作路径和暂存区中已更改文件的路径
// 使用场景：在审查更改和选择性暂存时识别受影响的文件
func (G *Gcm) GetModifiedFiles() ([]string, error) {
	output, err := G.execConfig.Exec("git", "diff", "--name-only", "HEAD")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// GetBranchTrackingBranch gets the upstream branch of specified branch
// Returns upstream branch name in format "remote/branch"
// Use case: check tracking configuration and understand remote connections
//
// GetBranchTrackingBranch 获取指定分支的上游分支
// 返回格式为 "remote/branch" 的上游分支名称
// 使用场景：验证跟踪配置和了解远程连接
func (G *Gcm) GetBranchTrackingBranch(branch string) (string, error) {
	output, err := G.execConfig.Exec("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	if err != nil {
		return "", erero.Wro(err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetIgnoredFiles gets files git ignores in the repo
// Returns paths of files matching gitignore rules
// Use case: find ignored files when cleaning workspace and checking config
//
// GetIgnoredFiles 获取仓库中被 git 忽略的文件
// 返回匹配 gitignore 规则的文件路径
// 使用场景：在清理工作空间和检查配置时识别被忽略的文件
func (G *Gcm) GetIgnoredFiles() ([]string, error) {
	output, err := G.execConfig.Exec("git", "status", "--ignored", "-s", "--", ".")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var regexpIgnore = regexp.MustCompile(`^!!\s*(.+)$`)
	var paths []string
	for _, item := range strings.Split(string(output), "\n") {
		item = strings.TrimSpace(item)
		if m := regexpIgnore.FindStringSubmatch(item); len(m) > 0 {
			paths = append(paths, m[1])
		}
	}
	return paths, nil
}
