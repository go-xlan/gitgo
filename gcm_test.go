package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/runpath"
)

func TestGcm_Status(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		MustDone()
}

func TestGcm_Submit(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		Add().
		WhenExec(func(gcm *gitgo.Gcm) (bool, error) {
			return gcm.HasStagingChanges()
		}, func(gcm *gitgo.Gcm) *gitgo.Gcm {
			return gcm.Commit("提交代码").Push()
		}).
		MustDone()
}
