package gogitxexec_test

import (
	"github.com/go-xlan/gogitxexec"
	"github.com/yyle88/runpath"
	"testing"
)

func TestGitCmx_CheckStagingChanges(t *testing.T) {
	gcx := gogitxexec.New(runpath.PARENT.Path())

	changes := gcx.CheckStagingChanges()
	t.Log(changes)
}
