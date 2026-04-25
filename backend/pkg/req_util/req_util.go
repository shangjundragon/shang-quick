package req_util

import (
	"backend/pkg/constants"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BindJson[T any](c *gin.Context) (T, error) {
	var result T
	if err := c.ShouldBindJSON(&result); err != nil {
		return result, err
	}
	return result, nil
}

func BindJsonWithObj(c *gin.Context, obj any) error {
	// ShouldBindJSON 接收的是 any(interface{}) 类型，obj 必须是结构体指针
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func BindQuery[T any](c *gin.Context) (T, error) {
	var result T
	if err := c.ShouldBindQuery(&result); err != nil {
		return result, err
	}
	return result, nil
}

func GetTraceLogger(c *gin.Context) (traceLogger *zap.Logger, traceID string) {
	// 从 Context 中获取带 trace_id 的 logger
	l, _ := c.Get(constants.ContextTraceLoggerKey)
	traceLogger = l.(*zap.Logger)
	value, _ := c.Get(constants.ContextTraceIDKey)
	return traceLogger, value.(string)
}
