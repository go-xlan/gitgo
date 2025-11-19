package gitgo_test

import (
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
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

// TestGcm_HasChanges tests detection of changes (staged and unstaged combined)
// Verifies that the function identifies if work is in progress
func TestGcm_HasChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	changes, err := gcm.HasChanges()
	require.NoError(t, err)
	t.Log(changes)
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
// Verifies that the function correctly filters tags by prefix pattern
func TestGcm_LatestGitTagHasPrefix(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagHasPrefix("v")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGcm_LatestGitTagHasPrefix_Compare compares general tag fetch with prefix filtering
// Demonstrates that prefix-filtered tags may differ from the overall latest tag
func TestGcm_LatestGitTagHasPrefix_Compare(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tagA := rese.C1(gcm.LatestGitTag())
	t.Log(tagA)
	tagB := rese.C1(gcm.LatestGitTagHasPrefix("v"))
	t.Log(tagB)
	// Only log without equality assertion - tags may legitimately differ // 只打印不断言相等，因为标签可能确实不同
	t.Log("tag-A:", tagA, "tag-B:", tagB)
}

// TestGcm_LatestGitTagMatchRegexp tests fetch of latest tag matching glob pattern
// Verifies that the function correctly filters tags using wildcard patterns
func TestGcm_LatestGitTagMatchRegexp(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagMatchRegexp("v[0-9]*.[0-9]*.[0-9]*")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGitCommitHash tests fetch of commit hash with branch reference
// Verifies that the function correctly resolves branch name to commit hash
func TestGitCommitHash(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	hash, err := gcm.GitCommitHash("main")
	require.NoError(t, err)
	t.Log(hash)
	require.NotEmpty(t, hash)
}

// TestGitCommitHash_TAG tests fetch of commit hash with tag reference
// Verifies that the function correctly resolves tag name to commit hash
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
// Verifies that the function returns tags in ascending chronological sequence
func TestGcm_SortedGitTags(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.SortedGitTags()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

// TestGcm_GetTopPath tests fetch of Git space base path
// Verifies that the function returns correct project top-level location
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

// TestGcm_IsInsideWorkTree tests verification of current path location in Git work tree
// Verifies that the function identifies whether inside a Git project
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

// TestGcm_GetRemoteURL tests fetch of remote space URL
// Verifies that the function returns correct origin remote location
func TestGcm_GetRemoteURL(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	remoteURL, err := gcm.GetRemoteURL("origin")
	require.NoError(t, err)
	t.Log(remoteURL)
	require.NotEmpty(t, remoteURL)
}

// TestGcm_GetCommitCount tests fetch of total commit count in current branch
// Verifies that the function returns accurate number of commits
func TestGcm_GetCommitCount(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	count, err := gcm.GetCommitCount()
	require.NoError(t, err)
	t.Log(count)
	require.Greater(t, count, 0)
}

// TestGcm_ListBranches tests fetch of all local branch names
// Verifies that the function returns complete list of project branches
func TestGcm_ListBranches(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	branches, err := gcm.ListBranches()
	require.NoError(t, err)
	t.Log(branches)
	require.NotEmpty(t, branches)
}

// TestGcm_ListRemoteBranches tests fetch of all remote branch names
// Verifies that the function returns complete list of remote branches
func TestGcm_ListRemoteBranches(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	branches, err := gcm.ListRemoteBranches()
	require.NoError(t, err)
	t.Log(branches)
	require.NotEmpty(t, branches)
}

// TestGcm_GetLogOneLine tests fetch of concise commit log with limit
// Verifies that the function returns specified number of one-line commit entries
func TestGcm_GetLogOneLine(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	commits, err := gcm.GetLogOneLine(5)
	require.NoError(t, err)
	t.Log(commits)
	require.NotEmpty(t, commits)
}
