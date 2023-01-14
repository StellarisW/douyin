package handler

import (
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/crud"
	"douyin/app/service/video/internal/sys"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"douyin/app/service/video/api/internal/logic"
	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishReq

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
					consts.ErrIdLogicCrud,
					crud.ErrIdOprPublish,
					0,
				),
				err.Error(),
			)
			return
		}
		log.Logger.Debug("recv:", zap.Reflect("args", req))

		l := logic.NewPublishLogic(r.Context(), svcCtx)

		fhs := r.MultipartForm.File["data"]

		var contentBytes []byte

		if len(fhs) > 0 {
			file, err := fhs[0].Open()
			if err != nil {
				log.Logger.Error(crud.ErrOpenFile, zap.Error(err))
				response.Fail(
					w,
					http.StatusInternalServerError,
					errx.Encode(
						errx.Sys,
						sys.SysId,
						douyin.Api,
						sys.ServiceIdApi,
						consts.ErrIdLogicCrud,
						crud.ErrIdOprPublish,
						crud.ErrIdOpenFile,
					),
					crud.ErrOpenFile,
				)
			}

			contentBytes, err = io.ReadAll(file)
			if err != nil {
				log.Logger.Error(errx.ReadBytes, zap.Error(err))
				response.Fail(
					w,
					http.StatusInternalServerError,
					errx.Encode(
						errx.Sys,
						sys.SysId,
						douyin.Api,
						sys.ServiceIdApi,
						consts.ErrIdLogicCrud,
						crud.ErrIdOprPublish,
						crud.ErrIdReadBytes,
					),
					errx.ReadBytes,
				)
			}

			// 判断文件类型
			contentType := http.DetectContentType(contentBytes)
			if !strings.HasPrefix(contentType, "video") {
				httpx.WriteJson(
					w,
					http.StatusForbidden,
					&types.PublishRes{
						StatusCode: errx.Encode(
							errx.Logic,
							sys.SysId,
							douyin.Api,
							sys.ServiceIdApi,
							consts.ErrIdLogicCrud,
							crud.ErrIdOprPublish,
							crud.ErrIdInvalidVideoType,
						),
						StatusMsg: crud.ErrInvalidVideoType,
					})
				return
			}
		}

		res, err := l.Publish(&req, contentBytes)
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
