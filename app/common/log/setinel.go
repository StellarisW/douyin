package log

import (
	"douyin/app/common/douyin"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alibaba/sentinel-golang/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type SentinelLogger struct {
	logger *zap.Logger
}

func GetSentinelLogger() (logging.Logger, error) {
	if Logger == nil {
		return nil, errors.New("logger is null, try user InitLogger to initialize a logger")
	}

	return &SentinelLogger{
		logger: Logger,
	}, nil
}

func (l *SentinelLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, getFields(keysAndValues)...)
}

func (l *SentinelLogger) DebugEnabled() bool {
	return zap.DebugLevel >= getLoggerLevel()
}

func (l *SentinelLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, getFields(keysAndValues)...)
}

func (l *SentinelLogger) InfoEnabled() bool {
	return zap.InfoLevel >= getLoggerLevel()
}

func (l *SentinelLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, getFields(keysAndValues)...)
}

func (l *SentinelLogger) WarnEnabled() bool {
	return zap.WarnLevel >= getLoggerLevel()
}

func (l *SentinelLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, append(getFields(keysAndValues), zap.Error(err))...)
}

func (l *SentinelLogger) ErrorEnabled() bool {
	return zap.ErrorLevel >= getLoggerLevel()
}

func getFields(keysAndValues ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues))

	size := len(keysAndValues)

	if size&1 != 0 {
		bytes, _ := json.Marshal(fmt.Sprintf("%+v", keysAndValues))
		fields = append(fields, zap.String("param", string(bytes)))

	} else if size != 0 {
		for i := 0; i < size; i += 2 {
			fields = append(fields, zap.Reflect(fmt.Sprintf("param[%d]", i>>1), keysAndValues[i+1]))
		}
	}

	return fields
}

func getLoggerLevel() zapcore.Level {
	mode := os.Getenv(douyin.ModeEnvName)

	switch mode {
	case douyin.DevMode:
		return zap.DebugLevel

	case douyin.FatMode, douyin.UatMode:
		return zap.InfoLevel

	case douyin.ProMode:
		return zap.ErrorLevel

	default:
		return zap.DebugLevel
	}
}
