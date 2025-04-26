package example1

import (
	"path/filepath"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
)

func TestGcm_GetTopPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	root, err := gcm.GetTopPath()
	require.NoError(t, err)
	t.Log(root)
	require.NotEmpty(t, root)

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

	root, err := gcm.GetTopPath()
	require.NoError(t, err)
	t.Log(root)
	require.NotEmpty(t, root)

	osmustexist.MustRoot(root)
	require.Equal(t, path, filepath.Join(root, ".git"))
}

func TestGcm_GetSubPathToRoot(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPathToRoot()
	require.NoError(t, err)
	t.Log(path)
	require.NotEmpty(t, path)
	require.Equal(t, "../../../", path)

	root, err := gcm.GetTopPath()
	require.NoError(t, err)
	t.Log(root)
	require.NotEmpty(t, root)

	osmustexist.MustRoot(root)
	require.Equal(t, root, runpath.PARENT.Join(path))
}

func TestGcm_GetSubPath(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()
	path, err := gcm.GetSubPath()
	require.NoError(t, err)
	t.Log(path)
	require.NotEmpty(t, path)
	require.Equal(t, "internal/examples/example1/", path)

	root, err := gcm.GetTopPath()
	require.NoError(t, err)
	t.Log(root)
	require.NotEmpty(t, root)

	osmustexist.MustRoot(root)
	require.Equal(t, runpath.PARENT.Path(), filepath.Join(root, path))
}
