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
