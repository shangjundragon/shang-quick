// Package fileutil 提供文件上传、类型检测、图片压缩和路径处理工具函数
package fileutil

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func GenerateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

func GetUploadDir(basePath string, now time.Time) string {
	return filepath.Join(basePath, "storage", "uploads", now.Format("2006"), now.Format("01"))
}

func GetRelativePath(fileName string, now time.Time) string {
	return "/" + path.Join(now.Format("2006"), now.Format("01"), fileName)
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func IsImage(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	default:
		return false
	}
}

func DetectFileType(reader io.Reader) (string, error) {
	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	buf = buf[:n]

	if len(buf) >= 2 && buf[0] == 0xFF && buf[1] == 0xD8 {
		return "image/jpeg", nil
	}
	if len(buf) >= 8 && string(buf[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png", nil
	}
	if len(buf) >= 6 && (string(buf[:6]) == "GIF87a" || string(buf[:6]) == "GIF89a") {
		return "image/gif", nil
	}
	if len(buf) >= 12 && string(buf[:4]) == "RIFF" && string(buf[8:12]) == "WEBP" {
		return "image/webp", nil
	}
	if len(buf) >= 2 && buf[0] == 'B' && buf[1] == 'M' {
		return "image/bmp", nil
	}

	return "application/octet-stream", nil
}

func NormalizeContentType(contentType string) string {
	if contentType == "image/jpg" {
		return "image/jpeg"
	}
	return contentType
}

func FormatFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
