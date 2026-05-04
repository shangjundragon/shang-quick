package model

type SysUserRole struct {
	Id     int64 `dbw:"primaryKey" json:"id"`
	UserId int64 `dbw:"column:user_id" json:"userId"`
	RoleId int64 `dbw:"column:role_id" json:"roleId"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
