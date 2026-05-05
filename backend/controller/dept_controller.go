package controller

import (
	"backend/model"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DeptList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	list, err := service.DeptService.List(c)
	if err != nil {
		traceLogger.Error("查询部门列表失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}
	res_util.Success(c, res_util.WithData(list))
}

func DeptAdd(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type AddReq struct {
		ParentId int64  `json:"parentId,string"`
		DeptName string `json:"deptName" binding:"required"`
		OrderNum int    `json:"orderNum"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	dept := &model.SysDept{
		ParentId: req.ParentId,
		DeptName: req.DeptName,
		OrderNum: req.OrderNum,
		Leader:   utils.StrPtr(req.Leader),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
		Status:   1,
	}

	err = service.DeptService.Add(c, dept)
	if err != nil {
		traceLogger.Error("新增部门失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

func DeptEdit(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type EditReq struct {
		Id       int64  `json:"id,string" binding:"required"`
		ParentId int64  `json:"parentId,string"`
		DeptName string `json:"deptName" binding:"required"`
		OrderNum int    `json:"orderNum"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	dept := &model.SysDept{
		Id:       req.Id,
		ParentId: req.ParentId,
		DeptName: req.DeptName,
		OrderNum: req.OrderNum,
		Leader:   utils.StrPtr(req.Leader),
		Phone:    utils.StrPtr(req.Phone),
		Email:    utils.StrPtr(req.Email),
	}

	err = service.DeptService.Update(c, dept)
	if err != nil {
		traceLogger.Error("编辑部门失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

func DeptDelete(c *gin.Context) {
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

	err = service.DeptService.Delete(c, req.Id)
	if err != nil {
		traceLogger.Error("删除部门失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	res_util.Success(c)
}
