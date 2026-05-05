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
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shangjundragon/dbw"
	"go.uber.org/zap"
)

func UserList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
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

	list, total, err := service.UserService.List(c, req.PageNum, req.PageSize, req.Username, req.Phone, req.Status, req.DeptId)
	if err != nil {
		traceLogger.Error("查询失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	type UserVO struct {
		model.SysUser
		DeptName  string   `json:"deptName"`
		RoleNames []string `json:"roleNames"`
	}

	voList := make([]UserVO, len(list))
	for i, user := range list {
		dept, _ := service.DeptService.GetById(c, user.DeptId)
		deptName := ""
		if dept != nil {
			deptName = dept.DeptName
		}

		roles, _ := service.UserService.GetUserRoles(c, user.Id)
		roleNames := make([]string, 0, len(roles))
		for _, role := range roles {
			roleNames = append(roleNames, role.RoleName)
		}

		voList[i] = UserVO{
			SysUser:   user,
			DeptName:  deptName,
			RoleNames: roleNames,
		}
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  voList,
		"total": total,
	}))
}

func UserAdd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type AddReq struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		Nickname string  `json:"nickname"`
		Phone    string  `json:"phone"`
		Email    string  `json:"email"`
		DeptId   int64   `json:"deptId,string"`
		RoleIds  []int64 `json:"roleIds"`
		Status   int     `json:"status"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		traceLogger.Error("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	if err := password.ValidatePassword(req.Password); err != nil {
		res_util.Fail(c, res_util.WithMsg(err.Error()))
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

	err = service.UserService.AddWithRoles(c, user, req.RoleIds)
	if err != nil {
		traceLogger.Error("新增失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

func UserEdit(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type EditReq struct {
		Id       int64   `json:"id,string" binding:"required"`
		Nickname string  `json:"nickname"`
		Phone    string  `json:"phone"`
		Email    string  `json:"email"`
		DeptId   int64   `json:"deptId,string"`
		RoleIds  []int64 `json:"roleIds"`
		Status   int     `json:"status"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
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

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("编辑失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(c)).Eq("user_id", req.Id).Delete()
	for _, roleId := range req.RoleIds {
		dbw.New[model.SysUserRole](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(c)).Insert(
			&model.SysUserRole{UserId: req.Id, RoleId: roleId})
	}

	res_util.Success(c)
}

func UserChangeStatus(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type StatusReq struct {
		Id     int64 `json:"id,string" binding:"required"`
		Status int   `json:"status" binding:"required"`
	}

	req, err := req_util.BindJson[StatusReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	user := &model.SysUser{
		Id:     req.Id,
		Status: req.Status,
	}
	updateBy, _ := c.Get(constants.ContextUserIDKey)
	user.UpdateBy = updateBy.(int64)

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("操作失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("操作失败"))
		return
	}

	res_util.Success(c)
}

func UserResetPwd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ResetReq struct {
		Id int64 `json:"id,string" binding:"required"`
	}

	req, err := req_util.BindJson[ResetReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	hashedPwd, _ := password.Hash("123456")
	user := &model.SysUser{
		Id:       req.Id,
		Password: hashedPwd,
	}

	err = service.UserService.Update(c, user)
	if err != nil {
		traceLogger.Error("重置失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("重置失败"))
		return
	}

	res_util.Success(c)
}

func UserDelete(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type DeleteReq struct {
		Id int64 `json:"id,string" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.UserService.Delete(c, req.Id)
	if err != nil {
		traceLogger.Error("删除失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("删除失败"))
		return
	}

	res_util.Success(c)
}

func UserRoleIds(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	var userId int64
	fmt.Sscanf(userIdStr, "%d", &userId)

	roles, err := service.UserService.GetUserRoles(c, userId)
	if err != nil {
		traceLogger.Error("查询失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	roleIds := make([]string, len(roles))
	for i, role := range roles {
		roleIds[i] = fmt.Sprintf("%d", role.Id)
	}

	res_util.Success(c, res_util.WithData(roleIds))
}
