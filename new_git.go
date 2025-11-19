// Package gitgo: Streamlined Git command execution engine with fluent chaining interface
// Provides comprehensive Git operations through os/exec with smart handling and debug support
// Features method chaining to handle complex Git workflows with integrated logging capabilities
// Supports main Git operations including commit, push, branch management, and status queries
//
// gitgo: 流式 Git 命令执行引擎，带有流畅的链式调用接口
// 通过 os/exec 提供全面的 Git 操作，具有智能处理和调试支持
// 具有方法链式调用功能以处理复杂 Git 工作流，集成日志记录功能
// 支持主要 Git 操作，包括提交、推送、分支管理和状态查询
package gitgo

import (
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

// Gcm represents a Git Command Engine with chaining support and integrated handling
// Maintains execution state, output capture, and debug information between method invocations
// Supports fluent interface to manage complex Git workflows with automatic propagation
//
// State Management:
// - execConfig: Execution configuration with working path context
// - output: Command output bytes from recent operations (both success and failures)
// - errorOnce: First error encountered in chain (becomes nil when operations succeed)
// - debugMode: Activates detailed debug logging with colored console output
//
// Gcm 代表 Git 命令引擎，支持链式调用和集成处理
// 在方法调用间维护执行状态、输出捕获和调试信息
// 支持流畅接口以管理复杂 Git 工作流，具有自动传播功能
//
// 状态管理：
// - execConfig: 执行配置和工作路径上下文
// - output: 来自最近操作的命令输出字节（成功和失败）
// - errorOnce: 链中遇到的第一个错误（操作成功时为 nil）
// - debugMode: 启用带有彩色控制台输出的详细调试日志
type Gcm struct {
	execConfig *osexec.ExecConfig // Execution configuration with path context // 执行配置和路径上下文
	output     []byte             // Last command output bytes // 最后命令的输出字节
	errorOnce  error              // First error in the chain // 链中的第一个错误
	debugMode  bool               // Debug logging flag // 调试日志标志
}

// New creates a new Gcm instance with default configuration at the specified path
// Initializes Git command engine with standard settings and execution context
// Returns Gcm instance that is configured and prepared to chain Git operations
//
// New 在指定路径创建具有默认配置的新 Gcm 实例
// 使用标准设置和执行上下文初始化 Git 命令引擎
// 返回已配置和准备好的 Gcm 实例以进行链式 Git 操作
func New(path string) *Gcm {
	return newOkGcm(osexec.NewCommandConfig().WithPath(path).WithDebugMode(osexec.NewDebugMode(debugModeOpen)), make([]byte, 0), debugModeOpen)
}

// NewGcm creates a new Gcm instance with custom execution configuration
// Allows advanced configuration of command execution environment and settings
// Provides adaptation when specialized Git operation needs arise
//
// NewGcm 使用自定义执行配置创建新的 Gcm 实例
// 允许高级配置命令执行环境和行为
// 在专门的 Git 操作需求出现时提供适配
func NewGcm(path string, execConfig *osexec.ExecConfig) *Gcm {
	return newOkGcm(execConfig.NewConfig().WithPath(path).WithDebugMode(osexec.NewDebugMode(debugModeOpen)), make([]byte, 0), debugModeOpen)
}

// newOkGcm creates success-state Gcm instance with green success logging in debug mode
// This function builds Gcm with no errors to continue command chains
// Shows green-tinted success messages with command output details when debugging
//
// newOkGcm 在调试模式下创建带有绿色成功日志的成功状态 Gcm 实例
// 该函数构建无错误的 Gcm 以继续命令链
// 调试时显示带有命令输出详情的绿色成功消息
func newOkGcm(execConfig *osexec.ExecConfig, output []byte, debugMode bool) *Gcm {
	if debugMode {
		if len(output) > 0 {
			zaplog.ZAPS.Skip3.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(output))+"\n", "-")
		} else {
			zaplog.ZAPS.Skip3.SUG.Debugln("done", "\n", "-")
		}
	}
	return &Gcm{
		execConfig: execConfig,
		output:     output,
		errorOnce:  nil,
		debugMode:  debugMode,
	}
}

// newWaGcm creates a failed-state Gcm instance with red logging in debug mode
// This function builds Gcm with captured errors to stop command chains
// Shows red-tinted messages with error details when debugging
//
// newWaGcm 在调试模式下创建带有红色日志的失败状态 Gcm 实例
// 该函数构建具有捕获错误的 Gcm 以停止命令链
// 调试时显示带有错误详情的红色消息
func newWaGcm(execConfig *osexec.ExecConfig, output []byte, errorOnce error, debugMode bool) *Gcm {
	if debugMode {
		if len(output) > 0 {
			zaplog.ZAPS.Skip3.SUG.Errorln("wrong", eroticgo.RED.Sprint(errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(output))+"\n", "-")
		} else {
			zaplog.ZAPS.Skip3.SUG.Errorln("wrong", eroticgo.RED.Sprint(errorOnce))
		}
	}
	return &Gcm{
		execConfig: execConfig,
		output:     output,
		errorOnce:  errorOnce,
		debugMode:  debugMode,
	}
}

