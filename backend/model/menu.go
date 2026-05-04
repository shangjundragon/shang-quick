package model

type SysMenu struct {
	Id         int64   `dbw:"primaryKey" json:"id"`
	ParentId   int64   `dbw:"column:parent_id;default:0" json:"parentId"`
	MenuName   string  `dbw:"column:menu_name" json:"menuName"`
	MenuType   int     `dbw:"column:menu_type" json:"menuType"`
	Icon       *string `dbw:"column:icon" json:"icon"`
	Path       *string `dbw:"column:path" json:"path"`
	Component  *string `dbw:"column:component" json:"component"`
	Perm       *string `dbw:"column:perm" json:"perm"`
	OrderNum   int     `dbw:"column:order_num;default:0" json:"orderNum"`
	IsFrame    int     `dbw:"column:is_frame;default:0" json:"isFrame"`
	IsCache    int     `dbw:"column:is_cache;default:0" json:"isCache"`
	IsVisible  int     `dbw:"column:is_visible;default:1" json:"isVisible"`
	Status     int     `dbw:"column:status;default:1" json:"status"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic" json:"delFlag"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli" json:"createTime"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
