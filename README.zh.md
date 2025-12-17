[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/gitgo/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/gitgo/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/gitgo)](https://pkg.go.dev/github.com/go-xlan/gitgo)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/gitgo/main.svg)](https://coveralls.io/github/go-xlan/gitgo?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/gitgo.svg)](https://github.com/go-xlan/gitgo/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/gitgo)](https://goreportcard.com/report/github.com/go-xlan/gitgo)

# gitgo

流式 Git 命令执行引擎，具有流畅的链式调用接口和全面的 Git 操作支持。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🔗 **流式链式接口**: 复杂 Git 工作流的方法链式调用，具有自动问题传播
⚡ **全面 Git 操作**: 完整覆盖 Git 命令，包括提交、推送、拉取和分支管理
🔍 **智能状态检测**: 智能检查暂存和未暂存更改、干净工作树和仓库状态
🎯 **问题处理**: 强健的问题传播，具有详细上下文和调试信息
📋 **仓库查询**: 高级仓库信息查询，包括分支、提交和状态信息

## 关联项目

- **[gogit](https://github.com/go-xlan/gogit)** - 增强型 Git 操作工具包，基于 go-git 实现，提供纯 Go 实现无需 CLI 依赖
- **[gitgo](https://github.com/go-xlan/gitgo)**（本项目）- 流式 Git 命令执行引擎，具有流畅的链式调用接口

## 安装

```bash
go get github.com/go-xlan/gitgo
```

## 使用方法

### 基础 Git 操作

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

⬆️ **Source:** [源码](internal/demos/demo1x/main.go)

### 仓库状态检测

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

⬆️ **Source:** [源码](internal/demos/demo2x/main.go)

### 标签和仓库信息

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

	latest, exists, err := gcm.GetLatestTag()
	must.Done(err)
	must.True(exists)
	zaplog.SUG.Info("latest tag:", latest)

	count := rese.V1(gcm.GetCommitCount())
	zaplog.SUG.Info("commit count:", count)
}
```

⬆️ **Source:** [源码](internal/demos/demo3x/main.go)

## API 参考

### 核心方法

- `New(path string) *Gcm` - 创建新的 Git 命令引擎
- `NewGcm(path, execConfig) *Gcm` - 使用自定义设置创建

### Git 操作

- `Status() *Gcm` - 显示工作树状态
- `Add() *Gcm` - 暂存更改
- `Commit(message) *Gcm` - 创建带消息的提交
- `Push() *Gcm` - 推送到远程仓库
- `Pull() *Gcm` - 从远程仓库获取并合并

### 分支管理

- `CheckoutNewBranch(name) *Gcm` - 创建并切换到新分支
- `Checkout(name) *Gcm` - 切换到现有分支
- `GetCurrentBranch() (string, error)` - 获取分支名称
- `ListBranches() ([]string, error)` - 获取分支列表

### 仓库状态

- `HasStagedChanges() (bool, error)` - 检查暂存更改是否存在
- `HasUnstagedChanges() (bool, error)` - 检查未暂存更改是否存在
- `HasChanges() (bool, error)` - 检查更改是否存在
- `GetCommitCount() (int, error)` - 获取提交数量
- `GetCommitHash(ref) (string, error)` - 使用引用获取提交哈希
- `GetRemoteURL(remote) (string, error)` - 获取远程仓库 URL
- `GetIgnoredFiles() ([]string, error)` - 获取 gitignore 忽略的文件
- `ConfigGet(key) (string, error)` - 获取 Git 配置值

### 标签操作

- `GetLatestTag() (string, bool, error)` - 获取最新标签名称并检查是否存在

### 问题处理

- `Result() ([]byte, error)` - 获取输出并检查问题
- `MustDone() *Gcm` - 当问题发生时触发 panic

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-xlan/gitgo.svg?variant=adaptive)](https://starchart.cc/go-xlan/gitgo)
