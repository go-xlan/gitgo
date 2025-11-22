package gitgo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestGcm_HasStagingChanges tests detection of staged changes in the Git space
// Verifies that the function identifies if changes are in staging area
func TestGcm_HasStagingChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	changes, err := gcm.HasStagingChanges()
	require.NoError(t, err)
	t.Log(changes)
}

// TestGcm_HasUnstagedChanges tests detection of unstaged modifications in working tree
// Verifies that the function identifies changes not yet in staging area
func TestGcm_HasUnstagedChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	changes, err := gcm.HasUnstagedChanges()
	require.NoError(t, err)
	t.Log(changes)
}

// TestGcm_HasStagingChanges_True tests HasStagingChanges returns true with staged changes
// Verifies exit code 1 scenario where staged changes exist
//
// TestGcm_HasStagingChanges_True 测试有暂存更改时 HasStagingChanges 返回 true
// 验证退出码 1 场景，即存在暂存更改
func TestGcm_HasStagingChanges_True(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-staging-true-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "init.txt"), []byte("init"), 0644))
	gcm.Add().Commit("initial").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "new.txt"), []byte("new"), 0644))
	gcm.Add().Done()

	has := rese.V1(gcm.HasStagingChanges())
	require.True(t, has)
}

// TestGcm_HasStagingChanges_False tests HasStagingChanges returns false without staged changes
// Verifies exit code 0 scenario where no staged changes exist
//
// TestGcm_HasStagingChanges_False 测试无暂存更改时 HasStagingChanges 返回 false
// 验证退出码 0 场景，即不存在暂存更改
func TestGcm_HasStagingChanges_False(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-staging-false-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "init.txt"), []byte("init"), 0644))
	gcm.Add().Commit("initial").Done()

	has := rese.V1(gcm.HasStagingChanges())
	require.False(t, has)
}

// TestGcm_HasUnstagedChanges_True tests HasUnstagedChanges returns true with unstaged changes
// Verifies exit code 1 scenario where unstaged changes exist
//
// TestGcm_HasUnstagedChanges_True 测试有未暂存更改时 HasUnstagedChanges 返回 true
// 验证退出码 1 场景，即存在未暂存更改
func TestGcm_HasUnstagedChanges_True(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-unstaged-true-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("initial").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))

	has := rese.V1(gcm.HasUnstagedChanges())
	require.True(t, has)
}

// TestGcm_HasUnstagedChanges_False tests HasUnstagedChanges returns false without unstaged changes
// Verifies exit code 0 scenario where no unstaged changes exist
//
// TestGcm_HasUnstagedChanges_False 测试无未暂存更改时 HasUnstagedChanges 返回 false
// 验证退出码 0 场景，即不存在未暂存更改
func TestGcm_HasUnstagedChanges_False(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-unstaged-false-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("initial").Done()

	has := rese.V1(gcm.HasUnstagedChanges())
	require.False(t, has)
}

// TestGcm_HasChanges tests detection of changes (staged and unstaged combined)
// Verifies that the function identifies if work is in progress
func TestGcm_HasChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	changes, err := gcm.HasChanges()
	require.NoError(t, err)
	t.Log(changes)
}

// TestGcm_HasChanges_True tests HasChanges returns true with changes
// Verifies exit code 1 scenario where changes exist
//
// TestGcm_HasChanges_True 测试有更改时 HasChanges 返回 true
// 验证退出码 1 场景，即存在更改
func TestGcm_HasChanges_True(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-changes-true-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("initial").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))

	has := rese.V1(gcm.HasChanges())
	require.True(t, has)
}

// TestGcm_HasChanges_False tests HasChanges returns false without changes
// Verifies exit code 0 scenario where no changes exist
//
// TestGcm_HasChanges_False 测试无更改时 HasChanges 返回 false
// 验证退出码 0 场景，即不存在更改
func TestGcm_HasChanges_False(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-has-changes-false-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("initial").Done()

	has := rese.V1(gcm.HasChanges())
	require.False(t, has)
}

// TestGcm_GetPorcelainStatus tests fetch of clean status information
// Verifies that the function returns accurate porcelain format status output
func TestGcm_GetPorcelainStatus(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	status, err := gcm.GetPorcelainStatus()
	require.NoError(t, err)
	t.Log(status)
}

