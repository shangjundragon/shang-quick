package model

// SysOperLog 操作日志，由 OperLog 中间件自动记录 POST 写操作
type SysOperLog struct {
	Id           int64   `dbw:"primaryKey" json:"id,string"`                 // 日志 ID
	Title        *string `dbw:"column:title" json:"title"`                   // 操作标题（如"用户管理"）
	OperType     int     `dbw:"column:oper_type" json:"operType"`           // 操作类型：1=新增 2=修改 3=删除 4=查询
	Method       *string `dbw:"column:method" json:"method"`                 // HTTP 方法
	RequestUrl   *string `dbw:"column:request_url" json:"requestUrl"`       // 请求 URL
	RequestData  *string `dbw:"column:request_data" json:"requestData"`     // 请求数据（JSON 字符串，密码已脱敏）
	ResponseData *string `dbw:"column:response_data" json:"responseData"`   // 响应数据（JSON 字符串，密码已脱敏）
	OperName     *string `dbw:"column:oper_name" json:"operName"`           // 操作人用户名
	OperIp       *string `dbw:"column:oper_ip" json:"operIp"`               // 操作人 IP
	OperTime     int64   `dbw:"column:oper_time" json:"operTime,string"`    // 操作时间（毫秒时间戳）
	Status       int     `dbw:"column:status" json:"status"`                // 状态：1=成功 0=失败
	ErrorMsg     *string `dbw:"column:error_msg" json:"errorMsg"`           // 错误信息
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}
