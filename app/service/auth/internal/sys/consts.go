package sys

import "douyin/app/common/douyin"

const (
	SysId = douyin.SysIdAuth

	ServiceIdApi = iota - 1

	ServiceIdRpcStore = iota - 2
	ServiceIdRpcEnhancer
)