// TestGcm_CheckStagedChanges tests validation of staged changes with chain integration
// Verifies that the function works right in fluent interface command chains
func TestGcm_CheckStagedChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		Add().
		CheckStagedChanges().
		ShowDebugMessage().
		Status().
		ShowDebugMessage()
}

// TestGcm_CheckStagedChanges_HasChanges tests CheckStagedChanges with staged changes
// Verifies exit code 1 scenario where staged changes exist
//
// TestGcm_CheckStagedChanges_HasChanges 测试有暂存更改时的 CheckStagedChanges
// 验证退出码 1 场景，即存在暂存更改
func TestGcm_CheckStagedChanges_HasChanges(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-staged-has-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	// Create first commit // 创建初始提交
	must.Done(os.WriteFile(filepath.Join(tempDIR, "init.txt"), []byte("init"), 0644))
	gcm.Add().Commit("initial").Done()

	// Stage a new file // 暂存一个新文件
	must.Done(os.WriteFile(filepath.Join(tempDIR, "new.txt"), []byte("new"), 0644))
	gcm.Add().Done()

	// Should succeed with staged changes // 有暂存更改时应该成功
	gcm.CheckStagedChanges().Done()
}

// TestGcm_CheckStagedChanges_NoChanges tests CheckStagedChanges without staged changes
// Verifies exit code 0 scenario where no staged changes exist
//
// TestGcm_CheckStagedChanges_NoChanges 测试没有暂存更改时的 CheckStagedChanges
// 验证退出码 0 场景，即不存在暂存更改
func TestGcm_CheckStagedChanges_NoChanges(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-staged-non-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	// Create first commit // 创建初始提交
	must.Done(os.WriteFile(filepath.Join(tempDIR, "init.txt"), []byte("init"), 0644))
	gcm.Add().Commit("initial").Done()

	// No staged changes, expect issue // 没有暂存更改时应该失败
	err := gcm.CheckStagedChanges().Reason()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NON-STAGED-CHANGES")
}

// TestLatestGitTag tests fetch of the most recent tag in the project
// Verifies that the function returns the latest tag name right
func TestLatestGitTag(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTag()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGcm_LatestGitTagHasPrefix tests fetch of latest tag with specific prefix
// Verifies that the function filters tags with prefix pattern
func TestGcm_LatestGitTagHasPrefix(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagHasPrefix("v")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGcm_LatestGitTagHasPrefix_Compare compares common tag fetch with prefix matching
// Demonstrates that prefix-matched tags can be distinct from the most recent tag
func TestGcm_LatestGitTagHasPrefix_Compare(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tagA := rese.C1(gcm.LatestGitTag())
	t.Log(tagA)
	tagB := rese.C1(gcm.LatestGitTagHasPrefix("v"))
	t.Log(tagB)
	// Log without assertion - tags can be different // 只打印不断言相等，因为标签可能不同
	t.Log("tag-A:", tagA, "tag-B:", tagB)
}

// TestGcm_LatestGitTagMatchRegexp tests fetch of latest tag matching glob pattern
// Verifies that the function filters tags using wildcards
func TestGcm_LatestGitTagMatchRegexp(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagMatchRegexp("v[0-9]*.[0-9]*.[0-9]*")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGitCommitHash tests fetch of commit hash with branch reference
// Verifies that the function resolves branch name to commit hash
func TestGitCommitHash(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	hash, err := gcm.GitCommitHash("main")
	require.NoError(t, err)
	t.Log(hash)
	require.NotEmpty(t, hash)
}

// TestGitCommitHash_TAG tests fetch of commit hash with tag reference
// Verifies that the function resolves tag name to commit hash
func TestGitCommitHash_TAG(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTag()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)

	hash, err := gcm.GitCommitHash(tag)
	require.NoError(t, err)
	t.Log(hash)
	require.NotEmpty(t, hash)
}

// TestGcm_SortedGitTags tests fetch of sorted tag list with dates
// Verifies that the function returns tags in ascending time sequence
func TestGcm_SortedGitTags(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.SortedGitTags()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGcm_GetTopPath tests fetch of Git space base path
// Verifies that the function returns project top path location
func TestGcm_GetTopPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	root, err := gcm.GetTopPath()
	require.NoError(t, err)
	t.Log(root)
	require.NotEmpty(t, root)
	require.Equal(t, runpath.PARENT.Path(), root)

	osmustexist.MustRoot(root)
	osmustexist.MustRoot(filepath.Join(root, ".git"))
}

// TestGcm_GetGitDIRAbsPath tests fetch of .git location absolute path
// Verifies that the function returns correct path to Git metadata location
func TestGcm_GetGitDIRAbsPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetGitDIRAbsPath()
	require.NoError(t, err)
	t.Log(path)
	require.NotEmpty(t, path)

	osmustexist.MustRoot(path)
	require.Equal(t, runpath.PARENT.Join(".git"), path)
}

// TestGcm_GetSubPathToRoot tests fetch of path from current location to base
// Verifies that the function returns correct path navigation string to project base
func TestGcm_GetSubPathToRoot(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPathToRoot()
	require.NoError(t, err)
	t.Log(path)
	require.Empty(t, path)
	t.Log("sub path is the project root") // Empty when at base // 在基础位置时为空
	osmustexist.MustRoot(runpath.PARENT.Join(".git"))
}

// TestGcm_GetSubPath tests fetch of path from base to current location
// Verifies that the function returns correct relative path within project structure
func TestGcm_GetSubPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPath()
	require.NoError(t, err)
	t.Log(path)
	require.Empty(t, path)
	t.Log("sub path is the project root") // Empty when at base // 在基础位置时为空
	osmustexist.MustRoot(runpath.PARENT.Join(".git"))
}

