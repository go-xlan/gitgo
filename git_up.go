package gogitxexec

import "github.com/pkg/errors"

func (G *Gcm) HasStagingChanges() (bool, error) {
	if data, err := G.Cmc.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(data) != 0 {
			return false, errors.New(string(data))
		}
		return true, nil
	}
	return false, nil
}

func (G *Gcm) CheckStagedChanges() *Gcm {
	if output, err := G.Cmc.Exec("git", "diff", "--cached", "--quiet"); err == nil && len(output) == 0 {
		return newWa(G.Cmc, []byte{}, errors.New("no-staged-changes"), G.DBG)
	}
	return G
}
