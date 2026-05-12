package model

// SysUserRole 用户-角色关联表（多对多）
type SysUserRole struct {
	Id     int64 `dbw:"primaryKey" json:"id,string"`                // 关联 ID
	UserId int64 `dbw:"column:user_id" json:"userId,string"`        // 用户 ID
	RoleId int64 `dbw:"column:role_id" json:"roleId,string"`        // 角色 ID
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
