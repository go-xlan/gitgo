package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
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

func TestGcm_LatestGitTagWithPrefixPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tag, err := gcm.LatestGitTagWithPrefixPath("v")
	require.NoError(t, err)
	t.Log(tag)
	require.NotEmpty(t, tag)
}

func TestGcm_LatestGitTagWithPrefixPath_Compare(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	tagA := rese.C1(gcm.LatestGitTag())
	t.Log(tagA)
	tagB := rese.C1(gcm.LatestGitTagWithPrefixPath("v"))
	t.Log(tagB)
	//这里只打印，不断言相等，因为确实存在不相等的场景
	t.Log("tag-A:", tagA, "tag-B:", tagB)
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
