package common

import "douyin/app/common/config/internal/types"

type Group struct {
	agollo *types.Agollo
}

var insCommon = &Group{}

func GetGroup() *Group {
	if insCommon.agollo == nil {
		insCommon = &Group{agollo: types.AgolloClient}
	}

	return insCommon
}
