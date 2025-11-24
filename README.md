[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/gitgo/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/gitgo/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/gitgo)](https://pkg.go.dev/github.com/go-xlan/gitgo)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/gitgo/main.svg)](https://coveralls.io/github/go-xlan/gitgo?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://github.com/go-xlan/gitgo)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/gitgo.svg)](https://github.com/go-xlan/gitgo/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/gitgo)](https://goreportcard.com/report/github.com/go-xlan/gitgo)

# gitgo

Streamlined Git command execution engine with fluent chaining interface and comprehensive Git operations support.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🔗 **Fluent Chaining Interface**: Method chaining with complex Git workflows and automatic issue propagation
⚡ **Comprehensive Git Operations**: Complete Git commands coverage including commit, push, fetch, and branch management
🔍 **Smart State Detection**: Intelligent checking on staged and unstaged changes, clean working trees, and repo status
🎯 **Issue Handling**: Robust issue propagation with detailed context and debug information
📋 **Repo Querying**: Advanced repo introspection with branch, commit, and status information

## Related Projects

- **[gogit](https://github.com/go-xlan/gogit)** - Enhanced Git operations toolkit with go-git foundation, providing pure Go implementation without CLI dependencies
- **[gitgo](https://github.com/go-xlan/gitgo)** (this project) - Streamlined Git command execution engine with fluent chaining interface

## Installation

```bash
go get github.com/go-xlan/gitgo
```

## Usage

### Basic Git Operations

```go
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
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### Repo State Detection

```go
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
	zaplog.SUG.Info("initialized repo")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v1"), 0644))
	gcm.Add().Commit("add file").Done()
	zaplog.SUG.Info("committed v1")

	must.Done(os.WriteFile(filepath.Join(tempDIR, "file.txt"), []byte("v2"), 0644))
	zaplog.SUG.Info("modified file to v2")

	hasChanges := rese.V1(gcm.HasUnstagedChanges())
	zaplog.SUG.Info("has unstaged changes:", hasChanges)

	if hasChanges {
		gcm.Add().Commit("update file").Done()
		zaplog.SUG.Info("committed v2 changes")
	}
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

### Tags and Repo Information

```go
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
```

⬆️ **Source:** [Source](internal/demos/demo3x/main.go)

## API Reference

### Core Methods

- `New(path string) *Gcm` - Create new Git command engine
- `NewGcm(path, execConfig) *Gcm` - Create with custom settings

### Git Operations

- `Status() *Gcm` - Show working tree status
- `Add() *Gcm` - Stage changes
- `Commit(message) *Gcm` - Create commit with message
- `Push() *Gcm` - Push to remote repo
- `Pull() *Gcm` - Fetch and merge from remote repo

### Branch Management

- `CheckoutNewBranch(name) *Gcm` - Create and switch to new branch
- `Checkout(name) *Gcm` - Switch to existing branch
- `GetCurrentBranch() (string, error)` - Get the branch name
- `ListBranches() ([]string, error)` - Get branches as a list

### Repo State

- `HasStagingChanges() (bool, error)` - Check staged changes existence
- `HasUnstagedChanges() (bool, error)` - Check unstaged changes existence
- `HasChanges() (bool, error)` - Check changes existence
- `GetCommitCount() (int, error)` - Get commit count
- `GitCommitHash(ref) (string, error)` - Get commit hash with reference
- `GetRemoteURL(remote) (string, error)` - Get remote repo URL
- `GetIgnoredFiles() ([]string, error)` - Get files ignored in gitignore

### Tag Operations

- `LatestGitTag() (string, error)` - Get latest tag name (fails when no tags exist)
- `GetLatestTag() (string, bool, error)` - Get latest tag name with existence check

### Issue Handling

- `Result() ([]byte, error)` - Get output and check issues
- `MustDone() *Gcm` - Panic when issues happen

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-06 04:53:24.895249 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the concept
- 📖 **Documentation confusing?** Report it and we can enhance
- 🚀 **Need new features?** Share the use cases to help us understand needs
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package helped the workflow
- 💬 **Feedback?** We welcome input and comments

---

## 🔧 Development

When contributing new code, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI)
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`)
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`)
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Make sure tests pass (`go test ./...`) and follow Go code conventions
7. **Documentation**: Update documentation to support changes and write commit messages with substance
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) with backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`)
11. **PR**: Open a PR on GitHub (on the GitHub webpage) with detailed description

Make sure tests pass and include documentation updates as needed.

---

## 🌟 Support

Welcome to contribute to this project through submitting PRs and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** when this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/gitgo.svg?variant=adaptive)](https://starchart.cc/go-xlan/gitgo)