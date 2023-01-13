package handler

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"douyin/app/service/chat/api/internal/consts"
	"douyin/app/service/chat/api/internal/consts/chat"
	"douyin/app/service/chat/internal/sys"
	"go.uber.org/zap"
	"net/http"

	"douyin/app/service/chat/api/internal/logic"
	"douyin/app/service/chat/api/internal/svc"
	"douyin/app/service/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMessageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMessageListReq

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
					chat.ErrIdOprGetMessageList,
					0,
				),
				err.Error(),
			)
			return
		}
		log.Logger.Debug("recv:", zap.Reflect("args", req))

		l := logic.NewGetMessageListLogic(r.Context(), svcCtx)

		res, err := l.GetMessageList(&req)
		if err != nil {
			log.Logger.Error(errx.ProcessHttpLogic, zap.Error(err))
		}

		log.Logger.Debug("send:", zap.Reflect("args", res))

		if res.StatusCode != 0 {
			if errx.IsSysErr(res.StatusCode) {
				httpx.WriteJson(w, http.StatusInternalServerError, res)
			} else {
				httpx.WriteJson(w, http.StatusBadRequest, res)
			}

		}

		httpx.WriteJson(w, http.StatusOK, res)
	}
}
