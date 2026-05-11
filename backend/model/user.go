package model

type SysUser struct {
	Id         int64   `dbw:"primaryKey" json:"id,string"`
	Username   string  `dbw:"column:username" json:"username"`
	Password   string  `dbw:"column:password" json:"-"`
	Nickname   *string `dbw:"column:nickname" json:"nickname"`
	Phone      *string `dbw:"column:phone" json:"phone"`
	Email      *string `dbw:"column:email" json:"email"`
	Avatar     *string `dbw:"column:avatar" json:"avatar"`
	DeptId     int64   `dbw:"column:dept_id;tableUpdateStrategy:always" json:"deptId,string"`
	Status     int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`
	CreateBy   int64   `dbw:"column:create_by;createBy" json:"createBy,string"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"`
	UpdateBy   int64   `dbw:"column:update_by;updateBy" json:"updateBy,string"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
