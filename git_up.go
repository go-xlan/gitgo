package gogitxexec

import "github.com/pkg/errors"

func (G *GitCmx) CheckStagingChanges() (bool, error) {
	if data, err := G.Cmx.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		if len(data) != 0 {
			return false, errors.New(string(data))
		}
		return true, nil
	}
	return false, nil
}
