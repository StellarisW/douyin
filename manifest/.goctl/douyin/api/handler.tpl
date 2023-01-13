package {{.PkgName}}

import (
	"go.uber.org/zap"
	"net/http"
	"douyin/app/common/log"
	"douyin/app/common/cqupt"
	"douyin/app/common/errx"
	"douyin/app/common/model/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		
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
					consts.ErrIdLogic, // TODO: add your logic id here
				    .ErrIdOpr, // TODO: add your opr id here
					0,
				),
				err.Error(),
			)
			return
		}
		log.Logger.Debug("recv:", zap.Reflect("args", req))

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		
		{{if .HasResp}}res, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
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
