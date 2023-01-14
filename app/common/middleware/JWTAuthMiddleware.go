package middleware

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-redis/redis/v9"
)

type JWTAuthMiddleware struct {
	Domain string
	Rdb    *redis.ClusterClient
}

func NewJWTAuthMiddleware(domain string, rdb *redis.ClusterClient) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{Domain: domain, Rdb: rdb}
}

func (m *JWTAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		if r.Method == http.MethodGet {
			token = r.URL.Query().Get("token")
			if token == "" {
				_ = r.ParseMultipartForm(1 << 12)
				token = r.FormValue("token")
			}
		} else {
			_ = r.ParseMultipartForm(1 << 12)
			token = r.FormValue("token")
		}

		// 检验 accessToken 有效性
		readTokenRes, err := req.NewRequest().
			SetFormData(map[string]string{
				"token_value": token,
			}).
			Post("http://douyin-auth-api:11120/douyin/token/read")
		if err != nil {
			log.Logger.Error(errx.RequestHttpSend, zap.Error(err))
			response.Fail(
				w,
				http.StatusInternalServerError,
				errx.Encode(
					errx.Sys,
					sysId,
					douyin.Api,
					serviceIdJWT,
					0,
					0,
					errIdRequestHttpSendAuth),
				errx.RequestHttpSend,
			)
			return
		}

		readTokenResJson := gjson.Parse(readTokenRes.String())

		if readTokenResJson.Get("code").Uint() != 0 {
			response.Fail(
				w,
				http.StatusUnauthorized,
				errx.Encode(
					errx.Logic,
					sysId,
					douyin.Api,
					serviceIdJWT,
					0,
					0,
					errIdInvalidToken,
				),
				readTokenResJson.Get("msg").String(),
			)
			return
		} else {
			token = readTokenResJson.Get("data.payload").String()
		}

		tokenPayloadJson := gjson.Parse(token)

		sub := tokenPayloadJson.Get("sub").String()
		scope := tokenPayloadJson.Get("scope").String()

		// 设置上下文
		r = r.WithContext(context.WithValue(r.Context(), KeyUserId, sub))
		r = r.WithContext(context.WithValue(r.Context(), KeyScope, scope))

		next(w, r)
	}
}
