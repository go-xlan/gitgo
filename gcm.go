package gogitxexec

func (G *Gcm) Status() *Gcm {
	return G.do("git", "status")
}

func (G *Gcm) Add() *Gcm {
	return G.do("git", "add", ".")
}

func (G *Gcm) Commit(m string) *Gcm {
	//当没有待提交文件时，这里也会报错，目前暂无解决方案。
	return G.do("git", "commit", "-m", m)
}

func (G *Gcm) Pull() *Gcm {
	return G.do("git", "pull")
}

func (G *Gcm) Push() *Gcm {
	return G.do("git", "push")
}

func (G *Gcm) PushSetUpstreamOriginBranch(newBranchName string) *Gcm {
	return G.do("git", "push", "--set-upstream", "origin", newBranchName)
}

func (G *Gcm) Reset() *Gcm {
	return G.do("git", "reset")
}

func (G *Gcm) ResetHard() *Gcm {
	return G.do("git", "reset", "--hard")
}

func (G *Gcm) Checkout(nameBranch string) *Gcm {
	return G.do("git", "checkout", nameBranch)
}

func (G *Gcm) CheckoutNewBranch(nameBranch string) *Gcm {
	return G.do("git", "checkout", "-b", nameBranch)
}

func (G *Gcm) Init() *Gcm {
	return G.do("git", "init")
}

func (G *Gcm) Merge(featureBranchName string) *Gcm {
	return G.do("git", "merge", featureBranchName)
}

func (G *Gcm) MergeAbort() *Gcm {
	return G.do("git", "merge", "--abort")
}

func (G *Gcm) TagList() *Gcm {
	return G.do("git", "tag", "--list")
}

func (G *Gcm) Tags() *Gcm {
	return G.do("git", "tag", "--list")
}

func (G *Gcm) Tag(tag string) *Gcm {
	return G.do("git", "tag", tag)
}

func (G *Gcm) PushTags() *Gcm {
	return G.do("git", "push", "--tags")
}
