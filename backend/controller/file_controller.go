package controller

import (
	"backend/model"
	"backend/pkg/constants"
	"backend/pkg/fileutil"
	"backend/pkg/global_vars"
	"backend/pkg/jwt"
	"backend/pkg/req_util"
	"backend/pkg/res_util"
	"backend/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FileUpload(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)

	maxSize := service.FileService.GetMaxSizeMB()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxSize)*1024*1024)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		traceLogger.Warn("获取上传文件失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("请选择要上传的文件"))
		return
	}
	defer file.Close()

	userID, _ := c.Get(constants.ContextUserIDKey)

	uploadedFile, err := service.FileService.Upload(c, header.Filename, header.Size, header.Header.Get("Content-Type"), file, userID.(int64))
	if err != nil {
		traceLogger.Error("上传文件失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("上传失败"))
		return
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"id":           uploadedFile.Id,
		"originalName": uploadedFile.OriginalName,
		"fileName":     uploadedFile.FileName,
		"filePath":     uploadedFile.FilePath,
		"fileSize":     uploadedFile.FileSize,
		"fileSizeStr":  fileutil.FormatFileSize(uploadedFile.FileSize),
		"fileType":     uploadedFile.FileType,
		"isImage":      uploadedFile.IsImage,
		"createTime":   uploadedFile.CreateTime,
	}))
}

func FileList(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type ListReq struct {
		PageNum      int    `form:"pageNum" binding:"required"`
		PageSize     int    `form:"pageSize" binding:"required"`
		OriginalName string `form:"originalName"`
	}

	req, err := req_util.BindQuery[ListReq](c)
	if err != nil {
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	list, total, err := service.FileService.List(c, req.PageNum, req.PageSize, req.OriginalName)
	if err != nil {
		traceLogger.Error("查询文件列表失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("查询失败"))
		return
	}

	type FileVO struct {
		model.SysFile
		FileSizeStr string `json:"fileSizeStr"`
	}

	voList := make([]FileVO, len(list))
	for i, f := range list {
		voList[i] = FileVO{
			SysFile:     f,
			FileSizeStr: fileutil.FormatFileSize(f.FileSize),
		}
	}

	res_util.Success(c, res_util.WithData(gin.H{
		"list":  voList,
		"total": total,
	}))
}

func FileDelete(c *gin.Context) {
	traceLogger, _ := req_util.GetTraceLogger(c)
	type DeleteReq struct {
		Id int64 `json:"id,string" binding:"required"`
	}

	req, err := req_util.BindJson[DeleteReq](c)
	if err != nil {
		traceLogger.Warn("参数错误", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("参数错误"))
		return
	}

	err = service.FileService.Delete(c, req.Id)
	if err != nil {
		traceLogger.Error("删除文件失败", zap.Error(err))
		res_util.Fail(c, res_util.WithMsg("删除失败"))
		return
	}

	res_util.Success(c)
}

func FileConfig(c *gin.Context) {
	res_util.Success(c, res_util.WithData(gin.H{
		"fileUrlPrefix": service.FileService.GetFileUrlPrefix(),
	}))
}

// FileAccess 提供文件访问（直接 HTTP 响应文件内容），含路径遍历防护
func FileAccess(c *gin.Context) {
	filePath := c.Param("filepath")
	if filePath == "" {
		c.AbortWithStatus(404)
		return
	}

	if !service.FileService.IsPublicAccess() {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(401)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(401)
			return
		}
		_, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
	}

	// 防止路径遍历攻击：解析符号链接后检查路径是否仍在上传目录内
	uploadsDir := filepath.Join(global_vars.BasePath, "storage", "uploads")
	fullPath := filepath.Join(uploadsDir, strings.TrimPrefix(filePath, "/"))

	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	resolvedUploads, _ := filepath.EvalSymlinks(uploadsDir)
	if !strings.HasPrefix(resolvedPath, resolvedUploads+string(filepath.Separator)) && resolvedPath != resolvedUploads {
		c.AbortWithStatus(403)
		return
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		c.AbortWithStatus(404)
		return
	}

	c.File(resolvedPath)
}
