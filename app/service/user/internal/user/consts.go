package user

import "douyin/app/service/user/internal/sys"

const (
	RdbKey         = sys.SysName + ":"
	RdbKeyRelation = RdbKey + "relation:"

	RdbKeyFollow   = "{" + RdbKeyRelation + "}follow:"
	RdbKeyFollower = "{" + RdbKeyRelation + "}follower:"
)
