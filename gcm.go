package gitgo

// Status displays the working tree status and staged changes
// Shows modified files, untracked files, and staging area contents
// Use case: check repository state before commit or to understand current changes
//
// Status 显示工作树状态和暂存更改
// 显示修改的文件、未跟踪的文件和暂存区内容
// 使用场景：提交前检查仓库状态或了解当前更改
func (G *Gcm) Status() *Gcm {
	return G.do("git", "status")
}

// Add stages all changes in the working DIR for the next commit
// Includes modified files, new files, and deleted files in staging area
// Use case: prepare all changes for commit in automated workflows
//
// Add 将工作目录中的所有更改暂存到下次提交
// 在暂存区包含修改的文件、新文件和删除的文件
// 使用场景：自动化工作流中准备所有更改进行提交
func (G *Gcm) Add() *Gcm {
	return G.do("git", "add", ".")
}

// Commit creates a new commit with the provided message
// Commits all staged changes to the local repository with descriptive message
// Use case: save staged changes as a permanent snapshot with meaningful description
// Note: will fail if no changes are staged
//
// Commit 使用提供的消息创建新的提交
// 将所有暂存更改提交到本地仓库并附带描述性消息
// 使用场景：将暂存更改保存为永久快照并附带有意义的描述
// 注意：如果没有暂存更改将会失败
func (G *Gcm) Commit(m string) *Gcm {
	return G.do("git", "commit", "-m", m) //当没有待提交文件时，这里也会报错
}

// Pull fetches and merges changes from the remote repository
// Downloads latest commits from remote and integrates them into current branch
// Use case: sync local branch with remote changes before starting new work
//
// Pull 从远程仓库获取并合并更改
// 从远程下载最新提交并将其集成到当前分支
// 使用场景：开始新工作前将本地分支与远程更改同步
func (G *Gcm) Pull() *Gcm {
	return G.do("git", "pull")
}

// Push uploads local commits to the remote repository
// Sends local branch commits to the corresponding remote branch
// Use case: share local work with team or backup commits to remote server
//
// Push 将本地提交上传到远程仓库
// 将本地分支提交发送到对应的远程分支
// 使用场景：与团队共享本地工作或将提交备份到远程服务器
func (G *Gcm) Push() *Gcm {
	return G.do("git", "push")
}

// PushSetUpstreamOriginBranch pushes new branch and sets upstream tracking
// Creates remote branch and establishes tracking relationship for future push/pull
// Use case: publish new local branch to remote and enable simple push/pull commands
//
// PushSetUpstreamOriginBranch 推送新分支并设置上游跟踪
// 创建远程分支并建立跟踪关系以便将来推送/拉取
// 使用场景：将新的本地分支发布到远程并启用简单的推送/拉取命令
func (G *Gcm) PushSetUpstreamOriginBranch(newBranchName string) *Gcm {
	return G.do("git", "push", "--set-upstream", "origin", newBranchName)
}

// Reset unstages changes while keeping them in working DIR
// Moves staged changes back to working DIR without losing modifications
// Use case: undo staging operations when you want to modify files before commit
//
// Reset 取消暂存更改但保留在工作目录中
// 将暂存的更改移回工作目录而不丢失修改
// 使用场景：撤销暂存操作，希望在提交前修改文件
func (G *Gcm) Reset() *Gcm {
	return G.do("git", "reset")
}

// ResetHard discards all changes and resets to last commit
// DANGEROUS: permanently deletes all uncommitted changes in working DIR and staging area
// Use case: completely abandon current work and return to clean repository state
//
// ResetHard 丢弃所有更改并重置到上次提交
// 危险：永久删除工作目录和暂存区中所有未提交的更改
// 使用场景：完全放弃当前工作并返回到干净的仓库状态
func (G *Gcm) ResetHard() *Gcm {
	return G.do("git", "reset", "--hard")
}

