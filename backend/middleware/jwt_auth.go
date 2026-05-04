package middleware

import (
	"backend/pkg/constants"
	"backend/pkg/jwt"
	"backend/pkg/res_util"
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

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			res_util.Fail(c, res_util.WithCode(401), res_util.WithMsg("Token无效或已过期"))
			return
		}

		c.Set(constants.ContextUserIDKey, claims.UserID)
		c.Set(constants.ContextUsernameKey, claims.Username)
		c.Set(constants.ContextRoleCodeKey, claims.RoleCode)
		c.Next()
	}
}
