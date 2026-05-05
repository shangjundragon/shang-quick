package model

type SysUserRole struct {
	Id     int64 `dbw:"primaryKey" json:"id,string"`
	UserId int64 `dbw:"column:user_id" json:"userId,string"`
	RoleId int64 `dbw:"column:role_id" json:"roleId,string"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
