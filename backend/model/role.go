package model

// SysRole 系统角色，role_code 唯一标识（admin 为超级管理员）
type SysRole struct {
	Id       int64   `dbw:"primaryKey" json:"id,string"`                                      // 角色 ID
	RoleName string  `dbw:"column:role_name" json:"roleName"`                                 // 角色名称
	RoleCode string  `dbw:"column:role_code" json:"roleCode"`                                 // 角色编码（唯一）
	Remark   *string `dbw:"column:remark" json:"remark"`                                      // 备注
	Status   int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"` // 状态：1=启用 0=禁用

	DelFlag    string `dbw:"column:del_flag;tableLogic" json:"delFlag"`                        // 逻辑删除标识：N=未删 Y=已删
	CreateBy   int64  `dbw:"column:create_by;createBy" json:"createBy,string"`                 // 创建人 ID
	CreateTime int64  `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"` // 创建时间（毫秒时间戳）
	UpdateBy   int64  `dbw:"column:update_by;updateBy" json:"updateBy,string"`                 // 更新人 ID
	UpdateTime int64  `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"` // 更新时间（毫秒时间戳）
}

func (SysRole) TableName() string {
	return "sys_role"
}
