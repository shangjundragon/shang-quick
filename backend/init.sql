-- 建表语句

CREATE TABLE IF NOT EXISTS sys_user (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    nickname TEXT,
    phone TEXT,
    email TEXT,
    avatar TEXT,
    dept_id INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    del_flag TEXT DEFAULT 'N',
    create_by INTEGER,
    create_time INTEGER,
    update_by INTEGER,
    update_time INTEGER
);

CREATE TABLE IF NOT EXISTS sys_dept (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER DEFAULT 0,
    dept_name TEXT NOT NULL,
    order_num INTEGER DEFAULT 0,
    leader TEXT,
    phone TEXT,
    email TEXT,
    status INTEGER DEFAULT 1,
    del_flag TEXT DEFAULT 'N',
    create_time INTEGER,
    update_time INTEGER
);

CREATE TABLE IF NOT EXISTS sys_menu (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER DEFAULT 0,
    menu_name TEXT NOT NULL,
    menu_type INTEGER,
    icon TEXT,
    path TEXT,
    component TEXT,
    perm TEXT,
    order_num INTEGER DEFAULT 0,
    is_frame INTEGER DEFAULT 0,
    is_cache INTEGER DEFAULT 0,
    is_visible INTEGER DEFAULT 1,
    status INTEGER DEFAULT 1,
    del_flag TEXT DEFAULT 'N',
    create_time INTEGER,
    update_time INTEGER
);

CREATE TABLE IF NOT EXISTS sys_role (
    id INTEGER PRIMARY KEY,
    role_name TEXT NOT NULL,
    role_code TEXT NOT NULL,
    remark TEXT,
    status INTEGER DEFAULT 1,
    del_flag TEXT DEFAULT 'N',
    create_time INTEGER,
    update_time INTEGER
);

CREATE TABLE IF NOT EXISTS sys_user_role (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    role_id INTEGER
);

CREATE TABLE IF NOT EXISTS sys_role_menu (
    id INTEGER PRIMARY KEY,
    role_id INTEGER,
    menu_id INTEGER
);

CREATE TABLE IF NOT EXISTS sys_oper_log (
    id INTEGER PRIMARY KEY,
    title TEXT,
    oper_type INTEGER,
    method TEXT,
    request_url TEXT,
    request_data TEXT,
    response_data TEXT,
    oper_name TEXT,
    oper_ip TEXT,
    oper_time INTEGER,
    status INTEGER,
    error_msg TEXT
);

-- 初始化数据

INSERT INTO sys_dept (parent_id, dept_name, order_num, status, del_flag) VALUES (0, '总公司', 0, 1, 'N');

INSERT INTO sys_role (role_name, role_code, status, del_flag) VALUES ('超级管理员', 'admin', 1, 'N');

-- admin123 的 bcrypt 加密密码
INSERT INTO sys_user (username, password, nickname, dept_id, status, del_flag) VALUES ('admin', '$2a$10$N.ZUPZ7GXYX7YX7GXYX7XuXYX7GXYX7GXYX7GXYX7GXYX7GXYX7Ge', '管理员', 1, 1, 'N');

INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1);

-- 菜单数据
INSERT INTO sys_menu (parent_id, menu_name, menu_type, icon, path, component, perm, order_num, status, del_flag) VALUES 
(0, '系统管理', 0, 'SettingsOutline', '/system', '', '', 1, 1, 'N'),
(0, '仪表盘', 1, 'SpeedometerOutline', '/dashboard', 'dashboard/index', 'dashboard:list', 0, 1, 'N'),
(0, '个人中心', 1, 'PersonOutline', '/profile', 'profile/index', 'profile:list', 99, 1, 'N'),
(1, '用户管理', 1, 'PeopleOutline', 'user', 'system/user/index', 'user:list', 1, 1, 'N'),
(1, '部门管理', 1, 'BusinessOutline', 'dept', 'system/dept/index', 'dept:list', 2, 1, 'N'),
(1, '菜单管理', 1, 'MenuOutline', 'menu', 'system/menu/index', 'menu:list', 3, 1, 'N'),
(1, '角色管理', 1, 'ShieldOutline', 'role', 'system/role/index', 'role:list', 4, 1, 'N'),
(1, '操作日志', 1, 'DocumentTextOutline', 'operLog', 'system/operLog/index', 'operLog:list', 5, 1, 'N');

-- 按钮权限
INSERT INTO sys_menu (parent_id, menu_name, menu_type, perm, order_num, status, del_flag) VALUES 
(4, '新增', 2, 'user:add', 1, 1, 'N'),
(4, '编辑', 2, 'user:edit', 2, 1, 'N'),
(4, '删除', 2, 'user:delete', 3, 1, 'N'),
(4, '重置密码', 2, 'user:resetPwd', 4, 1, 'N'),
(5, '新增', 2, 'dept:add', 1, 1, 'N'),
(5, '编辑', 2, 'dept:edit', 2, 1, 'N'),
(5, '删除', 2, 'dept:delete', 3, 1, 'N'),
(6, '新增', 2, 'menu:add', 1, 1, 'N'),
(6, '编辑', 2, 'menu:edit', 2, 1, 'N'),
(6, '删除', 2, 'menu:delete', 3, 1, 'N'),
(7, '新增', 2, 'role:add', 1, 1, 'N'),
(7, '编辑', 2, 'role:edit', 2, 1, 'N'),
(7, '删除', 2, 'role:delete', 3, 1, 'N'),
(7, '分配权限', 2, 'role:assign', 4, 1, 'N');

-- 角色菜单关联（超级管理员拥有所有权限）
INSERT INTO sys_role_menu (role_id, menu_id) VALUES 
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8),
(1, 9), (1, 10), (1, 11), (1, 12), (1, 13), (1, 14), (1, 15),
(1, 16), (1, 17), (1, 18), (1, 19), (1, 20), (1, 21), (1, 22);
