package middleware

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/service"
	"bytes"
	"context"
	"io"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于脱敏日志中的密码字段（正则匹配 JSON 中的 password/oldPassword/newPassword）
var (
	passwordRe    = regexp.MustCompile(`"password"\s*:\s*"[^"]*"`)
	oldPasswordRe = regexp.MustCompile(`"oldPassword"\s*:\s*"[^"]*"`)
	newPasswordRe = regexp.MustCompile(`"newPassword"\s*:\s*"[^"]*"`)
)

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func sanitizeSensitiveData(data string) string {
	if data == "" {
		return data
	}
	sanitized := passwordRe.ReplaceAllString(data, `"password":"***"`)
	sanitized = oldPasswordRe.ReplaceAllString(sanitized, `"oldPassword":"***"`)
	sanitized = newPasswordRe.ReplaceAllString(sanitized, `"newPassword":"***"`)
	return sanitized
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperLog 操作日志中间件：拦截 POST 请求，记录请求/响应数据、操作人和耗时
func OperLog(title string, operType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UTC().UnixMilli()

		// 拦截请求体（只能读取一次，需重新设置）
		var requestBody []byte
		if c.Request.Method == "POST" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 用自定义 ResponseWriter 拦截响应数据
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		responseData := blw.body.String()

		username, _ := c.Get(constants.ContextUsernameKey)
		if username == nil {
			username = ""
		}

		status := 1
		if c.Writer.Status() >= 400 {
			status = 0
		}

		// 脱敏密码字段后异步写入操作日志
		operLog := &model.SysOperLog{
			Title:        strPtr(title),
			OperType:     operType,
			Method:       strPtr(c.Request.Method),
			RequestUrl:   strPtr(c.Request.URL.Path),
			RequestData:  strPtr(sanitizeSensitiveData(string(requestBody))),
			ResponseData: strPtr(sanitizeSensitiveData(responseData)),
			OperName:     strPtr(username.(string)),
			OperIp:       strPtr(c.ClientIP()),
			OperTime:     startTime,
			Status:       status,
		}

		ctx := context.WithValue(context.Background(), constants.ContextTraceIDKey, c.GetString(constants.ContextTraceIDKey))
		go service.OperLogService.Save(ctx, operLog) // 异步写入，不阻塞请求
	}
}
