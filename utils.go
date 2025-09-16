package gitgo

// Global debug mode flag controls logging output across all Gcm instances
// When enabled, shows detailed command execution and output information
// 全局调试模式标志控制所有 Gcm 实例的日志输出
// 启用时显示详细的命令执行和输出信息
var debugModeOpen = false

// SetDebugMode enables or disables global debug logging for Git operations
// Controls verbose output for all subsequent Git command executions
// Use case: enable detailed logging during development or troubleshooting
//
// SetDebugMode 启用或禁用 Git 操作的全局调试日志
// 控制所有后续 Git 命令执行的详细输出
// 使用场景：开发或故障排除期间启用详细日志记录
func SetDebugMode(enable bool) {
	debugModeOpen = enable
}
