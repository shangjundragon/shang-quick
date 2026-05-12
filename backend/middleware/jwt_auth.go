// Package middleware 定义 Gin 中间件：CORS、JWT 认证、权限校验、操作日志、链路追踪
package middleware

import (
	"backend/pkg/constants"
	"backend/pkg/jwt"
	"backend/pkg/res_util"
	"backend/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			res_util.Fail(c, res_util.WithCode(401), res_util.WithMsg("未登录"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res_util.Fail(c, res_util.WithCode(401), res_util.WithMsg("Token格式错误"))
			return
		}

		token := parts[1]

		// 检查 Token 是否已被管理员踢下线
		if service.OnlineUserService.IsBlacklisted(c, token) {
			res_util.Fail(c, res_util.WithCode(401), res_util.WithMsg("您的会话已被管理员终止"))
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			res_util.Fail(c, res_util.WithCode(401), res_util.WithMsg("Token无效或已过期"))
			return
		}

		c.Set(constants.ContextUserIDKey, claims.UserID)
		c.Set(constants.ContextUsernameKey, claims.Username)
		c.Set(constants.ContextRoleCodeKey, claims.RoleCode)

		// 异步更新用户最后活跃时间（不阻塞请求）
		go service.OnlineUserService.UpdateActiveTime(c, token)

		c.Next()
	}
}
