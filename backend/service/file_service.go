package service

import (
	"backend/model"
	"backend/pkg/fileutil"
	"backend/pkg/global_vars"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shangjundragon/dbw"
)

var FileService = new(fileService)

type fileService struct{}

func (s *fileService) getAllowedExtensions() []string {
	exts := global_vars.ConfigYml.GetStringSlice("Upload.AllowedExtensions")
	if len(exts) == 0 {
		return []string{"*"}
	}
	return exts
}

func (s *fileService) isExtensionAllowed(ext string) bool {
	allowed := s.getAllowedExtensions()
	for _, allowedExt := range allowed {
		if allowedExt == "*" {
			return true
		}
		if strings.EqualFold(allowedExt, ext) {
			return true
		}
	}
	return false
}

// Upload 上传文件：先写临时文件，再校验类型，然后压缩图片（如需要），最后移至上传目录
func (s *fileService) Upload(ctx context.Context, originalName string, fileSize int64, fileType string, reader io.Reader, userID int64) (*model.SysFile, error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	if !s.isExtensionAllowed(ext) {
		return nil, fmt.Errorf("file extension '%s' is not allowed", ext)
	}

	now := time.Now()
	uniqueName := fileutil.GenerateUniqueFileName(originalName)
	uploadDir := fileutil.GetUploadDir(global_vars.BasePath, now)
	relativePath := fileutil.GetRelativePath(uniqueName, now)
	fullPath := filepath.Join(uploadDir, uniqueName)

	if err := fileutil.EnsureDir(uploadDir); err != nil {
		return nil, err
	}

	// 先写入临时文件，避免直接污染上传目录
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return nil, err
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath) // 无论成功或失败，清理临时文件

	if _, err := io.Copy(tempFile, reader); err != nil {
		tempFile.Close()
		return nil, err
	}
	tempFile.Close()

	isImg := 0
	if fileutil.IsImage(fileType) {
		isImg = 1
	}

	// 图片文件检测：比对 Content-Type 和实际文件头，防止类型伪造
	if isImg == 1 {
		f, err := os.Open(tempPath)
		if err != nil {
			return nil, err
		}
		detectedType, err := fileutil.DetectFileType(f)
		f.Close()
		if err != nil {
			return nil, err
		}
		normalizedType := fileutil.NormalizeContentType(fileType)
		if normalizedType != "" && normalizedType != detectedType {
			return nil, fmt.Errorf("file type mismatch: declared '%s' but detected '%s'", fileType, detectedType)
		}
	}

	if err := os.Rename(tempPath, fullPath); err != nil {
		return nil, err
	}

	if isImg == 1 {
		if err := s.compressAndSaveImage(fullPath, fileutil.NormalizeContentType(fileType)); err != nil {
			os.Remove(fullPath)
			return nil, err
		}
	}

	file := &model.SysFile{
		OriginalName: originalName,
		FileName:     uniqueName,
		FilePath:     relativePath,
		FileSize:     fileSize,
		FileType:     fileType,
		IsImage:      isImg,
		CreateBy:     userID,
	}

	_, err = dbw.New[model.SysFile](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(file)
	if err != nil {
		os.Remove(fullPath)
		return nil, err
	}

	return file, nil
}

const maxImageDimension = 10000

// compressAndSaveImage 压缩并重写图片（JPEG 压缩质量可配，PNG 仅重编码），通过临时文件+原子重命名保证安全
func (s *fileService) compressAndSaveImage(fullPath string, fileType string) error {
	if fileType != "image/jpeg" && fileType != "image/png" {
		return nil
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	if bounds.Dx() > maxImageDimension || bounds.Dy() > maxImageDimension {
		return fmt.Errorf("image dimensions too large, maximum supported is %dx%d", maxImageDimension, maxImageDimension)
	}

	quality := global_vars.ConfigYml.GetInt("Upload.ImageCompressQuality")
	if quality < 0 || quality > 100 {
		quality = 80
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(fullPath), "compress-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	var encodeErr error
	switch fileType {
	case "image/png":
		encodeErr = png.Encode(tmpFile, img)
	default:
		encodeErr = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: quality})
	}
	tmpFile.Close()

	if encodeErr != nil {
		return encodeErr
	}

	return os.Rename(tmpPath, fullPath)
}

func (s *fileService) List(ctx context.Context, pageNum, pageSize int, originalName string) ([]model.SysFile, int64, error) {
	wrapper := dbw.New[model.SysFile](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx))
	if originalName != "" {
		wrapper = wrapper.Like("original_name", "%"+originalName+"%")
	}
	list, total, err := wrapper.OrderByDesc("create_time").SelectPage(pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Delete 删除文件记录并移除物理文件
func (s *fileService) Delete(ctx context.Context, id int64) error {
	file, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}
	if file != nil {
		fullPath := filepath.Join(global_vars.BasePath, "storage", "uploads", strings.TrimPrefix(file.FilePath, "/"))
		os.Remove(fullPath) // 删除物理文件，忽略错误（文件可能已被手动删除）
	}
	_, err = dbw.New[model.SysFile](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).DeleteById(id)
	return err
}

func (s *fileService) GetById(ctx context.Context, id int64) (*model.SysFile, error) {
	return dbw.New[model.SysFile](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).SelectById(id)
}

func (s *fileService) GetFileUrlPrefix() string {
	return global_vars.ConfigYml.GetString("Upload.FileUrlPrefix")
}

func (s *fileService) IsPublicAccess() bool {
	return global_vars.ConfigYml.GetBool("Upload.PublicAccess")
}

func (s *fileService) GetMaxSizeMB() int {
	return global_vars.ConfigYml.GetInt("Upload.MaxSize")
}
