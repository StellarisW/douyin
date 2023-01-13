package middleware

import (
	"douyin/app/common/config/internal/common"
	"douyin/app/common/config/internal/consts"
	"douyin/app/common/config/internal/database"
	"douyin/app/common/middleware"
)

func (g *Group) NewJWTAuthMiddleware() (*middleware.JWTAuthMiddleware, error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	domain, err := common.GetGroup().GetDomain()
	if err != nil {
		return nil, err
	}

	rdb, err := database.GetGroup().GetRedisClusterClient()
	if err != nil {
		return nil, err
	}

	return middleware.NewJWTAuthMiddleware(domain, rdb), nil
}
