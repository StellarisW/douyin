package errx

const (
	Logic = iota
	Sys
)

const (
	Internal = "internal err"
)

const (
	// 初始化发生的系统错误

	EmptyProjectMode = "project mode not set, please consider configure env: DOUYIN_MODE"
	InitAgolloClient = "initialize Apollo Client failed."
	InitMysql        = "initialize mysql failed"
	InitRedis        = "initialize redis failed"
)

const (
	// Api层的独有错误

	ParseHttpRequest = "parse http request failed"
	ProcessHttpLogic = "process logic failed"
)

const (
	// Redis类

	RedisGet  = "get redis key failed"
	RedisSet  = "set redis key failed"
	RedisDel  = "del redis key failed"
	RedisScan = "scan redis key failed"
)

const (
	// 网络请求类

	RequestHttpSend       = "send http request failed"
	RequestHttpStatusCode = "status code err"
	RequestRpcRes         = "rpc res not ok"
)

const (
	// bytes类(序列化,解析等)

	UnmarshalServiceConfig = "unmarshal viper key into service config failed"

	JsonUnmarshal = "unmarshal json failed"
	JsonMarshal   = "marshal json failed"
)
