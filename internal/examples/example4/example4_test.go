package example4_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func TestTagOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	tag := "v1.0.0"
	_, err = git.Tag(tag).Result()
	require.NoError(t, err)

	output, err := git.Tags().Result()
	require.NoError(t, err)
	require.Contains(t, string(output), tag)

	latest, exists, err := git.GetLatestTag()
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, tag, latest)
}

func TestTagWithPrefix(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file1 := filepath.Join(tempDIR, "auth.go")
	must.Done(os.WriteFile(file1, []byte("package auth\n"), 0644))
	_, err = git.Add().Commit("auth").Result()
	require.NoError(t, err)
	_, err = git.Tag("auth/v1.0.0").Result()
	require.NoError(t, err)

	file2 := filepath.Join(tempDIR, "api.go")
	must.Done(os.WriteFile(file2, []byte("package api\n"), 0644))
	_, err = git.Add().Commit("api").Result()
	require.NoError(t, err)
	_, err = git.Tag("api/v1.0.0").Result()
	require.NoError(t, err)

	auth, err := git.GetLatestTagHasPrefix("auth/")
	require.NoError(t, err)
	require.Equal(t, "auth/v1.0.0", auth)

	api, err := git.GetLatestTagHasPrefix("api/")
	require.NoError(t, err)
	require.Equal(t, "api/v1.0.0", api)
}

func TestCommitHash(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	git := gitgo.New(tempDIR)
	_, err := git.Init().Result()
	require.NoError(t, err)

	file := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(file, []byte("test"), 0644))
	_, err = git.Add().Commit("init").Result()
	require.NoError(t, err)

	head, err := git.GetCommitHash("HEAD")
	require.NoError(t, err)
	require.NotEmpty(t, head)
	require.Len(t, head, 40)

	tag := "v1.0.0"
	_, err = git.Tag(tag).Result()
	require.NoError(t, err)

	hash, err := git.GetCommitHash(tag)
	require.NoError(t, err)
	require.Equal(t, head, hash)
}
