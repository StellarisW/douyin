package user

import "douyin/app/service/user/internal/sys"

const (
	RdbKey = sys.SysName + ":"

	RdbKeyFollow   = RdbKey + "follow:"
	RdbKeyFollower = RdbKey + "follower:"
)
