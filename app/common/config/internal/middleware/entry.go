package middleware

import "douyin/app/common/config/internal/types"

type Group struct {
	agollo *types.Agollo
}

var insMiddleware = &Group{}

func GetGroup() *Group {
	if insMiddleware.agollo == nil {
		insMiddleware = &Group{agollo: types.AgolloClient}
	}

	return insMiddleware
}
