package consts

import "errors"

// error 常量定义
var (
	ErrEmptyConfigClient = errors.New("configClient is null (try to initialize a new one)")
	ErrGetViper          = errors.New("get viper failed")
	ErrViperEmptyKey     = errors.New("get viper key failed")
)
