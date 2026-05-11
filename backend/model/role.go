package model

type SysRole struct {
	Id         int64   `dbw:"primaryKey" json:"id,string"`
	RoleName   string  `dbw:"column:role_name" json:"roleName"`
	RoleCode   string  `dbw:"column:role_code" json:"roleCode"`
	Remark     *string `dbw:"column:remark" json:"remark"`
	Status     int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
