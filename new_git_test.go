package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestGcm_UpdateCmc(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		UpdateCmc(func(cmc *osexec.CommandConfig) {
			cmc.WithShellType("bash").WithShellFlag("-c")
		}).
		Add().
		ShowDebugMessage().
		MustDone().
		Status().
		ShowDebugMessage().
		MustDone()
}
