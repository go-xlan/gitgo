package gitgo

import "github.com/pkg/errors"

func (G *Gcm) HasStagingChanges() (bool, error) {
	if output, err := G.execConfig.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(output) != 0 {
			return false, errors.New(string(output))
		}
		return true, nil
	}
	return false, nil
}

func (G *Gcm) CheckStagedChanges() *Gcm {
	if output, err := G.execConfig.Exec("git", "diff", "--cached", "--quiet"); err == nil {
		if len(output) != 0 {
			return newWaGcm(G.execConfig, output, err, G.debugMode)
		}
		return newWaGcm(G.execConfig, []byte{}, errors.New("no-staged-changes"), G.debugMode)
	}
	return G
}
