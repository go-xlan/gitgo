package gogitxexec_test

import (
	"github.com/go-xlan/gogitxexec"
	"github.com/yyle88/runpath"
	"testing"
)

func TestGitCmx_Status(t *testing.T) {
	gcx := gogitxexec.New(runpath.PARENT.Path())

	gcx.WithDebug().
		Status().
		MustDone()
}

func TestGitCmx_Commit(t *testing.T) {
	gcx := gogitxexec.New(runpath.PARENT.Path())

	gcx.WithDebug().
		Status().
		Add().
		Commit("提交代码").
		Push().
		MustDone()
}
