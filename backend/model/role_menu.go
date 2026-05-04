package model

type SysRoleMenu struct {
	Id     int64 `dbw:"primaryKey" json:"id"`
	RoleId int64 `dbw:"column:role_id" json:"roleId"`
	MenuId int64 `dbw:"column:menu_id" json:"menuId"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
