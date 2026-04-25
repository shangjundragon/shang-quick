package req_util

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetParam struct {
	c            *gin.Context
	require      bool
	defaultValue any
	key          string
}

type ParamOption func(param *GetParam)

func WithParamRequire(required bool) ParamOption {
	return func(param *GetParam) {
		param.require = required
	}
}

func WithParamDefaultValue(value any) ParamOption {
	return func(param *GetParam) {
		param.defaultValue = value
	}
}

// NewGetParam 创建新的 GetParam 实例
func NewGetParam(c *gin.Context, key string, opts ...ParamOption) *GetParam {
	gq := &GetParam{
		c:   c,
		key: key,
	}
	for _, opt := range opts {
		opt(gq)
	}
	return gq
}

// String 获取字符串类型参数
func (gq *GetParam) String() (string, error) {
	value := gq.c.Param(gq.key)
	if value == "" {
		if gq.require {
			return "", fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(string); ok {
				return v, nil
			}
			return "", errors.New("default value type mismatch")
		}
		return "", nil
	}
	return value, nil
}

// Int 获取整数类型参数
func (gq *GetParam) Int() (int, error) {
	value := gq.c.Param(gq.key)
	if value == "" {
		if gq.require {
			return 0, fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(int); ok {
				return v, nil
			}
			return 0, errors.New("default value type mismatch")
		}
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be an integer: %w", gq.key, err)
	}
	return intValue, nil
}

// Int64 获取 int64 类型参数
func (gq *GetParam) Int64() (int64, error) {
	value := gq.c.Param(gq.key)
	if value == "" {
		if gq.require {
			return 0, fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(int64); ok {
				return v, nil
			}
			return 0, errors.New("default value type mismatch")
		}
		return 0, nil
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be an integer: %w", gq.key, err)
	}
	return intValue, nil
}

// Bool 获取布尔类型参数
func (gq *GetParam) Bool() (bool, error) {
	value := gq.c.Param(gq.key)
	if value == "" {
		if gq.require {
			return false, fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(bool); ok {
				return v, nil
			}
			return false, errors.New("default value type mismatch")
		}
		return false, nil
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("parameter %s must be a boolean: %w", gq.key, err)
	}
	return boolValue, nil
}

// Float64 获取浮点数类型参数
func (gq *GetParam) Float64() (float64, error) {
	value := gq.c.Param(gq.key)

	if value == "" {
		if gq.require {
			return 0, fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(float64); ok {
				return v, nil
			}
			return 0, errors.New("default value type mismatch")
		}
		return 0, nil
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be a float: %w", gq.key, err)
	}
	return floatValue, nil
}

// Any 获取任意类型（原始字符串值）
func (gq *GetParam) Any() (string, error) {
	value := gq.c.Param(gq.key)
	if value == "" {
		if gq.require {
			return "", fmt.Errorf("parameter %s is required", gq.key)
		}
		if gq.defaultValue != nil {
			if v, ok := gq.defaultValue.(string); ok {
				return v, nil
			}
			return "", errors.New("default value type mismatch")
		}
		return "", nil
	}
	return value, nil
}
