package res_util

import (
	"backend/pkg/constants"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
	TraceId string `json:"trace_id"`
}

type Option func(opts *Response)

func NewR(opts ...Option) *Response {
	r := &Response{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithCode(code int) Option {
	return func(r *Response) {
		r.Code = code
	}
}

func WithMsg(msg string) Option {
	return func(r *Response) {
		r.Message = msg
	}
}

func WithData(data any) Option {
	return func(r *Response) {
		r.Data = data
	}
}

func Success(c *gin.Context, opts ...Option) {
	r := NewR(opts...)
	if r.Code == 0 {
		r.Code = 200
	}
	if r.Message == "" {
		r.Message = "success"
	}

	r.TraceId = c.GetString(constants.ContextTraceIDKey)
	c.JSON(200, r)
	c.Abort()
}

func Fail(c *gin.Context, opts ...Option) {
	r := NewR(opts...)
	if r.Code == 0 {
		r.Code = 500
	}
	if r.Message == "" {
		r.Message = "fail"
	}
	r.TraceId = c.GetString(constants.ContextTraceIDKey)
	c.JSON(200, r)
	c.Abort()
}
