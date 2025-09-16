package example2_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func TestBasicOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	ok, err := git.IsInsideWorkTree()
	require.NoError(t, err)
	require.True(t, ok)

	branch, err := git.GetCurrentBranch()
	require.NoError(t, err)
	require.Equal(t, "main", branch)
}

func TestBranchOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file1 := filepath.Join(tempDIR, "main.txt")
	must.Done(os.WriteFile(file1, []byte("main"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	_, err = git.CheckoutNewBranch("feature").Result()
	require.NoError(t, err)

	file2 := filepath.Join(tempDIR, "feature.txt")
	must.Done(os.WriteFile(file2, []byte("feature"), 0644))
	_, err = git.Add().Commit("add feature").Result()
	require.NoError(t, err)

	_, err = git.Checkout("main").Result()
	require.NoError(t, err)
	_, err = git.Merge("feature").Result()
	require.NoError(t, err)

	_, err = os.Stat(file2)
	require.NoError(t, err)
}

func TestLogOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	output, err := git.Status().Result()
	require.NoError(t, err)
	require.NotEmpty(t, output)
}
