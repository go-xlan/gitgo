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

func TestGcm_HasStagingChanges(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	changes, err := gcm.HasStagingChanges()
	require.NoError(t, err)
	t.Log(changes)
}

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

func TestLatestGitTag(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTag()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

func TestGcm_LatestGitTagHasPrefix(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagHasPrefix("v")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

func TestGcm_LatestGitTagHasPrefix_Compare(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tagA := rese.C1(gcm.LatestGitTag())
	t.Log(tagA)
	tagB := rese.C1(gcm.LatestGitTagHasPrefix("v"))
	t.Log(tagB)
	//这里只打印，不断言相等，因为确实存在不相等的场景
	t.Log("tag-A:", tagA, "tag-B:", tagB)
}

func TestGcm_LatestGitTagMatchRegexp(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagMatchRegexp("v[0-9]*.[0-9]*.[0-9]*")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

func TestGitCommitHash(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	hash, err := gcm.GitCommitHash("main")
	require.NoError(t, err)
	t.Log(hash)
	require.NotEmpty(t, hash)
}

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

func TestGcm_SortedGitTags(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.SortedGitTags()
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

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

func TestGcm_GetGitDirAbsPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetGitDirAbsPath()
	require.NoError(t, err)
	t.Log(path)
	require.NotEmpty(t, path)

	osmustexist.MustRoot(path)
	require.Equal(t, runpath.PARENT.Join(".git"), path)
}

func TestGcm_GetSubPathToRoot(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPathToRoot()
	require.NoError(t, err)
	t.Log(path)
	require.Empty(t, path)
	t.Log("sub path is the project root")
	osmustexist.MustRoot(runpath.PARENT.Join(".git"))
}

func TestGcm_GetSubPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPath()
	require.NoError(t, err)
	t.Log(path)
	require.Empty(t, path)
	t.Log("sub path is the project root")
	osmustexist.MustRoot(runpath.PARENT.Join(".git"))
}
