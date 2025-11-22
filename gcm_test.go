package gitgo_test

import (
	"os"
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestGcm_Status tests basic Git status command execution
// Verifies working tree state shown with debug output
// Uses parent path as the Git space location
func TestGcm_Status(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		MustDone()
}

// TestGcm_Submit tests condition-based commit and push workflow
// Demonstrates WhenThen pattern to commit and push when changes exist
// Disabled on default (if false) to prevent unintended commits in tests
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

// TestGcm_RemoteOperations tests remote repo management operations
// Verifies Remote, RemoteAdd, RemoteRemove, and RemoteSet functions
//
// TestGcm_RemoteOperations 测试远程仓库管理操作
// 验证 Remote、RemoteAdd、RemoteRemove 和 RemoteSet 函数
func TestGcm_RemoteOperations(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-remote-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()

	gcm.RemoteAdd("origin", "https://github.com/example/test.git").Done()

	output := gcm.Remote().Output()
	require.Contains(t, string(output), "origin")
	require.Contains(t, string(output), "https://github.com/example/test.git")

	gcm.RemoteSet("origin", "git@github.com:example/test.git").Done()
	output = gcm.Remote().Output()
	require.Contains(t, string(output), "git@github.com")

	gcm.RemoteRemove("origin").Done()
	output = gcm.Remote().Output()
	require.NotContains(t, string(output), "origin")
}
