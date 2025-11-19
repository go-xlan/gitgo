package gitgo

// Status displays the working tree state and staged changes
// Shows modified files, untracked files, and staging area contents
// Use case: check repo state before committing to understand current changes
//
// Status 显示工作树状态和暂存更改
// 显示修改的文件、未跟踪的文件和暂存区内容
// 使用场景：在提交前检查仓库状态以了解当前更改
func (G *Gcm) Status() *Gcm {
	return G.do("git", "status")
}

// Add stages changes in the working path to prepare the next commit
// Includes modified files, new files, and removed files in staging area
// Use case: prepare changes when committing in automated workflows
//
// Add 暂存工作路径中的更改以准备下次提交
// 在暂存区包含修改的文件、新文件和删除的文件
// 使用场景：在自动化工作流中准备更改以进行提交
func (G *Gcm) Add() *Gcm {
	return G.do("git", "add", ".")
}

// Commit creates a new commit with the provided message text
// Saves staged changes to the repo with descriptive message
// Use case: save staged changes as a snapshot with descriptive text
// Note: fails when no changes are staged
//
// Commit 使用提供的消息文本创建新的提交
// 将暂存更改保存到仓库并附带描述性消息
// 使用场景：将暂存更改保存为快照并附带描述性文本
// 注意：当没有暂存更改时会失败
func (G *Gcm) Commit(m string) *Gcm {
	return G.do("git", "commit", "-m", m) // Fails when no staged changes exist // 当没有待提交文件时会失败
}

// Fetch and merge and merges changes from the remote repo
// Downloads recent commits from remote and integrates them into the branch
// Use case: sync branch with remote changes when starting new work session
//
// Pull 从远程仓库获取并合并更改
// 从远程下载最新提交并将其集成到分支
// 使用场景：在新工作会话开始时将分支与远程更改同步
func (G *Gcm) Pull() *Gcm {
	return G.do("git", "pull")
}

// Push uploads commits to the remote repo
// Sends branch commits to the matching remote branch
// Use case: share work with team and backup commits to remote repo
//
// Push 将提交上传到远程仓库
// 将分支提交发送到对应的远程分支
// 使用场景：与团队共享工作和将提交备份到远程仓库
func (G *Gcm) Push() *Gcm {
	return G.do("git", "push")
}

// PushSetUpstreamOriginBranch pushes new branch and sets upstream tracking
// Creates remote branch and establishes tracking connections to handle push and fetch
// Use case: publish new branch to remote and enable simple push and fetch commands
//
// PushSetUpstreamOriginBranch 推送新分支并设置上游跟踪
// 创建远程分支并建立跟踪连接以处理推送和拉取
// 使用场景：将新分支发布到远程并启用简单的推送和拉取命令
func (G *Gcm) PushSetUpstreamOriginBranch(newBranchName string) *Gcm {
	return G.do("git", "push", "--set-upstream", "origin", newBranchName)
}

// Reset unstages changes and keeps them in working path
// Moves staged changes back to working path without losing modifications
// Use case: undo staging operations to adjust files before commit
//
// Reset 取消暂存更改但保留在工作路径中
// 将暂存的更改移回工作路径而不丢失修改
// 使用场景：撤销暂存操作以在提交前调整文件
func (G *Gcm) Reset() *Gcm {
	return G.do("git", "reset")
}

// ResetHard discards changes and resets to recent commit
// DANGEROUS: removes uncommitted changes in working path and staging area
// Use case: abandon work in progress and return to clean state
//
// ResetHard 丢弃更改并重置到最近提交
// 危险：删除工作路径和暂存区中未提交的更改
// 使用场景：放弃进行中的工作并返回到干净状态
func (G *Gcm) ResetHard() *Gcm {
	return G.do("git", "reset", "--hard")
}

