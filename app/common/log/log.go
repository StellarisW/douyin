package log

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/utils/file"
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func InitLogger(name string) error {
	options := Options{
		Mode:         os.Getenv(douyin.ModeEnvName),
		SavePath:     name + "/log",
		EncoderType:  ConsoleEncoder,
		EncodeLevel:  CapitalLevelEncoder,
		EncodeCaller: FullCallerEncoder,
	}
	return NewLoggerWithOptions(options)
}

func NewLoggerWithOptions(options Options) (err error) {
	// 创建日志保存的文件夹
	err = file.IsNotExistMkDir(options.SavePath)
	if err != nil {
		return err
	}

	dynamicLevel := zap.NewAtomicLevel()
	encoder := getEncoder(options)

	switch options.Mode {
	case douyin.DevMode:
		// 将当前日志等级设置为 Debug
		// 注意日志等级低于设置的等级，日志文件也不分记录
		dynamicLevel.SetLevel(zap.DebugLevel)

		// 调试级别
		debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.DebugLevel
		})
		// 日志级别
		infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.InfoLevel
		})
		// 警告级别
		warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.WarnLevel
		})
		// 错误级别(包含error,panic,fatal级别的日志)
		errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		})
		cores := [...]zapcore.Core{
			zapcore.NewCore(encoder, os.Stdout, dynamicLevel), //控制台输出
			//日志文件输出,按等级归档
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/all/server_all.log", options.SavePath)), zapcore.DebugLevel),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/debug/server_debug.log", options.SavePath)), debugPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/info/server_info.log", options.SavePath)), infoPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/warn/server_warn.log", options.SavePath)), warnPriority),
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/error/server_error.log", options.SavePath)), errorPriority),
		}
		Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(Logger)
	case douyin.FatMode, douyin.UatMode:
		// 将当前日志等级设置为 Info
		// 注意日志等级低于设置的等级，日志文件也不分记录
		dynamicLevel.SetLevel(zap.InfoLevel)

		cores := [...]zapcore.Core{
			zapcore.NewCore(encoder, os.Stdout, dynamicLevel), //控制台输出
			zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/server.log", options.SavePath)), dynamicLevel),
		}

		Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(Logger)
	case douyin.ProMode:
		dynamicLevel.SetLevel(zap.ErrorLevel)

		Logger = zap.New(zapcore.NewCore(
			encoder,
			getWriteSyncer(fmt.Sprintf("./%s/server.log", options.SavePath)), dynamicLevel),
			zap.AddCaller(),
		)
		defer func(zapLogger *zap.Logger) {
			_ = zapLogger.Sync()
		}(Logger)
	default:
		panic(errx.EmptyProjectMode)
	}

	//设置全局logger
	Logger.Info("initialize logger successfully")
	return nil
}

// getEncoder 获取编码器
func getEncoder(options Options) zapcore.Encoder {
	if options.EncoderType == JsonEncoder {
		return zapcore.NewJSONEncoder(getEncoderConfig(options))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(options))
}

// getEncoderConfig 编码器设置
func getEncoderConfig(options Options) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",                     // 日志消息键
		LevelKey:       "level",                       // 日志等级键
		TimeKey:        "time",                        // 时间键
		NameKey:        "logger",                      // 日志记录器名
		CallerKey:      "caller",                      // 日志文件信息键
		StacktraceKey:  "stacktrace",                  // 堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 友好日志等级名大小写（info INFO）
		EncodeTime:     CustomTimeEncoder,             // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.FullCallerEncoder,     // 日志文件信息 short（包/文件.go:行号） full (文件位置.go:行号)
	}
	switch {
	case options.EncodeLevel == LowercaseLevelEncoder: // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case options.EncodeLevel == LowercaseColorLevelEncoder: // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case options.EncodeLevel == CapitalLevelEncoder: // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case options.EncodeLevel == CapitalColorLevelEncoder: // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	if options.EncodeCaller == ShortCallerEncoder {
		config.EncodeCaller = zapcore.ShortCallerEncoder
	}
	return config
}

// getWriteSyncer 读写器设置
func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,  // 日志文件的位置
		MaxSize:    10,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 50000, // 保留旧文件的最大个数
		MaxAge:     1000,  // 保留旧文件的最大天数
		Compress:   true,  // 是否压缩/归档旧文件
		LocalTime:  true,  // 是否使用本地时间
	}
	return zapcore.AddSync(lumberJackLogger)
}

// CustomTimeEncoder 格式化时间
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}
