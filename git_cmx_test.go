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
		When(func(cmx *gogitxexec.GitCmx) bool {
			return cmx.CheckStagingChanges()
		}, func(cmx *gogitxexec.GitCmx) *gogitxexec.GitCmx {
			return cmx.Commit("提交代码").Push()
		}).
		MustDone()
}
