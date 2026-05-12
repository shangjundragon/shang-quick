package model

// SysFile 上传文件记录，FilePath 为相对路径，物理文件存储在 storage/uploads/
type SysFile struct {
	Id           int64   `dbw:"primaryKey" json:"id,string"`                                       // 文件 ID
	OriginalName string  `dbw:"column:original_name" json:"originalName"`                         // 原始文件名
	FileName     string  `dbw:"column:file_name" json:"fileName"`                                 // 存储文件名（UUID 重命名）
	FilePath     string  `dbw:"column:file_path" json:"filePath"`                                 // 文件相对路径（如 /2026/05/uuid.jpg）
	FileSize     int64   `dbw:"column:file_size" json:"fileSize"`                                 // 文件大小（字节）
	FileType     string  `dbw:"column:file_type" json:"fileType"`                                 // 文件 MIME 类型
	IsImage      int     `dbw:"column:is_image;default:0" json:"isImage"`                        // 是否图片：0=否 1=是
	CreateBy     int64   `dbw:"column:create_by" json:"createBy,string"`                          // 上传人 ID
	CreateTime   int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"` // 上传时间
	UpdateTime   int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"` // 更新时间
	DelFlag      string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`                        // 逻辑删除标识
}

func (SysFile) TableName() string {
	return "sys_file"
}