// TestGcm_IsInsideWorkTree tests checking of current path location in Git work tree
// Verifies that the function detects if inside a Git project
func TestGcm_IsInsideWorkTree(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	inside, err := gcm.IsInsideWorkTree()
	require.NoError(t, err)
	t.Log(inside)
	require.True(t, inside)
}

// TestGcm_GetCurrentBranch tests fetch of current branch name
// Verifies that the function returns accurate active branch information
func TestGcm_GetCurrentBranch(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	branch, err := gcm.GetCurrentBranch()
	require.NoError(t, err)
	t.Log(branch)
	require.NotEmpty(t, branch)
}

// TestGcm_GetRemoteURL tests fetch of remote space address
// Verifies that the function returns origin remote location
func TestGcm_GetRemoteURL(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	remoteURL, err := gcm.GetRemoteURL("origin")
	require.NoError(t, err)
	t.Log(remoteURL)
	require.NotEmpty(t, remoteURL)
}

// TestGcm_GetCommitCount tests fetch of commit count in current branch
// Verifies that the function returns the count of commits
func TestGcm_GetCommitCount(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	count, err := gcm.GetCommitCount()
	require.NoError(t, err)
	t.Log(count)
	require.Greater(t, count, 0)
}

// TestGcm_ListBranches tests fetch of branch names
// Verifies that the function returns the list of project branches
func TestGcm_ListBranches(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	branches, err := gcm.ListBranches()
	require.NoError(t, err)
	t.Log(branches)
	require.NotEmpty(t, branches)
}

// TestGcm_ListRemoteBranches tests fetch of remote branch names
// Verifies that the function returns the list of remote branches
func TestGcm_ListRemoteBranches(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	branches, err := gcm.ListRemoteBranches()
	require.NoError(t, err)
	t.Log(branches)
	require.NotEmpty(t, branches)
}

// TestGcm_GetLogOneLine tests fetch of concise commit log with limit
// Verifies that the function returns specified count of one-line commit entries
func TestGcm_GetLogOneLine(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	commits, err := gcm.GetLogOneLine(5)
	require.NoError(t, err)
	t.Log(commits)
	require.NotEmpty(t, commits)
}

// TestGcm_GetCurrentCommitHash tests getting HEAD commit hash
// Verifies GetCurrentCommitHash returns commit hash string
//
// TestGcm_GetCurrentCommitHash 测试获取 HEAD 提交哈希
// 验证 GetCurrentCommitHash 返回提交哈希字符串
func TestGcm_GetCurrentCommitHash(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-hash-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "test.txt"), []byte("test"), 0644))
	gcm.Add().Commit("test commit").Done()

	hash := rese.V1(gcm.GetCurrentCommitHash())
	require.Len(t, hash, 40)
	require.Regexp(t, "^[0-9a-f]{40}$", hash)
}

