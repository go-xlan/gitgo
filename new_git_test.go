package gogitxexec_test

import (
	"github.com/go-xlan/gogitxexec"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"testing"
)

func TestGitCmx_UpdateCmx(t *testing.T) {
	gcx := gogitxexec.New(runpath.PARENT.Path())

	gcx.WithDebug().
		UpdateCmx(func(cmx *osexec.CMX) {
			cmx.WithShellType("bash").WithShellFlag("-c")
		}).
		Add().
		ShowDebug().
		MustDone().
		Status().
		ShowDebug().
		MustDone()
}
