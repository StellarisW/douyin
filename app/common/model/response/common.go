package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Response struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

type WithData struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(w http.ResponseWriter, msg string) {
	httpx.WriteJson(w, http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
	})
}

func OkWithData(w http.ResponseWriter, msg string, data interface{}) {
	httpx.WriteJson(w, http.StatusOK, WithData{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

func Fail(w http.ResponseWriter, statusCode int, code uint32, msg string) {
	httpx.WriteJson(w, statusCode,
		Response{
			Code: code,
			Msg:  msg,
		},
	)
}
