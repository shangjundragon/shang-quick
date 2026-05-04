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
	JWTSecret          = "shang-quick-admin-secret-key-2026"
	JWTExpireHours     = 24

	MenuTypeDir  = 0
	MenuTypeMenu = 1
	MenuTypeBtn  = 2

	UserStatusEnable  = 1
	UserStatusDisable = 0

	RoleCodeAdmin = "admin"
)
