package gogitxexec_test

import (
	"testing"

	"github.com/go-xlan/gogitxexec"
	"github.com/yyle88/runpath"
)

func TestGcm_Status(t *testing.T) {
	gcm := gogitxexec.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		MustDone()
}

func TestGcm_Submit(t *testing.T) {
	gcm := gogitxexec.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		Add().
		WhenExec(func(gcm *gogitxexec.Gcm) (bool, error) {
			return gcm.HasStagingChanges()
		}, func(gcm *gogitxexec.Gcm) *gogitxexec.Gcm {
			return gcm.Commit("提交代码").Push()
		}).
		MustDone()
}
