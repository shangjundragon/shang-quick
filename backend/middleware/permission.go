package middleware

import (
	"backend/pkg/constants"
	"backend/pkg/res_util"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func Permission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode, _ := c.Get(constants.ContextRoleCodeKey)
		if roleCode == constants.RoleCodeAdmin {
			c.Next()
			return
		}

		userID, exists := c.Get(constants.ContextUserIDKey)
		if !exists {
			res_util.Fail(c, res_util.WithCode(403), res_util.WithMsg("无权限操作"))
			return
		}

		hasPerm, err := service.UserService.CheckPermission(c, userID.(int64), perm)
		if err != nil || !hasPerm {
			res_util.Fail(c, res_util.WithCode(403), res_util.WithMsg("无权限操作"))
			return
		}
		c.Next()
	}
}
