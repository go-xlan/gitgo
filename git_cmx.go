package gogitxexec

func (G *GitCmx) Status() *GitCmx {
	return G.do("git", "status")
}

func (G *GitCmx) Add() *GitCmx {
	return G.do("git", "add", ".")
}

func (G *GitCmx) Commit(m string) *GitCmx {
	//当没有待提交文件时，这里也会报错，目前暂无解决方案。
	return G.do("git", "commit", "-m", m)
}

func (G *GitCmx) Pull() *GitCmx {
	return G.do("git", "pull")
}

func (G *GitCmx) Push() *GitCmx {
	return G.do("git", "push")
}

func (G *GitCmx) PushSetUpstreamOriginBranch(newBranchName string) *GitCmx {
	return G.do("git", "push", "--set-upstream", "origin", newBranchName)
}

func (G *GitCmx) Reset() *GitCmx {
	return G.do("git", "reset")
}

func (G *GitCmx) ResetHard() *GitCmx {
	return G.do("git", "reset", "--hard")
}

func (G *GitCmx) Checkout(nameBranch string) *GitCmx {
	return G.do("git", "checkout", nameBranch)
}

func (G *GitCmx) CheckoutNewBranch(nameBranch string) *GitCmx {
	return G.do("git", "checkout", "-b", nameBranch)
}

func (G *GitCmx) Init() *GitCmx {
	return G.do("git", "init")
}

func (G *GitCmx) Merge(featureBranchName string) *GitCmx {
	return G.do("git", "merge", featureBranchName)
}

func (G *GitCmx) MergeAbort() *GitCmx {
	return G.do("git", "merge", "--abort")
}

func (G *GitCmx) TagList() *GitCmx {
	return G.do("git", "tag", "--list")
}

func (G *GitCmx) Tags() *GitCmx {
	return G.do("git", "tag", "--list")
}

func (G *GitCmx) Tag(tag string) *GitCmx {
	return G.do("git", "tag", tag)
}

func (G *GitCmx) PushTags() *GitCmx {
	return G.do("git", "push", "--tags")
}
