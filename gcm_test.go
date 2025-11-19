package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/runpath"
)

// TestGcm_Status tests basic Git status command execution
// Verifies working tree state display with debug output
// Uses parent path as the Git space location
func TestGcm_Status(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		MustDone()
}

// TestGcm_Submit tests condition-based commit and push workflow
// Demonstrates WhenThen pattern to commit and push when changes exist
// Disabled by default (if false) to prevent unintended commits in tests
func TestGcm_Submit(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	if false { // Disabled to prevent unintended commits // 禁用以防止意外提交
		gcm.WithDebug().
			Status().
			Add().
			WhenThen(func(gcm *gitgo.Gcm) (bool, error) {
				return gcm.HasStagingChanges()
			}, func(gcm *gitgo.Gcm) *gitgo.Gcm {
				return gcm.Commit("提交代码").Push()
			}).
			MustDone()
	}
}
