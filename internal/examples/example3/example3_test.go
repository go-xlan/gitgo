package example3_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func TestBranchOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file1 := filepath.Join(tempDIR, "main.go")
	must.Done(os.WriteFile(file1, []byte("package main\n"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	_, err = git.CheckoutNewBranch("feature").Result()
	require.NoError(t, err)

	file2 := filepath.Join(tempDIR, "feature.go")
	must.Done(os.WriteFile(file2, []byte("package main\n"), 0644))
	_, err = git.Add().Commit("add feature").Result()
	require.NoError(t, err)

	branches, err := git.ListBranches()
	require.NoError(t, err)
	require.Len(t, branches, 2)

	_, err = git.Checkout("main").Result()
	require.NoError(t, err)
	_, err = git.Merge("feature").Result()
	require.NoError(t, err)

	_, err = os.Stat(file2)
	require.NoError(t, err)
}

func TestResetOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("original"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	must.Done(os.WriteFile(file, []byte("modified"), 0644))
	_, err = git.Add().Result()
	require.NoError(t, err)

	staged, err := git.HasStagingChanges()
	require.NoError(t, err)
	require.True(t, staged)

	_, err = git.Reset().Result()
	require.NoError(t, err)

	staged, err = git.HasStagingChanges()
	require.NoError(t, err)
	require.False(t, staged)

	_, err = git.ResetHard().Result()
	require.NoError(t, err)

	changed, err := git.HasChanges()
	require.NoError(t, err)
	require.False(t, changed)
}

func TestDebugOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	git := gitgo.New(tempDIR).WithDebug()
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("test").Result()
	require.NoError(t, err)

	_, err = git.Status().Result()
	require.NoError(t, err)
}
