package controller

import (
	"backend/model"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/pkg/utils"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func DeptList(c *gin.Context) {
	list, err := service.DeptService.List()
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}
	res_util.Success(c, res_util.WithData(list))
}

func DeptAdd(c *gin.Context) {
	type AddReq struct {
		ParentId int64  `json:"parentId"`
		DeptName string `json:"deptName" binding:"required"`
		OrderNum int    `json:"orderNum"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[AddReq](c)
	if err != nil {
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

	err = service.DeptService.Add(dept)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("新增失败"))
		return
	}

	res_util.Success(c)
}

func DeptEdit(c *gin.Context) {
	type EditReq struct {
		Id       int64  `json:"id" binding:"required"`
		ParentId int64  `json:"parentId"`
		DeptName string `json:"deptName" binding:"required"`
		OrderNum int    `json:"orderNum"`
		Leader   string `json:"leader"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	req, err := req_util.BindJson[EditReq](c)
	if err != nil {
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

	err = service.DeptService.Update(dept)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("编辑失败"))
		return
	}

	res_util.Success(c)
}

func DeptDelete(c *gin.Context) {
	type DeleteReq struct {
		Id int64 `json:"id" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.DeptService.Delete(req.Id)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg(err.Error()))
		return
	}

	res_util.Success(c)
}
