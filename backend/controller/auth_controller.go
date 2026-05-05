package controller

import (
	"backend/pkg/constants"
	"backend/pkg/jwt"
	"backend/pkg/password"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthLogin(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type LoginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	req, err := req_util.BindJson[LoginReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	user, err := service.UserService.GetByUsername(c, req.Username)
	if err != nil || user == nil {
		traceLogger.Warn("用户不存在", zap.Any("req", req), zap.Any("err", err), zap.Any("user", user))
		res_util.Fail(c, res_util.WithMsg("用户不存在"))
		return
	}

	if !password.Verify(req.Password, user.Password) {
		res_util.Fail(c, res_util.WithMsg("密码错误"))
		return
	}

	if user.Status == 0 {
		res_util.Fail(c, res_util.WithMsg("账号已被禁用"))
		return
	}

	roles, _ := service.UserService.GetUserRoles(c, user.Id)
	roleCode := ""
	if len(roles) > 0 {
		roleCode = roles[0].RoleCode
	}

	token, err := jwt.GenerateToken(user.Id, user.Username, roleCode)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("Token生成失败"))
		return
	}

	perms, _ := service.UserService.GetUserPermissions(c, user.Id)
	menus, _ := service.UserService.GetUserMenus(c, user.Id)

	res_util.Success(c, res_util.WithData(gin.H{
		"token":       token,
		"userInfo":    user,
		"roleCode":    roleCode,
		"permissions": perms,
		"menus":       menus,
	}))
}

func AuthInfo(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	userId, _ := c.Get(constants.ContextUserIDKey)
	user, err := service.UserService.GetById(c, userId.(int64))
	if err != nil || user == nil {
		traceLogger.Error("获取用户信息失败", zap.Any("userId", userId), zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("用户不存在"))
		return
	}

	roles, _ := service.UserService.GetUserRoles(c, user.Id)
	roleCode := ""
	if len(roles) > 0 {
		roleCode = roles[0].RoleCode
	}

	perms, _ := service.UserService.GetUserPermissions(c, user.Id)
	menus, _ := service.UserService.GetUserMenus(c, user.Id)

	res_util.Success(c, res_util.WithData(gin.H{
		"userInfo":    user,
		"roleCode":    roleCode,
		"permissions": perms,
		"menus":       menus,
	}))
}
