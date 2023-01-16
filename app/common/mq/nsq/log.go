package nsq

import (
	"go.uber.org/zap"
)

type logger struct {
	zap *zap.SugaredLogger
}

func (l *logger) Output(_ int, s string) error {
	l.zap.Debug(s)
	return nil
}

func NewLogger(zap *zap.SugaredLogger) *logger {
	return &logger{zap: zap}
}
