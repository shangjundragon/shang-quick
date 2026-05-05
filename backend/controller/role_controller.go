package controller

import (
	"backend/model"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RoleList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ListReq struct {
		PageNum  int    `form:"pageNum" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		RoleName string `form:"roleName"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.RoleService.List(c, req.PageNum, req.PageSize, req.RoleName)
	if err != nil {
		traceLogger.Error("查询角色列表失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}

func RoleAdd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type AddReq struct {
		RoleName string `json:"roleName" binding:"required"`
		RoleCode string `json:"roleCode" binding:"required"`
		Remark   string `json:"remark"`
		Status   int    `json:"status"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	role := &model.SysRole{
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		Remark:   utils.StrPtr(req.Remark),
		Status:   req.Status,
	}

	err = service.RoleService.Add(c, role)
	if err != nil {
		traceLogger.Error("新增角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

func RoleEdit(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type EditReq struct {
		Id       int64  `json:"id,string" binding:"required"`
		RoleName string `json:"roleName" binding:"required"`
		RoleCode string `json:"roleCode" binding:"required"`
		Remark   string `json:"remark"`
		Status   int    `json:"status"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	role := &model.SysRole{
		Id:       req.Id,
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		Remark:   utils.StrPtr(req.Remark),
		Status:   req.Status,
	}

	err = service.RoleService.Update(c, role)
	if err != nil {
		traceLogger.Error("编辑角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

func RoleDelete(c *gin.Context) {
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

	err = service.RoleService.Delete(c, req.Id)
	if err != nil {
		traceLogger.Error("删除角色失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	res_util.Success(c)
}

func RoleMenuIds(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	roleIdStr := c.Query("roleId")
	if roleIdStr == "" {
		traceLogger.Warn("参数错误")
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	var roleId int64
	fmt.Sscanf(roleIdStr, "%d", &roleId)

	menuIds, err := service.RoleService.GetMenuIds(c, roleId)
	if err != nil {
		traceLogger.Error("查询角色菜单失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	menuIdStrs := make([]string, len(menuIds))
	for i, id := range menuIds {
		menuIdStrs[i] = fmt.Sprintf("%d", id)
	}

	res_util.Success(c, res_util.WithData(menuIdStrs))
}

func RoleAssignMenu(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type AssignReq struct {
		RoleId  int64   `json:"roleId,string" binding:"required"`
		MenuIds []int64 `json:"menuIds"`
	}

	req, err := req_util.BindJson[AssignReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.RoleService.AssignMenu(c, req.RoleId, req.MenuIds)
	if err != nil {
		traceLogger.Error("分配菜单失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("分配失败"))
		return
	}

	res_util.Success(c)
}
