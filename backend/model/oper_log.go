package model

type SysOperLog struct {
	Id           int64   `dbw:"primaryKey"`
	Title        *string `dbw:"column:title"`
	OperType     int     `dbw:"column:oper_type"`
	Method       *string `dbw:"column:method"`
	RequestUrl   *string `dbw:"column:request_url"`
	RequestData  *string `dbw:"column:request_data"`
	ResponseData *string `dbw:"column:response_data"`
	OperName     *string `dbw:"column:oper_name"`
	OperIp       *string `dbw:"column:oper_ip"`
	OperTime     int64   `dbw:"column:oper_time"`
	Status       int     `dbw:"column:status"`
	ErrorMsg     *string `dbw:"column:error_msg"`
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}
