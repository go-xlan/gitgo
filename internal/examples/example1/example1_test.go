package example1_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func TestBasicGitInfo(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	path, err := git.GetTopPath()
	require.NoError(t, err)
	require.NotEmpty(t, path)

	ok, err := git.IsInsideWorkTree()
	require.NoError(t, err)
	require.True(t, ok)
}

func TestBranchOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	branches, err := git.ListBranches()
	require.NoError(t, err)
	require.Len(t, branches, 1)

	branch, err := git.GetCurrentBranch()
	require.NoError(t, err)
	require.Equal(t, "main", branch)
}

func TestLogOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("add file").Result()
	require.NoError(t, err)

	logs, err := git.GetLogOneLine(1)
	require.NoError(t, err)
	require.Len(t, logs, 1)
}
