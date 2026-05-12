// Package global_vars 定义全局变量：项目根路径、日志、配置、数据库连接和 dbw 配置
package global_vars

import (
	"database/sql"

	"github.com/shangjundragon/dbw"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	BasePath string // 定义项目的根目录

	// 全局日志指针
	ZapLog    *zap.Logger
	ConfigYml *viper.Viper
	Db        *sql.DB
	DbConfig  *dbw.Config
)
