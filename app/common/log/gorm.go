package log

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormZapWriter struct {
	logger *zap.Logger
	logger.Config
}

func GetGormLoggerConfig() logger.Config {
	return logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	}
}

func GetGormZapWriter(config logger.Config) (logger.Interface, error) {
	if Logger == nil {
		return nil, errors.New("logger is null, try user InitLogger to initialize a logger")
	}

	return &GormZapWriter{
		logger: Logger,
		Config: config,
	}, nil
}

func (w *GormZapWriter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *w
	newLogger.LogLevel = level
	return &newLogger
}

func (w *GormZapWriter) Info(ctx context.Context, msg string, data ...interface{}) {
	if w.LogLevel >= logger.Info {
		w.logger.Info(msg, zap.String("line", utils.FileWithLineNum()), zap.Reflect("data", data))
	}
}

func (w *GormZapWriter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if w.LogLevel >= logger.Warn {
		w.logger.Warn(msg, zap.String("line", utils.FileWithLineNum()), zap.Reflect("data", data))
	}
}

func (w *GormZapWriter) Error(ctx context.Context, msg string, data ...interface{}) {
	if w.LogLevel >= logger.Error {
		w.logger.Error(msg, zap.String("line", utils.FileWithLineNum()), zap.Reflect("data", data))
	}
}

func (w *GormZapWriter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if w.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && w.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !w.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			w.logger.Error("",
				zap.String("line", utils.FileWithLineNum()),
				zap.Error(err),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.String("rows", "-"),
				zap.String("sql", sql))
		} else {
			w.logger.Error("",
				zap.String("line", utils.FileWithLineNum()),
				zap.Error(err),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}

	case elapsed > w.SlowThreshold && w.SlowThreshold != 0 && w.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", w.SlowThreshold)
		if rows == -1 {
			w.logger.Warn("",
				zap.String("line", utils.FileWithLineNum()),
				zap.String("slowLog", slowLog),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.String("rows", "-"),
				zap.String("sql", sql))
		} else {
			w.logger.Warn("",
				zap.String("line", utils.FileWithLineNum()),
				zap.String("slowLog", slowLog),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}

	case w.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			w.logger.Info("",
				zap.String("line", utils.FileWithLineNum()),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.String("rows", "-"),
				zap.String("sql", sql))
		} else {
			w.logger.Info("",
				zap.String("line", utils.FileWithLineNum()),
				zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	}
}
