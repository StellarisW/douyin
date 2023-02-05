package user

import "douyin/app/service/user/internal/sys"

const (
	RdbKey         = sys.SysName + ":"
	RdbKeyRegister = RdbKey + "register:"
	RdbKeyRelation = RdbKey + "relation:"

	RdbKeyRegisterSet = RdbKeyRegister + "set"

	RdbKeyFollow   = "{" + RdbKeyRelation + "}follow:"
	RdbKeyFollower = "{" + RdbKeyRelation + "}follower:"
)
