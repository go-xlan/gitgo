package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

// TestGcm_UpdateCommandConfig tests command configuration customization with bash shell
// Verifies that custom shell type and flag settings work with Git operations
// Uses debug mode to show command execution details and output
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

// TestNewGcm tests Gcm creation with custom execution configuration
// Verifies that explicit ExecConfig parameter works with debug mode
// Validates status command execution with message display
func TestNewGcm(t *testing.T) {
	execConfig := osexec.NewExecConfig()

	gcm := gitgo.NewGcm(runpath.PARENT.Path(), execConfig)

	gcm.WithDebug().
		Status().
		ShowDebugMessage().
		MustDone()
}
