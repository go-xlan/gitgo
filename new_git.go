// Package gitgo: Streamlined Git command execution engine with fluent chaining interface
// Provides comprehensive Git operations through os/exec wrapper with smart error handling and debug support
// Features method chaining for complex Git workflows and integrated logging for development and production use
// Supports all major Git operations including commit, push, pull, branch management, and advanced repository querying
//
// gitgo: 流式 Git 命令执行引擎，带有流畅的链式调用接口
// 通过 os/exec 包装器提供全面的 Git 操作，具有智能错误处理和调试支持
// 支持复杂 Git 工作流的方法链式调用，集成日志记录适用于开发和生产环境
// 支持所有主要 Git 操作，包括提交、推送、拉取、分支管理和高级仓库查询
package gitgo

import (
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must/mustslice"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

// Gcm represents a Git Command Manager with chaining support and integrated error handling
// Maintains execution state, output capture, and debug information across method calls
// Supports fluent interface for complex Git workflows with automatic error propagation
//
// Gcm 代表 Git 命令管理器，支持链式调用和集成的错误处理
// 在方法调用间维护执行状态、输出捕获和调试信息
// 支持复杂 Git 工作流的流畅接口，具有自动错误传播功能
type Gcm struct {
	execConfig *osexec.ExecConfig // Execution configuration and path context // 执行配置和路径上下文
	output     []byte             // Last command output bytes // 最后一个命令的输出字节
	errorOnce  error              // First error encountered in chain // 链中遇到的第一个错误
	debugMode  bool               // Enable debug logging output // 启用调试日志输出
}

// New creates a new Gcm instance with default configuration for the specified path
// Initializes Git command manager with standard debug settings and execution context
// Returns ready-to-use Gcm instance for chaining Git operations
//
// New 为指定路径创建具有默认配置的新 Gcm 实例
// 使用标准调试设置和执行上下文初始化 Git 命令管理器
// 返回可用于链式 Git 操作的 Gcm 实例
func New(path string) *Gcm {
	return newOkGcm(osexec.NewCommandConfig().WithPath(path).WithDebugMode(osexec.NewDebugMode(debugModeOpen)), make([]byte, 0), debugModeOpen)
}

// NewGcm creates a new Gcm instance with custom execution configuration
// Allows advanced configuration of command execution environment and debug behavior
// Provides flexibility for specialized Git operation requirements
//
// NewGcm 使用自定义执行配置创建新的 Gcm 实例
// 允许高级配置命令执行环境和调试行为
// 为专门的 Git 操作需求提供灵活性
func NewGcm(path string, execConfig *osexec.ExecConfig) *Gcm {
	return newOkGcm(execConfig.NewConfig().WithPath(path).WithDebugMode(osexec.NewDebugMode(debugModeOpen)), make([]byte, 0), debugModeOpen)
}

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

// Result returns the final output and error from the command chain
// Provides access to last command output and any error that occurred during execution
// Use case: extract results after completing a chain of Git operations
//
// Result 返回命令链的最终输出和错误
// 提供对最后一个命令输出和执行期间发生的任何错误的访问
// 使用场景：完成 Git 操作链后提取结果
func (G *Gcm) Result() ([]byte, error) {
	return G.output, G.errorOnce
}

func (G *Gcm) Output() []byte {
	return G.output
}

func (G *Gcm) Reason() error {
	return G.errorOnce
}

// 这个函数不要使用导出的，因为日志里是跳过3层调用栈，假如使用导出的，跳出的栈的层数就不正确啦
func (G *Gcm) do(name string, args ...string) *Gcm {
	if G.errorOnce != nil {
		return G //当出错时就不要再往下执行，直接在这里拦住，这样整个链式的后续动作就都不执行
	}
	output, err := G.execConfig.Exec(name, args...)
	if err != nil {
		return newWaGcm(G.execConfig, output, err, G.debugMode)
	}
	return newOkGcm(G.execConfig, output, G.debugMode)
}

// UpdateCommandConfig modifies the execution configuration using provided function
// Allows customization of command execution environment during chain operations
// Use case: adjust execution parameters for specific Git operations mid-chain
//
// UpdateCommandConfig 使用提供的函数修改执行配置
// 允许在链式操作期间自定义命令执行环境
// 使用场景：为特定 Git 操作调整执行参数
func (G *Gcm) UpdateCommandConfig(updateConfig func(cfg *osexec.CommandConfig)) *Gcm {
	updateConfig(G.execConfig)
	return G
}

// UpdateExecConfig modifies the execution configuration using provided function
// Provides fine-grained control over command execution settings
// Use case: dynamically configure execution environment for specialized Git operations
//
// UpdateExecConfig 使用提供的函数修改执行配置
// 提供对命令执行设置的细粒度控制
// 使用场景：为专门的 Git 操作动态配置执行环境
func (G *Gcm) UpdateExecConfig(updateConfig func(cfg *osexec.ExecConfig)) *Gcm {
	updateConfig(G.execConfig)
	return G
}

// WithDebug enables debug mode for the current Gcm instance
// Activates verbose logging to show detailed command execution information
// Use case: troubleshoot Git operations by seeing detailed command output
//
// WithDebug 为当前 Gcm 实例启用调试模式
// 激活详细日志记录以显示详细的命令执行信息
// 使用场景：通过查看详细命令输出来排查 Git 操作问题
func (G *Gcm) WithDebug() *Gcm {
	return G.WithDebugMode(true)
}

// WithDebugMode sets debug mode state for the current Gcm instance
// Controls verbose logging output based on the provided boolean flag
// Use case: conditionally enable debug output based on environment or user settings
//
// WithDebugMode 为当前 Gcm 实例设置调试模式状态
// 根据提供的布尔标志控制详细日志输出
// 使用场景：根据环境或用户设置有条件地启用调试输出
func (G *Gcm) WithDebugMode(debugMode bool) *Gcm {
	G.debugMode = debugMode
	G.execConfig.WithDebugMode(osexec.NewDebugMode(debugMode))
	return G
}

// ShowDebugMessage displays current execution state with colored output
// Shows success messages in green or error messages in red for visual debugging
// Use case: manually trigger debug output display at specific points in the chain
//
// ShowDebugMessage 显示当前执行状态和彩色输出
// 成功消息显示为绿色，错误消息显示为红色，便于可视化调试
// 使用场景：在链的特定点手动触发调试输出显示
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

// MustDone terminates execution with panic if any error occurred in the chain
// CRITICAL: causes immediate program termination if errors are present
// Use case: enforce strict error handling in critical Git operations that must succeed
//
// MustDone 如果链中发生任何错误则以 panic 终止执行
// 关键：如果存在错误则导致程序立即终止
// 使用场景：在必须成功的关键 Git 操作中强制执行严格的错误处理
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

func (G *Gcm) Must() *Gcm {
	return G.MustDone()
}

func (G *Gcm) Done() *Gcm {
	return G.MustDone()
}

func (G *Gcm) Nice() []byte {
	return mustslice.Nice(G.MustDone().Output())
}

func (G *Gcm) Zero() {
	mustslice.Zero(G.MustDone().Output())
}

func (G *Gcm) None() {
	mustslice.None(G.MustDone().Output())
}

// When executes provided function conditionally based on boolean evaluation
// Allows conditional execution within the method chain based on current state
// Use case: perform Git operations only when specific conditions are met
//
// When 根据布尔评估有条件地执行提供的函数
// 允许基于当前状态在方法链内有条件执行
// 使用场景：仅在满足特定条件时执行 Git 操作
func (G *Gcm) When(condition func(*Gcm) bool, run func(*Gcm) *Gcm) *Gcm {
	if condition(G) {
		return run(G)
	}
	return G
}

// WhenThen executes function conditionally with error handling
// Supports complex conditional logic with error propagation in the chain
// Use case: advanced conditional workflows with proper error management
//
// WhenThen 有条件地执行函数并进行错误处理
// 支持在链中进行错误传播的复杂条件逻辑
// 使用场景：具有适当错误管理的高级条件工作流
func (G *Gcm) WhenThen(condition func(*Gcm) (bool, error), run func(*Gcm) *Gcm) *Gcm {
	if success, err := condition(G); err != nil {
		return newWaGcm(G.execConfig, []byte{}, err, G.debugMode)
	} else if success {
		return run(G)
	}
	return G
}
