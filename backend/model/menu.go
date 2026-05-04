package model

type SysMenu struct {
	Id         int64   `dbw:"primaryKey"`
	ParentId   int64   `dbw:"column:parent_id;default:0"`
	MenuName   string  `dbw:"column:menu_name"`
	MenuType   int     `dbw:"column:menu_type"`
	Icon       *string `dbw:"column:icon"`
	Path       *string `dbw:"column:path"`
	Component  *string `dbw:"column:component"`
	Perm       *string `dbw:"column:perm"`
	OrderNum   int     `dbw:"column:order_num;default:0"`
	IsFrame    int     `dbw:"column:is_frame;default:0"`
	IsCache    int     `dbw:"column:is_cache;default:0"`
	IsVisible  int     `dbw:"column:is_visible;default:1"`
	Status     int     `dbw:"column:status;default:1"`
	DelFlag    string  `dbw:"column:del_flag;tableLogic"`
	CreateTime int64   `dbw:"column:create_time;autoCreateTime:milli"`
	UpdateTime int64   `dbw:"column:update_time;autoUpdateTime:milli"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
