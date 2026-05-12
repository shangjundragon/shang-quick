// Package utils 提供通用工具函数
package utils

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
