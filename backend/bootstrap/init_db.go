package bootstrap

import (
	"backend/pkg/global_vars"
	"backend/pkg/password"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func InitDatabase() {
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

	var count int64
	err = global_vars.Db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sys_user'").Scan(&count)
	
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

	var userCount int64
	err = global_vars.Db.QueryRow("SELECT COUNT(*) FROM sys_user").Scan(&userCount)
	if err == nil && userCount > 0 {
		updateAdminPassword()
	}

	if err == nil {
		log.Println("数据库初始化完成")
	}
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

func updateAdminPassword() {
	hashedPwd, err := password.Hash("admin123")
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return
	}
	
	_, err = global_vars.Db.Exec(
		"UPDATE sys_user SET password = ? WHERE username = ?",
		hashedPwd, "admin",
	)
	if err != nil {
		log.Printf("更新 admin 密码失败: %v", err)
	}
}
