package controller

import (
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func OperLogList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ListReq struct {
		PageNum  int    `form:"pageNum" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		Title    string `form:"title"`
		OperName string `form:"operName"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.OperLogService.List(c, req.PageNum, req.PageSize, req.Title, req.OperName)
	if err != nil {
		traceLogger.Error("查询操作日志失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}
