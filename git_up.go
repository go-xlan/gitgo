package gogitxexec

import "github.com/pkg/errors"

func (G *Gcm) CheckStagingChanges() (bool, error) {
	if data, err := G.Cmc.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(data) != 0 {
			return false, errors.New(string(data))
		}
		return true, nil
	}
	return false, nil
}
