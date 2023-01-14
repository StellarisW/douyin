package relation

const (
	// OprId

	ErrIdOprRelation = iota
	ErrIdOprGetFollowList
	ErrIdOprGetFollowerList
	ErrIdOprGetFriendList
)

// 通用操作

const (
	ErrIdParseInt = iota
	ErrIdCommon
)

const (
	ErrParseInt = "parse int failed"
)

// Relation 操作

const (
	ErrIdInvalidActionType = iota + ErrIdCommon
)

const (
	ErrInvalidActionType = "invalid action type"
)
