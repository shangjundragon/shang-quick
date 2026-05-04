package model

type SysRole struct {
	Id         int64   `dbw:"primaryKey"`
	RoleName   string  `dbw:"column:role_name"`
	RoleCode   string  `dbw:"column:role_code"`
	Remark     *string `dbw:"column:remark"`
	Status     int     `dbw:"column:status;default:1"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
