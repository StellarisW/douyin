package handler

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"douyin/app/service/auth/api/internal/consts"
	"douyin/app/service/auth/internal/sys"
	"go.uber.org/zap"
	"net/http"

	"douyin/app/service/auth/api/internal/logic"
	"douyin/app/service/auth/api/internal/svc"
	"douyin/app/service/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTokenByAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTokenByAuthReq

		if err := httpx.Parse(r, &req); err != nil {
			log.Logger.Error(errx.ParseHttpRequest, zap.Error(err), zap.Reflect("request", r))
			response.Fail(
				w,
				http.StatusBadRequest,
				errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Api,
					sys.ServiceIdApi,
					consts.ErrIdLogic,
					consts.ErrIdOprGetTokenByAuth,
					0,
				),
				err.Error(),
			)
			return
		}
		log.Logger.Debug("recv:", zap.Reflect("args", req))

		l := logic.NewGetTokenByAuthLogic(r.Context(), svcCtx)

		res, err := l.GetTokenByAuth(&req)
		if err != nil {
			log.Logger.Error(errx.ProcessHttpLogic, zap.Error(err))
		}

		log.Logger.Debug("send:", zap.Reflect("args", res))

		if res.Code != 0 {
			if errx.IsSysErr(res.Code) {
				httpx.WriteJson(w, http.StatusInternalServerError, res)
			} else {
				httpx.WriteJson(w, http.StatusBadRequest, res)
			}
		}

		httpx.WriteJson(w, http.StatusOK, res)
	}
}