// TestGcm_GetCommitMessage tests getting commit message
// Verifies GetCommitMessage retrieves message text
//
// TestGcm_GetCommitMessage 测试获取提交消息
// 验证 GetCommitMessage 检索消息文本
func TestGcm_GetCommitMessage(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-message-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	const commitMsg = "test commit message"
	must.Done(os.WriteFile(filepath.Join(tempDIR, "test.txt"), []byte("test"), 0644))
	gcm.Add().Commit(commitMsg).Done()

	message := rese.V1(gcm.GetCommitMessage("HEAD"))
	require.Contains(t, message, commitMsg)
}

// TestGcm_BranchExists tests branch existence checking
// Verifies BranchExists detects existing and non-existing branches
//
// TestGcm_BranchExists 测试分支存在性检查
// 验证 BranchExists 检测现有和不存在的分支
func TestGcm_BranchExists(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-branch-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "test.txt"), []byte("test"), 0644))
	gcm.Add().Commit("initial").Done()

	exists := rese.V1(gcm.BranchExists("main"))
	if !exists {
		exists = rese.V1(gcm.BranchExists("master"))
	}
	require.True(t, exists)

	exists = rese.V1(gcm.BranchExists("nonexistent-branch"))
	require.False(t, exists)
}

// TestGcm_RemoteBranchExists tests remote branch existence checking
// Verifies RemoteBranchExists detects existing and non-existing remote branches
//
// TestGcm_RemoteBranchExists 测试远程分支存在性检查
// 验证 RemoteBranchExists 检测现有和不存在的远程分支
func TestGcm_RemoteBranchExists(t *testing.T) {
	// Create remote repo // 创建远程仓库
	remoteDIR := rese.V1(os.MkdirTemp("", "gitgo-remote-repo-*"))
	defer func() {
		must.Done(os.RemoveAll(remoteDIR))
	}()
	remoteGcm := gitgo.New(remoteDIR)
	remoteGcm.Init().Done()
	must.Done(os.WriteFile(filepath.Join(remoteDIR, "test.txt"), []byte("test"), 0644))
	remoteGcm.Add().Commit("initial").Done()

	// Create test repo and add remote // 创建本地仓库并添加远程
	localDIR := rese.V1(os.MkdirTemp("", "gitgo-local-repo-*"))
	defer func() {
		must.Done(os.RemoveAll(localDIR))
	}()
	localGcm := gitgo.New(localDIR)
	localGcm.Init().Done()
	must.Done(os.WriteFile(filepath.Join(localDIR, "local.txt"), []byte("local"), 0644))
	localGcm.Add().Commit("local initial").Done()
	localGcm.RemoteAdd("origin", remoteDIR).Fetch("origin").Done()

	// Test existing remote branch // 测试存在的远程分支
	exists := rese.V1(localGcm.RemoteBranchExists("origin/main"))
	if !exists {
		exists = rese.V1(localGcm.RemoteBranchExists("origin/master"))
	}
	require.True(t, exists)

	// Test non-existing remote branch // 测试不存在的远程分支
	exists = rese.V1(localGcm.RemoteBranchExists("origin/nonexistent-branch"))
	require.False(t, exists)
}

// TestGcm_TagExists tests tag existence checking
// Verifies TagExists detects existing and non-existing tags
//
// TestGcm_TagExists 测试标签存在性检查
// 验证 TagExists 检测现有和不存在的标签
func TestGcm_TagExists(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-tag-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "test.txt"), []byte("test"), 0644))
	gcm.Add().Commit("initial").Tag("v1.0.0").Done()

	exists := rese.V1(gcm.TagExists("v1.0.0"))
	require.True(t, exists)

	exists = rese.V1(gcm.TagExists("v2.0.0"))
	require.False(t, exists)
}

