package log

// Options 日志器设置
type Options struct {
	Mode         string // 日志模式(debug,release)
	SavePath     string // 日志保存位置
	EncoderType  string // 日志编码器类型("json","console")
	EncodeLevel  string // 日志编码器风格
	EncodeCaller string // 打印调用函数风格
}

const (
	// 编码器类型

	// JsonEncoder Json编码器
	JsonEncoder = "json"
	// ConsoleEncoder 控制台编码器
	ConsoleEncoder = "console"

	// 编码器等级

	// LowercaseLevelEncoder 小写编码器
	LowercaseLevelEncoder = "LowercaseLevelEncoder"
	// LowercaseColorLevelEncoder 小写编码器带颜色
	LowercaseColorLevelEncoder = "LowercaseColorLevelEncoder"
	// CapitalLevelEncoder 大写编码器
	CapitalLevelEncoder = "CapitalLevelEncoder"
	// CapitalColorLevelEncoder 大写编码器带颜色
	CapitalColorLevelEncoder = "CapitalColorLevelEncoder"

	// 编码器打印调用函数方式
	// 日志文件信息 short（包/文件.go:行号） full (文件位置.go:行号)

	// ShortCallerEncoder 调用函数的大概位置
	ShortCallerEncoder = "ShortCallerEncoder"
	// FullCallerEncoder 调用函数的详细位置
	FullCallerEncoder = "FullCallerEncoder"
)
