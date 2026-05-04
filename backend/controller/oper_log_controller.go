package controller

import (
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func OperLogList(c *gin.Context) {
	type ListReq struct {
		PageNum  int    `form:"pageNum" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		Title    string `form:"title"`
		OperName string `form:"operName"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.OperLogService.List(req.PageNum, req.PageSize, req.Title, req.OperName)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}
