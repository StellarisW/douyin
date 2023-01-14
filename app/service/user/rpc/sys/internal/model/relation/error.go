package relation

import "douyin/app/common/errx"

var (
	// 通用错误

	errRedisAdd   = errx.New(ErrIdRedisAdd, errx.RedisAdd)
	errRedisRem   = errx.New(ErrIdRedisRem, errx.RedisRem)
	errRedisRange = errx.New(ErrIdRedisRange, errx.RedisRange)
	errRedisInter = errx.New(ErrIdRedisInter, errx.RedisInter)
	errRedisGet   = errx.New(ErrIdRedisGet, errx.RedisGet)
	errMysqlGet   = errx.New(ErrIdMysqlGet, errx.MysqlGet)
)

var (
	// Relation 操作

	errInvalidActionType = errx.New(ErrIdInvalidActionType, ErrInvalidActionType)
)
