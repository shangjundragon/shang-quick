package controller

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/pkg/password"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProfileGet(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	userID, _ := c.Get(constants.ContextUserIDKey)
	user, err := service.UserService.GetById(c, userID.(int64))
	if err != nil {
		traceLogger.Error("查询个人信息失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}
	res_util.Success(c, res_util.WithData(user))
}

func ProfileUpdate(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type UpdateReq struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[UpdateReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	userID, _ := c.Get(constants.ContextUserIDKey)
	user := &model.SysUser{
		Id:       userID.(int64),
		Nickname: utils.StrPtr(req.Nickname),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
	}

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("更新个人信息失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("更新失败"))
		return
	}

	updatedUser, _ := service.UserService.GetById(c, userID.(int64))
	res_util.Success(c, res_util.WithData(updatedUser))
}

func ProfileUpdatePwd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type UpdatePwdReq struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	req, err := req_util.BindJson[UpdatePwdReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	if err := password.ValidatePassword(req.NewPassword); err != nil {
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	userID, _ := c.Get(constants.ContextUserIDKey)
	user, err := service.UserService.GetById(c, userID.(int64))
	if err != nil {
		traceLogger.Error("获取用户失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("用户不存在"))
		return
	}

	if !password.Verify(req.OldPassword, user.Password) {
		res_util.Fail(c, res_util.WithMsg("旧密码错误"))
		return
	}

	hashedPwd, _ := password.Hash(req.NewPassword)
	updateUser := &model.SysUser{
		Id:       userID.(int64),
		Password: hashedPwd,
	}

	err = service.UserService.Update(c, updateUser)
	if err != nil {
		traceLogger.Error("修改密码失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("修改失败"))
		return
	}

	res_util.Success(c)
}
