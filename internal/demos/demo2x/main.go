package main

import (
	"os"
	"path/filepath"

	"github.com/go-xlan/gitgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	tempDIR := rese.V1(os.MkdirTemp("", "gitgo-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	zaplog.SUG.Debug("working in:", tempDIR)

	gcm := gitgo.New(tempDIR)
	gcm.Init().Done()
	zaplog.SUG.Info("initialized repo")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("add file").Done()
	zaplog.SUG.Info("committed v1")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))
	zaplog.SUG.Info("modified file to v2")

	hasChanges := rese.V1(gcm.HasUnstagedChanges())
	zaplog.SUG.Info("has unstaged changes:", hasChanges)

	if hasChanges {
		gcm.Add().Commit("update file").Done()
		zaplog.SUG.Info("committed v2 changes")
	}
}
