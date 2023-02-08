package middleware

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
	"net/http"
)

type SentinelMiddleware struct{}

func NewSentinelMiddleware(entity *config.Entity, rules []*flow.Rule) *SentinelMiddleware {
	err := sentinel.InitWithConfig(entity)
	if err != nil {
		panic("invalid config")
	}

	_, err = flow.LoadRules(rules)
	if err != nil {
		panic("invalid flow rule")
	}

	return &SentinelMiddleware{}
}

func (m *SentinelMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry, err := sentinel.Entry(
			fmt.Sprintf("%s:%s", r.Method, r.RequestURI),
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)
		defer entry.Exit()
		if err != nil {
			log.Logger.Debug(errFlow, zap.String("resource", fmt.Sprintf("%s:%s", r.Method, r.RequestURI)))
			response.Fail(
				w,
				http.StatusInternalServerError,
				errx.Encode(
					errx.Sys,
					sysId,
					douyin.Api,
					douyin.SysIdMiddleware,
					serviceIdSentinel,
					errIdFlow,
					0,
				),
				errx.Internal,
			)
			return
		}

		next(w, r)
	}
}
