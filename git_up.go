package gogitxexec

func (G *GitCmx) CheckStagingChanges() bool {
	if _, err := G.Cmx.Exec("git", "diff-index", "--cached", "--quiet", "HEAD"); err != nil {
		return true
	} else {
		return false
	}
}
