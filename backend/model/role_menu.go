package model

type SysRoleMenu struct {
	Id     int64 `dbw:"primaryKey"`
	RoleId int64 `dbw:"column:role_id"`
	MenuId int64 `dbw:"column:menu_id"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
