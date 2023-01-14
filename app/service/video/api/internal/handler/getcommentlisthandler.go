package handler

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/comment"
	"douyin/app/service/video/internal/sys"
	"go.uber.org/zap"
	"net/http"

	"douyin/app/service/video/api/internal/logic"
	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCommentListReq

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
					consts.ErrIdLogicComment,
					comment.ErrIdOprGetCommentList,
					0,
				),
				err.Error(),
			)
			return
		}
		log.Logger.Debug("recv:", zap.Reflect("args", req))

		l := logic.NewGetCommentListLogic(r.Context(), svcCtx)

		res, err := l.GetCommentList(&req)
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
			return
		}

		httpx.WriteJson(w, http.StatusOK, res)
	}
}
