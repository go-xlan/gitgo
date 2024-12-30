package gitgo_test

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/runpath"
)

func TestSetDebugMode(t *testing.T) {
	gitgo.SetDebugMode(true)

	gcm := gitgo.New(runpath.PARENT.Path())

	gcm.Status()
}
