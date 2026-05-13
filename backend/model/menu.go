package model

// SysMenu 系统菜单，menu_type: 0=目录 1=菜单 2=按钮，perm 为权限标识
type SysMenu struct {
	Id        int64   `dbw:"primaryKey" json:"id,string"`                                                  // 菜单 ID
	ParentId  int64   `dbw:"column:parent_id;default:0;tableUpdateStrategy:always" json:"parentId,string"` // 父菜单 ID（0 为根）
	MenuName  string  `dbw:"column:menu_name" json:"menuName"`                                             // 菜单名称
	MenuType  int     `dbw:"column:menu_type;default:0;tableUpdateStrategy:always" json:"menuType"`        // 类型：0=目录 1=菜单 2=按钮
	Icon      *string `dbw:"column:icon" json:"icon"`                                                      // 图标（ionicons 名称）
	Path      *string `dbw:"column:path" json:"path"`                                                      // 路由路径
	Component *string `dbw:"column:component" json:"component"`                                            // 组件路径（相对 views/）
	Perm      *string `dbw:"column:perm" json:"perm"`                                                      // 权限标识（如 user:add）
	OrderNum  int     `dbw:"column:order_num;default:0;tableUpdateStrategy:always" json:"orderNum"`        // 排序号
	IsFrame   int     `dbw:"column:is_frame;default:0;tableUpdateStrategy:always" json:"isFrame"`          // 是否外链：0=否 1=是
	IsCache   int     `dbw:"column:is_cache;default:0;tableUpdateStrategy:always" json:"isCache"`          // 是否缓存：0=否 1=是
	IsVisible int     `dbw:"column:is_visible;default:1;tableUpdateStrategy:always" json:"isVisible"`      // 是否显示：0=隐藏 1=显示
	Status    int     `dbw:"column:status;default:1;tableUpdateStrategy:always" json:"status"`             // 状态：1=启用 0=禁用

	DelFlag    string `dbw:"column:del_flag;tableLogic" json:"delFlag"`                        // 逻辑删除标识：N=未删 Y=已删
	CreateBy   int64  `dbw:"column:create_by;createBy" json:"createBy,string"`                 // 创建人 ID
	CreateTime int64  `dbw:"column:create_time;autoCreateTime:milli" json:"createTime,string"` // 创建时间（毫秒时间戳）
	UpdateBy   int64  `dbw:"column:update_by;updateBy" json:"updateBy,string"`                 // 更新人 ID
	UpdateTime int64  `dbw:"column:update_time;autoUpdateTime:milli" json:"updateTime,string"` // 更新时间（毫秒时间戳）
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
