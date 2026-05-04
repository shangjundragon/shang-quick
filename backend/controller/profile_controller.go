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
)

func ProfileGet(c *gin.Context) {
	userID, _ := c.Get(constants.ContextUserIDKey)
	user, err := service.UserService.GetById(userID.(int64))
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}
	res_util.Success(c, res_util.WithData(user))
}

func ProfileUpdate(c *gin.Context) {
	type UpdateReq struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[UpdateReq](c)
	if err != nil {
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

	err = service.UserService.Update(user)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("更新失败"))
		return
	}

	res_util.Success(c)
}

func ProfileUpdatePwd(c *gin.Context) {
	type UpdatePwdReq struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	req, err := req_util.BindJson[UpdatePwdReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	userID, _ := c.Get(constants.ContextUserIDKey)
	user, err := service.UserService.GetById(userID.(int64))
	if err != nil {
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

	err = service.UserService.Update(updateUser)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("修改失败"))
		return
	}

	res_util.Success(c)
}
