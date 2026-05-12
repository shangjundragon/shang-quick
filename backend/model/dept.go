package model

// SysDept 部门，树形结构通过 parent_id 关联
type SysDept struct {
	Id       int64   `dbw:"primaryKey" json:"id,string"`                                                  // 部门 ID
	ParentId int64   `dbw:"column:parent_id;default:0;tableUpdateStrategy:always" json:"parentId,string"` // 父部门 ID（0 为根）
	DeptName string  `dbw:"column:dept_name" json:"deptName"`                                             // 部门名称
	OrderNum int     `dbw:"column:order_num;default:0;tableUpdateStrategy:always" json:"orderNum"`        // 排序号
	Leader   *string `dbw:"column:leader" json:"leader"`                                                  // 负责人
	Phone    *string `dbw:"column:phone" json:"phone"`                                                    // 联系电话
	Email    *string `dbw:"column:email" json:"email"`                                                    // 邮箱
	Status   int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"`             // 状态：1=启用 0=禁用

	DelFlag    string `dbw:"column:del_flag;tableLogic" json:"delFlag"`                        // 逻辑删除标识
	CreateTime int64  `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"` // 创建时间
	UpdateTime int64  `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"` // 更新时间
}

func (SysDept) TableName() string {
	return "sys_dept"
}
