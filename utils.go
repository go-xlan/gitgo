package gitgo

// Global debug mode flag controls logging output during Gcm instance execution
// Shows detailed command execution and output information upon enabling
// 全局调试模式标志在 Gcm 实例执行期间控制日志输出
// 启用时显示详细的命令执行和输出信息
var debugModeOpen = false

// SetDebugMode enables and disables package-level debug logging on Git operations
// Controls verbose output during subsequent Git command executions
// Use case: enable detailed logging in development and troubleshooting
//
// SetDebugMode 在 Git 操作上启用和禁用全局调试日志
// 在后续 Git 命令执行期间控制详细输出
// 使用场景：在开发和故障排除中启用详细日志记录
func SetDebugMode(enable bool) {
	debugModeOpen = enable
}
