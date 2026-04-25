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
