package router

import (
	"backend/middleware"
	"backend/pkg/res_util"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用 CORS 中间件
	r.Use(middleware.CORS())
	// 使用 TraceLogger 中间件
	r.Use(middleware.TraceLogger())

	// API 路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			res_util.Success(c, res_util.WithData("ok"))
		})

	}

	return r
}
