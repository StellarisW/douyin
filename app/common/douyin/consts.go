package douyin

const (
	ModeEnvName = "DOUYIN_MODE"
	DevMode     = "dev" // 开发环境, 后端人员使用
	FatMode     = "fat" // 功能验收测试环境, 前端人员使用
	UatMode     = "uat" // 用户验收测试环境, 模拟上线环境
	ProMode     = "pro" // 生产环境, 正式环境
)

const (
	Api = iota
	Rpc
)

const (
	SysIdMiddleware = iota
	SysIdMq
	SysIdAuth
	SysIdUser
	SysIdVideo
	SysIdChat
)

const (
	SysNameMiddleware = "middleware"
	SysNameMq         = "mq"
	SysNameAuth       = "auth"
	SysNameUser       = "user"
	SysNameVideo      = "video"
	SysNameChat       = "chat"
)
