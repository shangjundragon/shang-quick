package service

import (
	"backend/model"
	"backend/pkg/fileutil"
	"backend/pkg/global_vars"
	"context"
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

func (s *fileService) Upload(ctx context.Context, originalName string, fileSize int64, fileType string, reader io.Reader, userID int64) (*model.SysFile, error) {
	now := time.Now()
	uniqueName := fileutil.GenerateUniqueFileName(originalName)
	uploadDir := fileutil.GetUploadDir(global_vars.BasePath, now)
	relativePath := fileutil.GetRelativePath(uniqueName, now)
	fullPath := filepath.Join(uploadDir, uniqueName)

	if err := fileutil.EnsureDir(uploadDir); err != nil {
		return nil, err
	}

	isImg := 0
	if fileutil.IsImage(fileType) {
		isImg = 1
	}

	if isImg == 1 {
		if err := s.compressAndSaveImage(reader, fullPath, fileType); err != nil {
			return nil, err
		}
	} else {
		file, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if _, err := io.Copy(file, reader); err != nil {
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

	_, err := dbw.New[model.SysFile](dbw.WithConfig(global_vars.DbConfig), dbw.WithContext(ctx)).Insert(file)
	if err != nil {
		os.Remove(fullPath)
		return nil, err
	}

	return file, nil
}

func (s *fileService) compressAndSaveImage(reader io.Reader, fullPath string, fileType string) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	quality := global_vars.ConfigYml.GetInt("Upload.ImageCompressQuality")
	if quality <= 0 || quality > 100 {
		quality = 80
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch fileType {
	case "image/png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	}
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

func (s *fileService) Delete(ctx context.Context, id int64) error {
	file, err := s.GetById(ctx, id)
	if err != nil {
		return err
	}
	if file != nil {
		fullPath := filepath.Join(global_vars.BasePath, "storage", "uploads", strings.TrimPrefix(file.FilePath, "/"))
		os.Remove(fullPath)
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
