package bootstrap

import (
	"backend/pkg/cache"
	"backend/pkg/constants"
	"backend/pkg/global_vars"
	"backend/pkg/logger"
	"backend/pkg/req_util"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/shangjundragon/dbw"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Bootstrap() {
	// 1. 初始化 项目根路径，
	initBasePath()
	// 2.检查配置文件以及日志目录等非编译性的必要条件
	checkRequiredFolders()
	// 3.加载配置文件
	loadConfigFile()
	// 4.初始化日志
	logger.InitLogger()
	// 5.初始化缓存
	cache.InitCache()
	// 6.初始化数据库
	initDatabase()
}

// 初始化项目根目录
func initBasePath() {
	curPath, err := os.Getwd()
	if err != nil {
		log.Fatal("初始化项目根目录失败")
	}

	// 判断是否在运行测试
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
		// 如果当前路径的最后一层是 tests 或 test，则使用父目录作为根路径
		if base := filepath.Base(curPath); base == "test" || base == "tests" {
			global_vars.BasePath = filepath.Dir(curPath)
		} else {
			global_vars.BasePath = curPath
		}
	} else {
		global_vars.BasePath = curPath
	}

	log.Printf("global.BasePath = %s", global_vars.BasePath)
}

// 检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录
func checkRequiredFolders() {
	log.Printf("global.BasePath %s", global_vars.BasePath)
	//1.检查配置文件是否存在
	if _, err := os.Stat(filepath.Join(global_vars.BasePath, "config", "config.yml")); err != nil {
		log.Fatalf("检查配置文件不存在: %s", err.Error())
	}
	//2.检查storage/logs 目录是否存在
	if _, err := os.Stat(filepath.Join(global_vars.BasePath, "storage", "logs")); err != nil {
		log.Fatalf("storage/logs 日志目录不存在: %s", err.Error())
	}
}

// 加载配置文件
func loadConfigFile() {
	global_vars.ConfigYml = viper.New()
	global_vars.ConfigYml.SetConfigName("config")
	global_vars.ConfigYml.SetConfigFile(filepath.Join(global_vars.BasePath, "config", "config.yml"))
	// 读取配置文件
	if err := global_vars.ConfigYml.ReadInConfig(); err != nil {
		log.Fatalf("加载配置文件不存在: %s", err.Error())
	}
}

func initDatabase() {

	dbw.SetLogFn(func(sqlStr string, args []any, ctx context.Context) {
		contextTraceId := ctx.Value(constants.ContextTraceIDKey)
		var traceLogger *zap.Logger
		if contextTraceId == nil || contextTraceId == "" {
			contextTraceId = ""
			traceLogger, _ = logger.GetTraceLogger()
		} else {
			traceLogger, _ = req_util.GetTraceLogger(ctx)
		}

		if sqlStr == "" {
			traceLogger.Warn("sqlStr is empty")
			return
		}
		if len(args) == 0 {
			traceLogger.Warn("args is empty")
		}

		// escapeSQLString 转义字符串中的特殊字符
		escapeSQLString := func(s string) string {
			s = strings.ReplaceAll(s, "'", "''")          // 单引号转义
			s = strings.ReplaceAll(s, "\\", "\\\\")       // 反斜杠转义
			s = strings.ReplaceAll(s, "\n", "\\n")        // 换行符
			s = strings.ReplaceAll(s, "\r", "\\r")        // 回车符
			s = strings.ReplaceAll(s, "\t", "\\t")        // 制表符
			return s
		}

		// formatArg 将参数格式化为SQL中可显示的形式
		formatArg := func(arg any) string {
			if arg == nil {
				return "NULL"
			}

			v := reflect.ValueOf(arg)

			// 处理指针类型
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					return "NULL"
				}
				v = v.Elem()
			}

			switch v.Kind() {
			case reflect.String:
				// 字符串需要转义并加引号
				return fmt.Sprintf("'%s'", escapeSQLString(v.String()))

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return fmt.Sprintf("%d", v.Int())

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return fmt.Sprintf("%d", v.Uint())

			case reflect.Float32, reflect.Float64:
				return fmt.Sprintf("%v", v.Float())

			case reflect.Bool:
				if v.Bool() {
					return "TRUE"
				}
				return "FALSE"

			case reflect.Struct:
				// 处理时间类型
				if t, ok := arg.(time.Time); ok {
					return fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05"))
				}
				return fmt.Sprintf("'%v'", arg)

			default:
				// 其他类型尝试直接转换
				return fmt.Sprintf("'%v'", arg)
			}
		}
		// ReplaceSQLParams 将SQL中的问号占位符替换为实际参数值，用于开发环境调试
		// 注意：此函数仅用于日志/调试，不可用于实际执行（有SQL注入风险）
		ReplaceSQLParams := func(sql string, args []any) string {
			if len(args) == 0 {
				return sql
			}

			result := sql
			for _, arg := range args {
				// 找到第一个问号的位置
				idx := strings.Index(result, "?")
				if idx == -1 {
					break
				}

				// 将参数转换为可显示的字符串
				replaced := formatArg(arg)

				// 替换第一个问号
				result = result[:idx] + replaced + result[idx+1:]
			}

			return result
		}
		traceLogger.Info("Debug SQL", zap.String(constants.ContextTraceIDKey, contextTraceId.(string)), zap.String("SQL", ReplaceSQLParams(sqlStr, args)))
	})

	driverName := global_vars.ConfigYml.GetString("DB.DriverName")
	// 打开数据库连接
	db, err := sql.Open(driverName, global_vars.ConfigYml.GetString("Db.Dsn"))
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	// 尝试连接数据库
	err = db.Ping()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	global_vars.ZapLog.Info("成功连接到数据库")
	global_vars.Db = db

	global_vars.DbConfig = dbw.NewConfig(func(config *dbw.Config) {
		config.Db = global_vars.Db
		config.Debug = global_vars.ConfigYml.GetBool("AppDebug")
		config.DriverName = driverName
		config.LogicDeleteValue = constants.Y
		config.LogicNotDeleteValue = constants.N
	})

	dbw.RegisterEntityHook(func(ctx context.Context, point dbw.HookPoint, entity any) error {
		if point != dbw.HookBeforeInsert && point != dbw.HookBeforeUpdate {
			return nil
		}
		v := reflect.ValueOf(entity).Elem()
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			tagMap := dbw.ResolveDbwTag(t.Field(i).Tag.Get("dbw"))
			if tagMap["createBy"] == "true" || tagMap["updateBy"] == "true" {
				operationUserId := ctx.Value(constants.ContextUserIDKey)
				if operationUserId == nil || operationUserId == "" {
					global_vars.ZapLog.Warn("操作人用户id为空 设置为超管id")
					operationUserId = 1
				}
				dbw.SetFieldValue(v.Field(i), operationUserId)
			}
		}
		return nil
	})

	InitDatabase()
}
