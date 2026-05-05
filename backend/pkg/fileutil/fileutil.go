package fileutil

import (
	"fmt"
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