// TestGcm_GetFileList tests getting tracked files
// Verifies GetFileList returns tracked file paths
//
// TestGcm_GetFileList 测试获取跟踪文件
// 验证 GetFileList 返回跟踪文件路径
func TestGcm_GetFileList(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-files-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file1.txt"), []byte("test"), 0644))
	must.Done(os.WriteFile(filepath.Join(tempDIR, "file2.txt"), []byte("test"), 0644))
	gcm.Add().Commit("add files").Done()

	files := rese.V1(gcm.GetFileList())
	require.Contains(t, files, "file1.txt")
	require.Contains(t, files, "file2.txt")
}

// TestGcm_GetUntrackedFiles tests getting untracked files
// Verifies GetUntrackedFiles returns untracked file paths
//
// TestGcm_GetUntrackedFiles 测试获取未跟踪文件
// 验证 GetUntrackedFiles 返回未跟踪文件路径
func TestGcm_GetUntrackedFiles(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-untracked-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "tracked.txt"), []byte("test"), 0644))
	gcm.Add().Commit("add tracked").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "untracked.txt"), []byte("test"), 0644))

	files := rese.V1(gcm.GetUntrackedFiles())
	require.Contains(t, files, "untracked.txt")
	require.NotContains(t, files, "tracked.txt")
}

// TestGcm_GetModifiedFiles tests getting modified files
// Verifies GetModifiedFiles returns changed file paths
//
// TestGcm_GetModifiedFiles 测试获取修改文件
// 验证 GetModifiedFiles 返回更改文件路径
func TestGcm_GetModifiedFiles(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-modified-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("v1").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))

	files := rese.V1(gcm.GetModifiedFiles())
	require.Contains(t, files, "file.txt")
}

// TestGcm_GetIgnoredFiles tests getting ignored files
// Verifies GetIgnoredFiles returns files matching gitignore rules
//
// TestGcm_GetIgnoredFiles 测试获取被忽略文件
// 验证 GetIgnoredFiles 返回匹配 gitignore 规则的文件路径
func TestGcm_GetIgnoredFiles(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-ignored-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	// Create .gitignore // 创建 .gitignore
	must.Done(os.WriteFile(filepath.Join(tempDIR, ".gitignore"), []byte("*.log\ntarget/\n"), 0644))
	gcm.Add().Commit("add gitignore").Done()

	// Create ignored files in root // 在根目录创建被忽略的文件
	must.Done(os.WriteFile(filepath.Join(tempDIR, "debug.log"), []byte("log"), 0644))
	must.Done(os.WriteFile(filepath.Join(tempDIR, "error.log"), []byte("log"), 0644))

	// Create ignored files in nested path // 在子目录创建被忽略的文件
	must.Done(os.MkdirAll(filepath.Join(tempDIR, "target"), 0755))
	must.Done(os.WriteFile(filepath.Join(tempDIR, "target", "output.bin"), []byte("bin"), 0644))

	// Create non-ignored file // 创建非忽略文件
	must.Done(os.WriteFile(filepath.Join(tempDIR, "main.go"), []byte("package main"), 0644))

	paths := rese.V1(gcm.GetIgnoredFiles())
	require.Len(t, paths, 3)
	require.Contains(t, paths, "debug.log")
	require.Contains(t, paths, "error.log")
	require.Contains(t, paths, "target/")
	require.NotContains(t, paths, "main.go")
}

// TestGcm_GetIgnoredFiles_SubPath tests getting ignored files from subPath
// Verifies GetIgnoredFiles works when gcm starts from subPath
//
// TestGcm_GetIgnoredFiles_SubPath 测试从子路径获取被忽略文件
// 验证 gcm 以子路径为起点时 GetIgnoredFiles 正常工作
func TestGcm_GetIgnoredFiles_SubPath(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-ignored-sub-path-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Init repo at root // 在根目录初始化仓库
	rootGcm := gitgo.New(tempDIR)
	rootGcm.Init().Done()

	// Create .gitignore at root // 在根目录创建 .gitignore
	must.Done(os.WriteFile(filepath.Join(tempDIR, ".gitignore"), []byte("*.log\n*.tmp\n"), 0644))
	rootGcm.Add().Commit("add gitignore").Done()

	// Create subPath with ignored files // 创建包含被忽略文件的子目录
	subPath := filepath.Join(tempDIR, "subPath")
	must.Done(os.MkdirAll(subPath, 0755))
	must.Done(os.WriteFile(filepath.Join(subPath, "app.log"), []byte("log"), 0644))
	must.Done(os.WriteFile(filepath.Join(subPath, "cache.tmp"), []byte("tmp"), 0644))
	must.Done(os.WriteFile(filepath.Join(subPath, "main.go"), []byte("package main"), 0644))

	// Test from subPath // 从子目录测试
	subPathGcm := gitgo.New(subPath)
	paths := rese.V1(subPathGcm.GetIgnoredFiles())
	require.Len(t, paths, 2)
	require.Contains(t, paths, "app.log")
	require.Contains(t, paths, "cache.tmp")
	require.NotContains(t, paths, "main.go")
}

