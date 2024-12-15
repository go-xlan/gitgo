package gogitxexec_test

import (
	"testing"

	"github.com/go-xlan/gogitxexec"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

func TestGcm_HasStagingChanges(t *testing.T) {
	gcm := gogitxexec.New(runpath.PARENT.Path())

	changes, err := gcm.HasStagingChanges()
	require.NoError(t, err)
	t.Log(changes)
}

func TestGcm_CheckStagedChanges(t *testing.T) {
	gcm := gogitxexec.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		Add().
		CheckStagedChanges().
		ShowDebugMessage().
		Status().
		ShowDebugMessage()
}
