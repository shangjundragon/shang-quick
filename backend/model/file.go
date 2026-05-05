package model

type SysFile struct {
	Id           int64   `dbw:"primaryKey" json:"id,string"`
	OriginalName string  `dbw:"column:original_name" json:"originalName"`
	FileName     string  `dbw:"column:file_name" json:"fileName"`
	FilePath     string  `dbw:"column:file_path" json:"filePath"`
	FileSize     int64   `dbw:"column:file_size" json:"fileSize"`
	FileType     string  `dbw:"column:file_type" json:"fileType"`
	IsImage      int     `dbw:"column:is_image;default:0" json:"isImage"`
	CreateBy     int64   `dbw:"column:create_by" json:"createBy,string"`
	CreateTime   int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"`
	UpdateTime   int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"`
	DelFlag      string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`
}

func (SysFile) TableName() string {
	return "sys_file"
}
