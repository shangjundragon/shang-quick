// Package model 定义数据模型结构体，映射数据库表，配合 dbw ORM 使用
package model

// SysUser 系统用户，密码用 bcrypt 哈希存储，逻辑删除使用 del_flag
type SysUser struct {
	Id       int64   `dbw:"primaryKey" json:"id,string"`                                      // 用户 ID
	Username string  `dbw:"column:username" json:"username"`                                  // 用户名
	Password string  `dbw:"column:password" json:"-"`                                         // 密码（bcrypt 哈希，不返回前端）
	Nickname *string `dbw:"column:nickname" json:"nickname"`                                  // 昵称
	Phone    *string `dbw:"column:phone" json:"phone"`                                        // 手机号
	Email    *string `dbw:"column:email" json:"email"`                                        // 邮箱
	Avatar   *string `dbw:"column:avatar" json:"avatar"`                                      // 头像 URL
	DeptId   int64   `dbw:"column:dept_id;tableUpdateStrategy:always" json:"deptId,string"`   // 部门 ID
	Status   int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"` // 状态：1=启用 0=禁用

	DelFlag    string `dbw:"column:del_flag;tableLogic" json:"delFlag"`                        // 逻辑删除标识：N=未删 Y=已删
	CreateBy   int64  `dbw:"column:create_by;createBy" json:"createBy,string"`                 // 创建人 ID
	CreateTime int64  `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"` // 创建时间（毫秒时间戳）
	UpdateBy   int64  `dbw:"column:update_by;updateBy" json:"updateBy,string"`                 // 更新人 ID
	UpdateTime int64  `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"` // 更新时间（毫秒时间戳）
}

func (SysUser) TableName() string {
	return "sys_user"
}
