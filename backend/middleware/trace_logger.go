package middleware

import (
	"backend/pkg/constants"
	"backend/pkg/global_vars"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TraceLogger 是核心中间件
func TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// A. 获取或生成 Trace ID
		traceID := c.Request.Header.Get(constants.TraceIDHeaderKey)
		if traceID == "" {
			// 如果请求头没有，尝试从 traceparent (W3C standard) 解析，这里简化处理直接生成
			// 生产环境建议解析 traceparent: 00-<trace-id>-<span-id>-<flags>
			traceID = uuid.New().String()
		}

		// B. 将 trace_id 放入 Context，方便 Handler 获取
		c.Set(constants.ContextTraceIDKey, traceID)

		// 可选：将 trace_id 写入响应头，方便前端或调用方追踪
		c.Writer.Header().Set(constants.TraceIDHeaderKey, traceID)

		// C. 创建带有 trace_id 字段的新 Logger
		// 这样所有通过该 logger 打印的日志都会自动带上 {"trace_id": "xxx"}

		traceLogger := global_vars.ZapLog.With(
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		// D. 将新 logger 放入 Context
		c.Set(constants.ContextTraceLoggerKey, traceLogger)

		// E. 执行后续处理
		c.Next()
		appDebug := global_vars.ConfigYml.GetBool("AppDebug")
		if appDebug {
			// F. 请求结束后的日志，如耗时、状态码
			costTime := time.Since(startTime)
			statusCode := c.Writer.Status()

			// 根据状态码选择日志级别
			if statusCode >= 500 {
				traceLogger.Error("request completed",
					zap.Int("status", statusCode),
					zap.Duration("costTime", costTime),
				)
			} else if statusCode >= 400 {
				traceLogger.Warn("request completed",
					zap.Int("status", statusCode),
					zap.Duration("costTime", costTime),
				)
			} else {
				traceLogger.Info("request completed",
					zap.Int("status", statusCode),
					zap.Duration("costTime", costTime),
				)
			}
		}

	}
}
