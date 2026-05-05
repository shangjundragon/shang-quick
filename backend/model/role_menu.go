package model

type SysRoleMenu struct {
	Id     int64 `dbw:"primaryKey" json:"id,string"`
	RoleId int64 `dbw:"column:role_id" json:"roleId,string"`
	MenuId int64 `dbw:"column:menu_id" json:"menuId,string"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
