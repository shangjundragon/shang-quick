package model

type SysOperLog struct {
	Id           int64   `dbw:"primaryKey" json:"id"`
	Title        *string `dbw:"column:title" json:"title"`
	OperType     int     `dbw:"column:oper_type" json:"operType"`
	Method       *string `dbw:"column:method" json:"method"`
	RequestUrl   *string `dbw:"column:request_url" json:"requestUrl"`
	RequestData  *string `dbw:"column:request_data" json:"requestData"`
	ResponseData *string `dbw:"column:response_data" json:"responseData"`
	OperName     *string `dbw:"column:oper_name" json:"operName"`
	OperIp       *string `dbw:"column:oper_ip" json:"operIp"`
	OperTime     int64   `dbw:"column:oper_time" json:"operTime"`
	Status       int     `dbw:"column:status" json:"status"`
	ErrorMsg     *string `dbw:"column:error_msg" json:"errorMsg"`
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}
