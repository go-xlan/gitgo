package gogitxexec_test

import (
	"github.com/go-xlan/gogitxexec"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
	"testing"
)

func TestGitCmx_CheckStagingChanges(t *testing.T) {
	gcx := gogitxexec.New(runpath.PARENT.Path())

	changes, err := gcx.CheckStagingChanges()
	require.NoError(t, err)
	t.Log(changes)
}
