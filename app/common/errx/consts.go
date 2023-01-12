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
)
