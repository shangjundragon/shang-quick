package logger

import (
	"backend/pkg/global_vars"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {

	ZapLogHandler := func(entry zapcore.Entry) error {

		// 参数 entry 介绍
		// entry  参数就是单条日志结构体，主要包括字段如下：
		//Level      日志等级
		//Time       当前时间
		//LoggerName  日志名称
		//Message    日志内容
		//Caller     各个文件调用路径
		//Stack      代码调用栈

		//这里启动一个协程，hook丝毫不会影响程序性能，
		go func(paramEntry zapcore.Entry) {
			//fmt.Println(" GoSkeleton  hook ....，你可以在这里继续处理系统日志....")
			//fmt.Printf("%#+v\n", paramEntry)
		}(entry)
		return nil
	}

	CreateZapFactory := func(entry func(zapcore.Entry) error) *zap.Logger {

		// 获取程序所处的模式：  开发调试 、 生产
		//global_vars.ConfigYml := yml_config.CreateYamlFactory()
		appDebug := global_vars.ConfigYml.GetBool("AppDebug")

		// 判断程序当前所处的模式，调试模式直接返回一个便捷的zap日志管理器地址，所有的日志打印到控制台即可
		if appDebug == true {
			if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
				return logger
			} else {
				log.Fatal("创建zap日志包失败，详情：" + err.Error())
			}
		}

		// 以下才是 非调试（生产）模式所需要的代码
		encoderConfig := zap.NewProductionEncoderConfig()

		timePrecision := global_vars.ConfigYml.GetString("Logs.TimePrecision")
		var recordTimeFormat string
		switch timePrecision {
		case "second":
			// ISO 8601 秒级：2026-03-25T12:05:47+08:00
			recordTimeFormat = "2006-01-02T15:04:05-07:00"
		case "millisecond":
			// ISO 8601 毫秒级：2026-03-25T12:05:47.123+08:00
			// 注意：Go 的参考时间中，毫秒用 .000 表示，时区偏移用 -07:00 表示（带冒号）
			recordTimeFormat = "2006-01-02T15:04:05.000-07:00"
		default:
			recordTimeFormat = "2006-01-02T15:04:05-07:00"
		}
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(recordTimeFormat))
		}
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.TimeKey = "created_at" // 生成json格式日志的时间键字段，默认为 ts,修改以后方便日志导入到 ELK 服务器

		var encoder zapcore.Encoder
		switch global_vars.ConfigYml.GetString("Logs.TextFormat") {
		case "console":
			encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
		case "json":
			encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
		default:
			encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式
		}

		//写入器
		fileName := global_vars.BasePath + global_vars.ConfigYml.GetString("Logs.AppFileLogName")
		lumberJackLogger := &lumberjack.Logger{
			Filename:   fileName,                                        //日志文件的位置
			MaxSize:    global_vars.ConfigYml.GetInt("Logs.MaxSize"),    //在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: global_vars.ConfigYml.GetInt("Logs.MaxBackups"), //保留旧文件的最大个数
			MaxAge:     global_vars.ConfigYml.GetInt("Logs.MaxAge"),     //保留旧文件的最大天数
			Compress:   global_vars.ConfigYml.GetBool("Logs.Compress"),  //是否压缩/归档旧文件
		}
		writer := zapcore.AddSync(lumberJackLogger)
		// 开始初始化zap日志核心参数，
		//参数一：编码器
		//参数二：写入器
		//参数三：参数级别，debug级别支持后续调用的所有函数写日志，如果是 fatal 高级别，则级别>=fatal 才可以写日志
		zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
		return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry), zap.AddStacktrace(zap.WarnLevel))
	}
	global_vars.ZapLog = CreateZapFactory(ZapLogHandler)
}

func GetTraceLogger() (traceLogger *zap.Logger, traceID string) {
	traceID = uuid.New().String()
	traceLogger = global_vars.ZapLog.With(
		zap.String("trace_id", traceID),
	)
	return traceLogger, traceID
}
