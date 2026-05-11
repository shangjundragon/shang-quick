package model

type SysDept struct {
	Id         int64   `dbw:"primaryKey" json:"id,string"`
	ParentId   int64   `dbw:"column:parent_id;default:0;tableUpdateStrategy:always" json:"parentId,string"`
	DeptName   string  `dbw:"column:dept_name" json:"deptName"`
	OrderNum   int     `dbw:"column:order_num;default:0;tableUpdateStrategy:always" json:"orderNum"`
	Leader     *string `dbw:"column:leader" json:"leader"`
	Phone      *string `dbw:"column:phone" json:"phone"`
	Email      *string `dbw:"column:email" json:"email"`
	Status     int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}