// Result returns the output and error from the command chains
// Provides access to command output with errors that happened in execution
// Use case: extract results when completing Git operation chains
//
// Result 返回命令链的输出和错误
// 提供对命令输出和执行期间发生错误的访问
// 使用场景：完成 Git 操作链时提取结果
func (G *Gcm) Result() ([]byte, error) {
	return G.output, G.errorOnce
}

// Output returns the raw output bytes from command execution
// Provides access to command output without error information
// Use case: extract command results when error handling happens in different places
//
// Output 返回命令执行的原始输出字节
// 提供对命令输出的访问，不包含错误信息
// 使用场景：当错误处理在其他地方进行时提取命令结果
func (G *Gcm) Output() []byte {
	return G.output
}

// Reason returns the first error that happened in the command chains
// Provides access to error information without command output data
// Use case: check error state when output data is not needed
//
// Reason 返回命令链中发生的第一个错误
// 提供对错误信息的访问，不包含命令输出数据
// 使用场景：当不需要输出数据时检查错误状态
func (G *Gcm) Reason() error {
	return G.errorOnce
}

// do executes Git commands with automatic error propagation in chains
// This method prevents execution when previous operations have errors
// Uses Skip3 stack trace to keep log source locations accurate
//
// Design Notes:
// - Must remain unexported to maintain stack trace depths in debug logs
// - Short-circuits when errors exist to prevent cascading failures
// - Returns newWaGcm when commands have errors, newOkGcm when success
//
// do 在链中执行 Git 命令，自动传播错误
// 该方法在之前操作有错误时阻止执行
// 使用 Skip3 堆栈跟踪以保持日志源位置准确
//
// 设计说明：
// - 必须保持非导出以在调试日志中维护堆栈跟踪深度
// - 存在错误时短路以防止级联失败
// - 命令有错误时返回 newWaGcm，成功时返回 newOkGcm
func (G *Gcm) do(name string, args ...string) *Gcm {
	if G.errorOnce != nil {
		return G // Short-circuit: halt execution on existing errors // 短路：存在错误时停止执行
	}
	output, err := G.execConfig.Exec(name, args...)
	if err != nil {
		return newWaGcm(G.execConfig, output, err, G.debugMode)
	}
	return newOkGcm(G.execConfig, output, G.debugMode)
}

// UpdateCommandConfig modifies the execution configuration using provided functions
// Customizes command execution environment when chaining operations
// Use case: adjust execution settings within specific Git operations in chains
//
// UpdateCommandConfig 使用提供的函数修改执行配置
// 在链式操作时自定义命令执行环境
// 使用场景：在链中的特定 Git 操作内调整执行设置
func (G *Gcm) UpdateCommandConfig(updateConfig func(cfg *osexec.CommandConfig)) *Gcm {
	updateConfig(G.execConfig)
	return G
}

// UpdateExecConfig modifies the execution configuration using provided function
// Provides fine-grained management of command execution settings
// Use case: set up execution environment during specialized Git operations
//
// UpdateExecConfig 使用提供的函数修改执行配置
// 提供对命令执行设置的细粒度管理
// 使用场景：在专门的 Git 操作期间配置执行环境
func (G *Gcm) UpdateExecConfig(updateConfig func(cfg *osexec.ExecConfig)) *Gcm {
	updateConfig(G.execConfig)
	return G
}

// WithDebug enables debug mode on the current Gcm instance
// Activates verbose logging to show detailed command execution information
// Use case: troubleshoot Git operations to view detailed command output
//
// WithDebug 在当前 Gcm 实例上启用调试模式
// 激活详细日志记录以显示详细的命令执行信息
// 使用场景：通过查看详细命令输出来排查 Git 操作问题
func (G *Gcm) WithDebug() *Gcm {
	return G.WithDebugMode(true)
}

// WithDebugMode sets debug mode state on the current Gcm instance
// Controls verbose logging output based on the provided boolean flag
// Use case: enable debug output based on environment and settings
//
// WithDebugMode 在当前 Gcm 实例上设置调试模式状态
// 根据提供的布尔标志控制详细日志输出
// 使用场景：根据环境和设置启用调试输出
func (G *Gcm) WithDebugMode(debugMode bool) *Gcm {
	G.debugMode = debugMode
	G.execConfig.WithDebugMode(osexec.NewDebugMode(debugMode))
	return G
}

