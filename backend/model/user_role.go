package model

type SysUserRole struct {
	Id     int64 `dbw:"primaryKey"`
	UserId int64 `dbw:"column:user_id"`
	RoleId int64 `dbw:"column:role_id"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
