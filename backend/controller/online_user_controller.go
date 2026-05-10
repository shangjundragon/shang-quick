package controller

import (
	"backend/pkg/constants"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func OnlineUserList(c *gin.Context) {
	list, err := service.OnlineUserService.List(c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("获取在线用户失败"))
		return
	}

	res_util.Success(c, res_util.WithData(list))
}

func OnlineUserKick(c *gin.Context) {
	type KickReq struct {
		TokenId string `json:"tokenId" binding:"required"`
	}

	req, err := req_util.BindJson[KickReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	currentUsername, _ := c.Get(constants.ContextUsernameKey)

	list, err := service.OnlineUserService.List(c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("获取在线用户失败"))
		return
	}
	var targetUser *string
	var targetUserID int64
	for _, user := range list {
		if user.TokenId == req.TokenId {
			targetUser = &user.Username
			targetUserID = user.UserId
			break
		}
	}

	if targetUser == nil {
		res_util.Fail(c, res_util.WithMsg("用户不在线"))
		return
	}

	if *targetUser == currentUsername.(string) {
		res_util.Fail(c, res_util.WithMsg("不能踢出自己"))
		return
	}

	isAdmin, err := service.UserService.IsAdmin(c, targetUserID)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}
	if isAdmin {
		res_util.Fail(c, res_util.WithMsg("不能踢出超级管理员"))
		return
	}

	err = service.OnlineUserService.KickUser(c, req.TokenId)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	res_util.Success(c, res_util.WithMsg("踢出成功"))
}
