package gitgo

import "github.com/pkg/errors"

func (G *Gcm) HasStagingChanges() (bool, error) {
	if output, err := G.cmdConfig.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(output) != 0 {
			return false, errors.New(string(output))
		}
		return true, nil
	}
	return false, nil
}

func (G *Gcm) CheckStagedChanges() *Gcm {
	if output, err := G.cmdConfig.Exec("git", "diff", "--cached", "--quiet"); err == nil {
		if len(output) != 0 {
			return newWaGcm(G.cmdConfig, output, err, G.debugMode)
		}
		return newWaGcm(G.cmdConfig, []byte{}, errors.New("no-staged-changes"), G.debugMode)
	}
	return G
}