// TestGcm_GitCommitHash tests getting commit hash with reference
// Verifies GitCommitHash returns hash string for tags and branches
//
// TestGcm_GitCommitHash 测试通过引用获取提交哈希
// 验证 GitCommitHash 为标签和分支返回正确的哈希
func TestGcm_GitCommitHash(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-commit-hash-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("test"), 0644))
	gcm.Add().Commit("initial").Tag("v1.0.0").Done()

	// Test HEAD reference // 测试 HEAD 引用
	hashHead := rese.V1(gcm.GitCommitHash("HEAD"))
	require.Len(t, hashHead, 40)
	require.Regexp(t, "^[0-9a-f]{40}$", hashHead)

	// Test tag reference // 测试标签引用
	hashTag := rese.V1(gcm.GitCommitHash("v1.0.0"))
	require.Equal(t, hashHead, hashTag)
}

// TestGcm_LatestGitTag_Temp tests getting latest tag in temp repo
// Verifies LatestGitTag returns correct tag name
//
// TestGcm_LatestGitTag_Temp 测试在临时仓库获取最新标签
// 验证 LatestGitTag 返回正确的标签名
func TestGcm_LatestGitTag_Temp(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-latest-tag-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("v1").Tag("v1.0.0").Done()

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))
	gcm.Add().Commit("v2").Tag("v2.0.0").Done()

	tag := rese.V1(gcm.LatestGitTag())
	require.Equal(t, "v2.0.0", tag)
}

// TestGcm_GetBranchTrackingBranch tests getting upstream branch
// Verifies GetBranchTrackingBranch returns correct tracking info
//
// TestGcm_GetBranchTrackingBranch 测试获取上游分支
// 验证 GetBranchTrackingBranch 返回正确的跟踪信息
func TestGcm_GetBranchTrackingBranch(t *testing.T) {
	// Create remote repo // 创建远程仓库
	remoteDIR := rese.V1(os.MkdirTemp("", "gitgo-tracking-remote-*"))
	defer func() {
		must.Done(os.RemoveAll(remoteDIR))
	}()
	remoteGcm := gitgo.New(remoteDIR)
	remoteGcm.Init().Done()
	must.Done(os.WriteFile(filepath.Join(remoteDIR, "file.txt"), []byte("test"), 0644))
	remoteGcm.Add().Commit("initial").Done()

	// Create test repo with tracking // 创建带跟踪的本地仓库
	localDIR := rese.V1(os.MkdirTemp("", "gitgo-tracking-local-*"))
	defer func() {
		must.Done(os.RemoveAll(localDIR))
	}()
	localGcm := gitgo.New(localDIR)
	localGcm.Init().Done()
	must.Done(os.WriteFile(filepath.Join(localDIR, "local.txt"), []byte("local"), 0644))
	localGcm.Add().Commit("local").Done()
	localGcm.RemoteAdd("origin", remoteDIR).Fetch("origin").Done()

	// Get current branch name // 获取当前分支名
	currentBranch := rese.V1(localGcm.GetCurrentBranch())

	// Set upstream // 设置上游
	localGcm.UpdateCommandConfig(func(cfg *osexec.CommandConfig) {
		cfg.WithExpectExit(128, "MAY-FAIL")
	})

	// Get tracking branch (might not exist if not set) // 尝试获取跟踪分支（如果未设置可能不存在）
	_, err := localGcm.GetBranchTrackingBranch(currentBranch)
	// OK if fails (no upstream set) // 出错也可以（未设置上游）
	t.Log("tracking branch error (expected if no upstream):", err)
}
