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
	zaplog.SUG.Info("git repo initialized")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "demo.txt"), []byte("hello"), 0644))
	zaplog.SUG.Info("created demo.txt")

	gcm.Add().Commit("demo commit").Done()
	zaplog.SUG.Info("committed changes")
}
