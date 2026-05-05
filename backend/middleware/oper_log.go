package middleware

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/service"
	"bytes"
	"context"
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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperLog(title string, operType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UTC().UnixMilli()

		var requestBody []byte
		if c.Request.Method == "POST" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

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

		operLog := &model.SysOperLog{
			Title:        strPtr(title),
			OperType:     operType,
			Method:       strPtr(c.Request.Method),
			RequestUrl:   strPtr(c.Request.URL.Path),
			RequestData:  strPtr(string(requestBody)),
			ResponseData: strPtr(responseData),
			OperName:     strPtr(username.(string)),
			OperIp:       strPtr(c.ClientIP()),
			OperTime:     startTime,
			Status:       status,
		}

		go service.OperLogService.Save(context.Background(), operLog)
	}
}
