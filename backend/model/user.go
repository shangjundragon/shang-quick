package model

type SysUser struct {
	Id         int64   `dbw:"primaryKey"`
	Username   string  `dbw:"column:username"`
	Password   string  `dbw:"column:password"`
	Nickname   *string `dbw:"column:nickname"`
	Phone      *string `dbw:"column:phone"`
	Email      *string `dbw:"column:email"`
	Avatar     *string `dbw:"column:avatar"`
	DeptId     int64   `dbw:"column:dept_id"`
	Status     int     `dbw:"column:status;default:1"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic"`
	CreateBy   int64   `dbw:"column:create_by"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli"`
	UpdateBy   int64   `dbw:"column:update_by"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