// Checkout switches to an existing branch or commit
// Changes the working path to match the specified branch or commit state
// Use case: switch between development branches or examine past commits
//
// Checkout 切换到现有分支或提交
// 更改工作路径以匹配指定分支或提交状态
// 使用场景：在开发分支间切换或检查过去提交
func (G *Gcm) Checkout(nameBranch string) *Gcm {
	return G.do("git", "checkout", nameBranch)
}

// CheckoutNewBranch creates and switches to a new branch
// Creates new branch from the HEAD position and switches to the branch
// Use case: start new feature development and create test branches
//
// CheckoutNewBranch 创建并切换到新分支
// 从 HEAD 位置创建新分支并切换到该分支
// 使用场景：开始新功能开发和创建测试分支
func (G *Gcm) CheckoutNewBranch(nameBranch string) *Gcm {
	return G.do("git", "checkout", "-b", nameBranch)
}

// Init initializes a new Git repo in the path
// Creates .git path and sets up the structure
// Use case: start version management on new projects and convert existing ones
//
// Init 在路径中初始化新的 Git 仓库
// 创建 .git 路径并设置结构
// 使用场景：在新项目上开始版本管理和转换现有项目
func (G *Gcm) Init() *Gcm {
	return G.do("git", "init")
}

// Merge integrates changes from specified branch into the branch
// Combines commit records from feature branch with the branch
// Use case: integrate completed features into main development branches
//
// Merge 将指定分支的更改集成到分支
// 将功能分支的提交记录与分支合并
// 使用场景：将完成的功能集成到主开发分支
func (G *Gcm) Merge(featureBranchName string) *Gcm {
	return G.do("git", "merge", featureBranchName)
}

// MergeAbort cancels an in-progress merge operation
// Restores the state when encountering merge attempt conflicts
// Use case: abandon problematic merges and attempt with different tactics
//
// MergeAbort 取消正在进行的合并操作
// 在遇到合并尝试冲突时恢复状态
// 使用场景：放弃有问题的合并并用不同策略尝试
func (G *Gcm) MergeAbort() *Gcm {
	return G.do("git", "merge", "--abort")
}

// TagList shows existing tags in the repo
// Displays version tags and release points as a complete list
// Use case: examine release records and find specific version tags
//
// TagList 显示仓库中现有标签
// 显示版本标签和发布点的完整列表
// 使用场景：检查发布记录和查找特定版本标签
func (G *Gcm) TagList() *Gcm {
	return G.do("git", "tag", "--list")
}

// Tags shows existing tags in the repo (TagList name alternative)
// Displays version tags and release points as a complete list
// Use case: examine release records and find specific version tags
//
// Tags 显示仓库中现有标签（TagList 的名称替代）
// 显示版本标签和发布点的完整列表
// 使用场景：检查发布记录和查找特定版本标签
func (G *Gcm) Tags() *Gcm {
	return G.do("git", "tag", "--list")
}

// Tag creates a new tag at the commit position
// Marks the commit with version naming and release indications
// Use case: mark release points and important milestones in development
//
// Tag 在提交位置创建新标签
// 用版本命名和发布指示标记提交
// 使用场景：标记发布点和开发中的重要里程碑
func (G *Gcm) Tag(tag string) *Gcm {
	return G.do("git", "tag", tag)
}

// PushTags uploads tags to remote repo
// Synchronizes tag information with remote when team shares
// Use case: share version tags and release points with team members
//
// PushTags 将标签上传到远程仓库
// 在团队共享时与远程同步标签信息
// 使用场景：与团队成员共享版本标签和发布点
func (G *Gcm) PushTags() *Gcm {
	return G.do("git", "push", "--tags")
}

// PushTag uploads specific tag to remote repo
// Shares one tag with remote without affecting alternate tags
// Use case: publish specific release tag without pushing tags in bulk
//
// PushTag 将特定标签上传到远程仓库
// 与远程共享一个标签而不影响其他标签
// 使用场景：发布特定发布标签而不批量推送标签
func (G *Gcm) PushTag(tag string) *Gcm {
	return G.do("git", "push", "origin", tag)
}
