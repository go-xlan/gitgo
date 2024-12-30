package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestGcm_UpdateCommandConfig(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebugMode(true).
		UpdateCommandConfig(func(cfg *osexec.CommandConfig) {
			cfg.WithShellType("bash").WithShellFlag("-c")
		}).
		Add().
		ShowDebugMessage().
		MustDone().
		Status().
		ShowDebugMessage().
		MustDone()
}
