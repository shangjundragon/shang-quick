// Package constants 定义全局常量和上下文键，保持字符串一致
package constants

var (
	TraceIDHeaderKey      = "X-Trace-Id"
	ContextTraceIDKey     = "trace_id"
	ContextTraceLoggerKey = "trace_logger"
	Y                     = "Y"
	N                     = "N"

	ContextUserIDKey   = "user_id"
	ContextUsernameKey = "username"
	ContextRoleCodeKey = "role_code"
	AdminUserId        = 1

	MenuTypeDir  = 0
	MenuTypeMenu = 1
	MenuTypeBtn  = 2

	UserStatusEnable  = 1
	UserStatusDisable = 0

	RoleCodeAdmin = "admin"
)