// Checkout switches to an existing branch or commit
// Changes the working DIR to match the specified branch or commit state
// Use case: switch between development branches or examine historical commits
//
// Checkout 切换到现有分支或提交
// 更改工作目录以匹配指定分支或提交状态
// 使用场景：在开发分支间切换或检查历史提交
func (G *Gcm) Checkout(nameBranch string) *Gcm {
	return G.do("git", "checkout", nameBranch)
}

// CheckoutNewBranch creates and switches to a new branch
// Creates new branch from current HEAD and immediately switches to it
// Use case: start new feature development or create experimental branches
//
// CheckoutNewBranch 创建并切换到新分支
// 从当前 HEAD 创建新分支并立即切换到该分支
// 使用场景：开始新功能开发或创建实验性分支
func (G *Gcm) CheckoutNewBranch(nameBranch string) *Gcm {
	return G.do("git", "checkout", "-b", nameBranch)
}

// Init initializes a new Git repository in the current DIR
// Creates .git DIR and sets up initial repository structure
// Use case: start version control for new projects or convert existing projects
//
// Init 在当前目录中初始化新的 Git 仓库
// 创建 .git 目录并设置初始仓库结构
// 使用场景：为新项目开始版本控制或转换现有项目
func (G *Gcm) Init() *Gcm {
	return G.do("git", "init")
}

// Merge integrates changes from specified branch into current branch
// Combines commit history from feature branch with current branch
// Use case: integrate completed feature development into main development line
//
// Merge 将指定分支的更改集成到当前分支
// 将功能分支的提交历史与当前分支合并
// 使用场景：将完成的功能开发集成到主开发线
func (G *Gcm) Merge(featureBranchName string) *Gcm {
	return G.do("git", "merge", featureBranchName)
}

// MergeAbort cancels an in-progress merge operation
// Restores repository to state before merge attempt when conflicts occur
// Use case: abandon problematic merge and retry with different strategy
//
// MergeAbort 取消正在进行的合并操作
// 当发生冲突时将仓库恢复到合并尝试前的状态
// 使用场景：放弃有问题的合并并用不同策略重试
func (G *Gcm) MergeAbort() *Gcm {
	return G.do("git", "merge", "--abort")
}

// TagList displays all existing tags in the repository
// Shows complete list of version tags and release markers
// Use case: review release history or find specific version tags
//
// TagList 显示仓库中所有现有标签
// 显示版本标签和发布标记的完整列表
// 使用场景：查看发布历史或查找特定版本标签
func (G *Gcm) TagList() *Gcm {
	return G.do("git", "tag", "--list")
}

// Tags displays all existing tags in the repository (alias for TagList)
// Shows complete list of version tags and release markers
// Use case: review release history or find specific version tags
//
// Tags 显示仓库中所有现有标签（TagList 的别名）
// 显示版本标签和发布标记的完整列表
// 使用场景：查看发布历史或查找特定版本标签
func (G *Gcm) Tags() *Gcm {
	return G.do("git", "tag", "--list")
}

// Tag creates a new lightweight tag at current commit
// Marks current commit with version label or release identifier
// Use case: mark release points or important milestones in development
//
// Tag 在当前提交处创建新的轻量级标签
// 用版本标签或发布标识符标记当前提交
// 使用场景：标记发布点或开发中的重要里程碑
func (G *Gcm) Tag(tag string) *Gcm {
	return G.do("git", "tag", tag)
}

// PushTags uploads all local tags to remote repository
// Synchronizes all tag information with remote for team sharing
// Use case: share version tags and release markers with team members
//
// PushTags 将所有本地标签上传到远程仓库
// 与远程同步所有标签信息以便团队共享
// 使用场景：与团队成员共享版本标签和发布标记
func (G *Gcm) PushTags() *Gcm {
	return G.do("git", "push", "--tags")
}

// PushTag uploads specific tag to remote repository
// Shares single tag with remote without affecting other tags
// Use case: publish specific release tag without pushing all local tags
//
// PushTag 将特定标签上传到远程仓库
// 与远程共享单个标签而不影响其他标签
// 使用场景：发布特定发布标签而不推送所有本地标签
func (G *Gcm) PushTag(tag string) *Gcm {
	return G.do("git", "push", "origin", tag)
}
