package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/runpath"
)

// TestSetDebugMode tests package-level debug mode activation and status command output
// Verifies that debug logging displays Git command execution details
func TestSetDebugMode(t *testing.T) {
	gitgo.SetDebugMode(true)

	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.Status()
}
