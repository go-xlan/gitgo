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
	zaplog.SUG.Info("repo setup complete")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "app.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("v1").Tag("v1.0.0").Done()
	zaplog.SUG.Info("tagged v1.0.0")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "app.txt"), []byte("v2"), 0644))
	gcm.Add().Commit("v2").Tag("v1.1.0").Done()
	zaplog.SUG.Info("tagged v1.1.0")

	latest := rese.V1(gcm.LatestGitTag())
	zaplog.SUG.Info("latest tag:", latest)

	count := rese.V1(gcm.GetCommitCount())
	zaplog.SUG.Info("commit count:", count)
}
