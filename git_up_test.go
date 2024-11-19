package gogitxexec_test

import (
	"testing"

	"github.com/go-xlan/gogitxexec"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

func TestGcm_CheckStagingChanges(t *testing.T) {
	gcm := gogitxexec.New(runpath.PARENT.Path())

	changes, err := gcm.CheckStagingChanges()
	require.NoError(t, err)
	t.Log(changes)
}
