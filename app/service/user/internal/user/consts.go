package user

import "douyin/app/service/user/internal/sys"

const (
	rdbKey         = sys.SysName + ":"
	rdbKeyRegister = rdbKey + "register:"
	rdbKeyLogin    = rdbKey + "login:"
	rdbKeyRelation = rdbKey + "relation:"

	RdbKeyRegisterSet = rdbKeyRegister + "set"

	rdbKeyLoginFrozen         = rdbKeyLogin + "frozen:"
	RdbKeyLoginFrozenTime     = rdbKeyLoginFrozen + "time:"
	RdbKeyLoginFrozenTimeLast = rdbKeyLoginFrozen + "time:last:"
	RdbKeyLoginFrozenLoginCnt = rdbKeyLoginFrozen + "cnt:"

	RdbKeyFollow   = "{" + rdbKeyRelation + "}follow:"
	RdbKeyFollower = "{" + rdbKeyRelation + "}follower:"
)
