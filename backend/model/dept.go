package model

type SysDept struct {
	Id         int64   `dbw:"primaryKey"`
	ParentId   int64   `dbw:"column:parent_id;default:0"`
	DeptName   string  `dbw:"column:dept_name"`
	OrderNum   int     `dbw:"column:order_num;default:0"`
	Leader     *string `dbw:"column:leader"`
	Phone      *string `dbw:"column:phone"`
	Email      *string `dbw:"column:email"`
	Status     int     `dbw:"column:status;default:1"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}