// ShowDebugMessage shows current execution state with tinted output
// Success messages in green and problem messages in red to console
// Use case: show debug output at specific points in chains
//
// ShowDebugMessage 显示当前执行状态和着色输出
// 成功消息显示为绿色，问题消息显示为红色到控制台
// 使用场景：在链的特定点显示调试输出
func (G *Gcm) ShowDebugMessage() *Gcm {
	switch {
	case G.errorOnce != nil && len(G.output) > 0:
		zaplog.ZAPS.Skip1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(G.output))+"\n", "-")
	case G.errorOnce != nil:
		zaplog.ZAPS.Skip1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.errorOnce), "\n", "-")
	case len(G.output) > 0:
		zaplog.ZAPS.Skip1.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(G.output))+"\n", "-")
	default:
		zaplog.ZAPS.Skip1.SUG.Debugln("done", "\n", "-")
	}
	return G
}

// MustDone terminates execution with panic when errors happen in chains
// Causes immediate program termination when detecting errors in operations
// Use case: enforce strict handling in Git operations that must succeed
//
// MustDone 在链中发生错误时以 panic 终止执行
// 在操作中检测到错误时导致程序立即终止
// 使用场景：在必须成功的 Git 操作中强制执行严格处理
func (G *Gcm) MustDone() *Gcm {
	if G.errorOnce != nil {
		if len(G.output) > 0 {
			zaplog.ZAPS.Skip1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(G.output))+"\n", "-")
		} else {
			zaplog.ZAPS.Skip1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.errorOnce), "\n", "-")
		}
	}
	return G
}

// Must validates command chains success with automatic panic when errors happen
// Naming alternative to MustDone with concise API usage preference
// Use case: enforce strict handling in Git operations that must succeed
//
// Must 验证命令链成功，错误发生时自动 panic
// MustDone 的命名替代，提供简洁 API 使用偏好
// 使用场景：在必须成功的 Git 操作中强制执行严格处理
func (G *Gcm) Must() *Gcm {
	return G.MustDone()
}

// Done validates command chains success with automatic panic when errors happen
// Naming alternative to MustDone with smooth command chains expression
// Use case: enforce strict handling in Git operations that must succeed
//
// Done 验证命令链成功，错误发生时自动 panic
// MustDone 的命名替代，提供流畅的命令链表达
// 使用场景：在必须成功的 Git 操作中强制执行严格处理
func (G *Gcm) Done() *Gcm {
	return G.MustDone()
}

// Nice extracts non-blank output bytes with automatic panic when errors happen
// Combines MustDone validation with non-blank output extraction in one operation
// Use case: obtain command output when blank results indicate errors
//
// Nice 提取非空输出字节，错误发生时自动 panic
// 在一次操作中结合 MustDone 验证和非空输出提取
// 使用场景：当空白结果表示错误时获取命令输出
func (G *Gcm) Nice() []byte {
	return mustslice.Nice(G.MustDone().Output())
}

// Zero validates output is blank with automatic panic when errors happen and non-blank results
// Makes sure command completed with success and produced no output data
// Use case: check commands that should finish without producing output
//
// Zero 验证输出为空，错误发生或非空结果时自动 panic
// 确保命令成功完成且未产生输出数据
// 使用场景：检查应该完成而不产生输出的命令
func (G *Gcm) Zero() {
	mustslice.Zero(G.MustDone().Output())
}

// None validates output is blank with automatic panic when errors happen and non-blank results
// Naming alternative to Zero with API design patterns preference
// Use case: check commands that should finish without producing output
//
// None 验证输出为空，错误发生或非空结果时自动 panic
// Zero 的命名替代，提供 API 设计模式偏好
// 使用场景：检查应该完成而不产生输出的命令
func (G *Gcm) None() {
	mustslice.None(G.MustDone().Output())
}

// When executes provided functions based on condition checks
// Enables condition-based execution within method chains based on states
// Use case: perform Git operations when specific conditions match
//
// When 根据条件检查执行提供的函数
// 基于状态在方法链内启用条件执行
// 使用场景：在满足特定条件时执行 Git 操作
func (G *Gcm) When(condition func(*Gcm) bool, run func(*Gcm) *Gcm) *Gcm {
	if condition(G) {
		return run(G)
	}
	return G
}

// WhenThen executes functions with condition checks and error handling
// Supports complex condition-based logic with error propagation in chains
// Use case: condition-based workflows with error management and validation
//
// WhenThen 执行函数并进行条件检查和错误处理
// 支持在链中进行错误传播的复杂条件逻辑
// 使用场景：具有错误管理和验证的条件工作流
func (G *Gcm) WhenThen(condition func(*Gcm) (bool, error), run func(*Gcm) *Gcm) *Gcm {
	if success, err := condition(G); err != nil {
		return newWaGcm(G.execConfig, []byte{}, err, G.debugMode)
	} else if success {
		return run(G)
	}
	return G
}
