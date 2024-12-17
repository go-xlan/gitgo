package gogitosgcm_test

import (
	"testing"

	"github.com/go-xlan/gogitosgcm"
	"github.com/yyle88/runpath"
)

func TestGcm_Status(t *testing.T) {
	gcm := gogitosgcm.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		MustDone()
}

func TestGcm_Submit(t *testing.T) {
	gcm := gogitosgcm.New(runpath.PARENT.Path())

	gcm.WithDebug().
		Status().
		Add().
		WhenExec(func(gcm *gogitosgcm.Gcm) (bool, error) {
			return gcm.HasStagingChanges()
		}, func(gcm *gogitosgcm.Gcm) *gogitosgcm.Gcm {
			return gcm.Commit("提交代码").Push()
		}).
		MustDone()
}
