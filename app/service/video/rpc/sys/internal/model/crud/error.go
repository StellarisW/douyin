package crud

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlInsert       = errx.New(ErrIdMysqlInsert, errx.MysqlInsert)
	errMinioPut          = errx.New(ErrIdMinioPut, errx.MinioPut)
	errRedisAdd          = errx.New(ErrIdRedisAdd, errx.RedisAdd)
	errRedisIncr         = errx.New(ErrIdRedisIncr, errx.RedisIncr)
	errRedisRem          = errx.New(ErrIdRedisRem, errx.RedisRem)
	errRedisDecr         = errx.New(ErrIdRedisDecr, errx.RedisDecr)
	errInvalidActionType = errx.New(ErrIdInvalidActionType, ErrInvalidActionType)
	errMysqlDelete       = errx.New(ErrIdMysqlDelete, errx.MysqlDelete)
)

var (
	// Publish 操作

	errInvalidContentType = errx.New(ErrIdInvalidContentType, ErrInvalidContentType)
)
