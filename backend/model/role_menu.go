package model

// SysRoleMenu 角色-菜单关联表（多对多），用于 RBAC 权限控制
type SysRoleMenu struct {
	Id     int64 `dbw:"primaryKey" json:"id,string"`                // 关联 ID
	RoleId int64 `dbw:"column:role_id" json:"roleId,string"`        // 角色 ID
	MenuId int64 `dbw:"column:menu_id" json:"menuId,string"`        // 菜单 ID
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
