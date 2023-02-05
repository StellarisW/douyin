package handler

import (
	"bytes"
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/model/response"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/crud"
	"douyin/app/service/video/internal/sys"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"douyin/app/service/video/api/internal/logic"
	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishReq

		if err := r.ParseForm(); err != nil {
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

		if err := r.ParseMultipartForm(256 << 20); err != nil {
			if err != http.ErrNotMultipart {
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
		}

		req.Token = r.Form.Get("token")
		req.Title = r.Form.Get("title")

		log.Logger.Debug("recv:", zap.Reflect("args", req))

		l := logic.NewPublishLogic(r.Context(), svcCtx)

		fhs := r.MultipartForm.File["data"]

		var contentBytes []byte
		var contentType string
		var videoBuffer *bytes.Buffer
		var imgBuffer *bytes.Buffer

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
			contentType = http.DetectContentType(contentBytes)
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

			videoBuffer = bytes.NewBuffer(contentBytes)
			imgBuffer = bytes.NewBuffer(nil)
			cmdOutput := bytes.NewBuffer(nil)

			err = ffmpeg.Input("-", ffmpeg.KwArgs{"loglevel": "debug"}).
				Filter("select", ffmpeg.Args{"gte(n,1)"}).
				Output("-", ffmpeg.KwArgs{"vframes": 1, "q:v": "2", "f": "image2"}).
				OverWriteOutput().
				WithOutput(imgBuffer, os.Stdout).
				WithErrorOutput(cmdOutput).
				WithInput(videoBuffer).
				WithTimeout(3 * time.Second).
				Run()
			if err != nil {
				log.Logger.Error(crud.ErrGetVideoImage, zap.Error(errors.New(cmdOutput.String())))
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
							crud.ErrIdGetVideoImage,
						),
						StatusMsg: errx.Internal,
					})
				return
			}

			videoBuffer.Reset()
			videoBuffer.Write(contentBytes)
		}

		res, videoId, err := l.Publish(&req)
		if err != nil {
			log.Logger.Error(errx.ProcessHttpLogic, zap.Error(err))
		}

		_, err = svcCtx.MinioClient.PutObject(context.Background(),
			douyin.MinioBucket,
			fmt.Sprintf("video/%d/cover", videoId),
			imgBuffer,
			int64(imgBuffer.Len()),
			minio.PutObjectOptions{
				ContentType: "image/jpeg",
			},
		)
		if err != nil {
			log.Logger.Error(errx.MinioPut, zap.Error(err))
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
						crud.ErrIdMinioPut,
					),
					StatusMsg: errx.Internal,
				})
			return
		}

		_, err = svcCtx.MinioClient.PutObject(context.Background(),
			douyin.MinioBucket,
			fmt.Sprintf("video/%d/video", videoId),
			videoBuffer,
			int64(videoBuffer.Len()),
			minio.PutObjectOptions{
				ContentType: contentType,
			},
		)
		if err != nil {
			log.Logger.Error(errx.MinioPut, zap.Error(err))
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
						crud.ErrIdMinioPut,
					),
					StatusMsg: errx.Internal,
				})
			return
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
