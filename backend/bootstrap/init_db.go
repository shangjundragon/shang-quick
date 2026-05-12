package bootstrap

import (
	"backend/pkg/global_vars"
	"backend/pkg/password"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InitDatabase 执行建表 SQL，检查 admin 用户是否存在，不存在则插入初始化数据
func InitDatabase() {
	executeSQLFile()

	var adminCount int64
	err := global_vars.Db.QueryRow("SELECT COUNT(*) FROM sys_user WHERE username = 'admin'").Scan(&adminCount)
	if err == nil && adminCount > 0 {
		log.Println("admin 用户已存在，跳过初始化数据插入")
		return
	}

	insertInitData()
	log.Println("数据库初始化数据插入完成")
}

func executeSQLFile() {
	initSQLPath := filepath.Join(global_vars.BasePath, "init.sql")

	if _, err := os.Stat(initSQLPath); os.IsNotExist(err) {
		log.Printf("init.sql 文件不存在: %s", initSQLPath)
		return
	}

	sqlBytes, err := os.ReadFile(initSQLPath)
	if err != nil {
		log.Printf("读取 init.sql 失败: %v", err)
		return
	}

	sqlContent := string(sqlBytes)
	statements := splitSQLStatements(sqlContent)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := global_vars.Db.Exec(stmt)
		if err != nil {
			log.Printf("执行 SQL 失败: %v\nSQL: %s", err, stmt)
		}
	}

	log.Println("数据库表结构初始化完成")
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	lines := strings.Split(sql, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	return statements
}

func insertInitData() {
	now := time.Now().UTC().UnixMilli()

	_, err := global_vars.Db.Exec(
		`INSERT INTO sys_dept (parent_id, dept_name, order_num, status, del_flag, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		0, "总公司", 0, 1, "N", now, now)
	if err != nil {
		log.Fatalf("初始化部门失败: %v", err)
	}

	_, err = global_vars.Db.Exec(
		`INSERT INTO sys_role (role_name, role_code, status, del_flag, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)`,
		"超级管理员", "admin", 1, "N", now, now)
	if err != nil {
		log.Fatalf("初始化角色失败: %v", err)
	}

	hashedPwd, err := password.Hash("admin123")
	if err != nil {
		log.Fatalf("初始化密码失败: %v", err)
	}
	_, err = global_vars.Db.Exec(
		`INSERT INTO sys_user (id ,username, password, nickname, dept_id, status, del_flag, create_by, create_time, update_by, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		1, "admin", hashedPwd, "管理员", 1, 1, "N", 1, now, 1, now)
	if err != nil {
		log.Fatalf("初始化用户失败: %v", err)
	}

	_, err = global_vars.Db.Exec(`INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1)`)
	if err != nil {
		log.Fatalf("初始化用户角色关联失败: %v", err)
	}

	// 内置菜单定义：目录(menuType=0)、菜单(menuType=1)、按钮(menuType=2)
	menus := []struct {
		parentId  int64
		menuName  string
		menuType  int
		icon      string
		path      string
		component string
		perm      string
		orderNum  int
	}{
		{0, "系统管理", 0, "SettingsOutline", "/system", "", "", 1},
		{0, "仪表盘", 1, "SpeedometerOutline", "/dashboard", "dashboard/index", "dashboard:list", 0},
		{0, "个人中心", 1, "PersonOutline", "/profile", "profile/index", "profile:list", 99},
		{1, "用户管理", 1, "PeopleOutline", "user", "system/user/index", "user:list", 1},
		{1, "部门管理", 1, "BusinessOutline", "dept", "system/dept/index", "dept:list", 2},
		{1, "菜单管理", 1, "MenuOutline", "menu", "system/menu/index", "menu:list", 3},
		{1, "角色管理", 1, "ShieldOutline", "role", "system/role/index", "role:list", 4},
		{1, "操作日志", 1, "DocumentTextOutline", "operLog", "system/operLog/index", "operLog:list", 5},
		{1, "文件管理", 1, "FolderOutline", "file", "system/file/index", "file:list", 6},
		{1, "在线用户", 1, "WifiOutline", "onlineUser", "system/onlineUser/index", "onlineUser:list", 7},
	}

	menuIds := make([]int64, len(menus))
	for i, menu := range menus {
		result, err := global_vars.Db.Exec(
			`INSERT INTO sys_menu (parent_id, menu_name, menu_type, icon, path, component, perm, order_num, status, del_flag, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			menu.parentId, menu.menuName, menu.menuType, menu.icon, menu.path, menu.component, menu.perm, menu.orderNum, 1, "N", now, now)
		if err != nil {
			log.Fatalf("初始化菜单失败: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatalf("获取菜单ID失败: %v", err)
		}
		menuIds[i] = id
	}

	// 按钮级权限（menu_type=2），按父菜单 ID 关联
	buttons := []struct {
		parentId int64
		menuName string
		perm     string
	}{
		{menuIds[3], "新增", "user:add"},
		{menuIds[3], "编辑", "user:edit"},
		{menuIds[3], "删除", "user:delete"},
		{menuIds[3], "重置密码", "user:resetPwd"},
		{menuIds[4], "新增", "dept:add"},
		{menuIds[4], "编辑", "dept:edit"},
		{menuIds[4], "删除", "dept:delete"},
		{menuIds[5], "新增", "menu:add"},
		{menuIds[5], "编辑", "menu:edit"},
		{menuIds[5], "删除", "menu:delete"},
		{menuIds[6], "新增", "role:add"},
		{menuIds[6], "编辑", "role:edit"},
		{menuIds[6], "删除", "role:delete"},
		{menuIds[6], "分配权限", "role:assign"},
		{menuIds[7], "上传", "file:upload"},
		{menuIds[7], "删除", "file:delete"},
		{menuIds[8], "踢出", "onlineUser:kick"},
	}

	for _, btn := range buttons {
		result, err := global_vars.Db.Exec(
			`INSERT INTO sys_menu (parent_id, menu_name, menu_type, perm, order_num, status, del_flag, create_time, update_time) VALUES (?, ?, 2, ?, 0, 1, 'N', ?, ?)`,
			btn.parentId, btn.menuName, btn.perm, now, now)
		if err != nil {
			log.Fatalf("初始化按钮失败: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatalf("获取按钮ID失败: %v", err)
		}
		_, err = global_vars.Db.Exec(
			`INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)`,
			1, id)
		if err != nil {
			log.Fatalf("初始化角色菜单关联失败: %v", err)
		}
	}

	for _, menuId := range menuIds {
		_, err := global_vars.Db.Exec(
			`INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)`,
			1, menuId)
		if err != nil {
			log.Fatalf("初始化角色菜单关联失败: %v", err)
		}
	}
}
