package middleware

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/service"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func OperLog(title string, operType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UTC().UnixMilli()

		var requestBody []byte
		if c.Request.Method == "POST" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		username, _ := c.Get(constants.ContextUsernameKey)
		if username == nil {
			username = ""
		}

		status := 1
		if c.Writer.Status() >= 400 {
			status = 0
		}

		operLog := &model.SysOperLog{
			Title:       strPtr(title),
			OperType:    operType,
			Method:      strPtr(c.Request.Method),
			RequestUrl:  strPtr(c.Request.URL.Path),
			RequestData: strPtr(string(requestBody)),
			OperName:    strPtr(username.(string)),
			OperIp:      strPtr(c.ClientIP()),
			OperTime:    startTime,
			Status:      status,
		}

		go service.OperLogService.Save(operLog)
	}
}
