package controller

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/pkg/global_vars"
	"backend/pkg/password"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/shangjundragon/dbw"
)

func UserList(c *gin.Context) {
	type ListReq struct {
		PageNum  int    `form:"pageNum" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		Username string `form:"username"`
		Phone    string `form:"phone"`
		Status   *int   `form:"status"`
		DeptId   *int64 `form:"deptId"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.UserService.List(req.PageNum, req.PageSize, req.Username, req.Phone, req.Status, req.DeptId)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}

func UserAdd(c *gin.Context) {
	type AddReq struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		Nickname string  `json:"nickname"`
		Phone    string  `json:"phone"`
		Email    string  `json:"email"`
		DeptId   int64   `json:"deptId"`
		RoleIds  []int64 `json:"roleIds"`
		Status   int     `json:"status"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

		hashedPwd, _ := password.Hash(req.Password)
	user := &model.SysUser{
		Username: req.Username,
		Password: hashedPwd,
		Nickname: utils.StrPtr(req.Nickname),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
		DeptId:   req.DeptId,
		Status:   req.Status,
	}

	createBy, _ := c.Get(constants.ContextUserIDKey)
	user.CreateBy = createBy.(int64)

	err = service.UserService.Add(user)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	for _, roleId := range req.RoleIds {
		dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig)).Insert(
			&model.SysUserRole{UserId: user.Id, RoleId: roleId})
	}

	res_util.Success(c)
}

func UserEdit(c *gin.Context) {
	type EditReq struct {
		Id       int64   `json:"id" binding:"required"`
		Nickname string  `json:"nickname"`
		Phone    string  `json:"phone"`
		Email    string  `json:"email"`
		DeptId   int64   `json:"deptId"`
		RoleIds  []int64 `json:"roleIds"`
		Status   int     `json:"status"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	user := &model.SysUser{
		Id:       req.Id,
		Nickname: utils.StrPtr(req.Nickname),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
		DeptId:   req.DeptId,
		Status:   req.Status,
	}

	updateBy, _ := c.Get(constants.ContextUserIDKey)
	user.UpdateBy = updateBy.(int64)

	err = service.UserService.Update(user)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

func UserChangeStatus(c *gin.Context) {
	type StatusReq struct {
		Id     int64 `json:"id" binding:"required"`
		Status int   `json:"status" binding:"required"`
	}

	req, err := req_util.BindJson[StatusReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	user := &model.SysUser{
		Id:     req.Id,
		Status: req.Status,
	}
	updateBy, _ := c.Get(constants.ContextUserIDKey)
	user.UpdateBy = updateBy.(int64)

	err = service.UserService.Update(user)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}

	res_util.Success(c)
}

func UserResetPwd(c *gin.Context) {
	type ResetReq struct {
		Id int64 `json:"id" binding:"required"`
	}

	req, err := req_util.BindJson[ResetReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	hashedPwd, _ := password.Hash("123456")
	user := &model.SysUser{
		Id:       req.Id,
		Password: hashedPwd,
	}

	err = service.UserService.Update(user)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("重置失败"))
		return
	}

	res_util.Success(c)
}

func UserDelete(c *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.UserService.Delete(req.Id)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("删除失败"))
		return
	}

	res_util.Success(c)
}
