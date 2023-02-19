package relation

const (
	// OprId

	ErrIdOprRelation = iota
	ErrIdOprGetFollowList
	ErrIdOprGetFollowerList
	ErrIdOprGetFriendList
)

// 通用错误

const (
	ErrIdRedisAdd = iota
	ErrIdRedisRem
	ErrIdRedisRange
	ErrIdRedisInter
	ErrIdRedisGet
	ErrIdMysqlGet
	ErrIdRedisPipeExec
	ErrIdCommon
)

// Relation 操作

const (
	ErrIdInvalidActionType = iota + ErrIdCommon
	ErrIdAlreadyFollow
	ErrIdAlreadyUnfollow
	ErrIdUserNotFound
)

const (
	ErrInvalidActionType = "invalid action type"
	ErrAlreadyFollow     = "already followed"
	ErrAlreadyUnfollow   = "already unfollowed"
	ErrUserNotFound      = "user not found"
)
