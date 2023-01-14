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
	InitMinio        = "initialize minio client failed"

	GetViper       = "get viper failed"
	GetIdGenerator = "get idGenerator failed"
)

const (
	// Api层的独有错误

	ParseHttpRequest = "parse http request failed"
	ProcessHttpLogic = "process logic failed"
)

const (
	// Mysql类

	MysqlExec   = "exec sql failed"
	MysqlGet    = "query mysql record failed"
	MysqlInsert = "insert mysql record failed"
	MysqlUpdate = "update mysql record failed"
	MysqlDelete = "delete mysql record failed"
)

const (
	// Redis类

	RedisGet   = "get redis key failed"
	RedisSet   = "set redis key failed"
	RedisDel   = "del redis key failed"
	RedisScan  = "scan redis key failed"
	RedisAdd   = "add redis set member failed"
	RedisRem   = "remove redis set member failed"
	RedisRange = "get redis set members in range failed"
	RedisInter = "get redis set inter members failed"
	RedisIncr  = "incr redis key failed"
	RedisDecr  = "decr redis key failed"
)

const (
	// 网络请求类

	RequestHttpSend       = "send http request failed"
	RequestHttpStatusCode = "status code err"
	RequestRpcRes         = "rpc res not ok"
)

const (

	// Minio

	MinioPut    = "minio client put object failed"
	MinioRemove = "minio client remove object failed"
)

const (
	// bytes类(序列化,解析等)

	UnmarshalServiceConfig = "unmarshal viper key into service config failed"

	JsonUnmarshal = "unmarshal json failed"
	JsonMarshal   = "marshal json failed"
	ReadBytes     = "read bytes failed"
)
