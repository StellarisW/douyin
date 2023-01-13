package middleware

import (
	"net/http"
)

type CORSHeader struct {
	AllowMethods     string `mapstructure:"allow-methods" yaml:"allow-methods"`
	AllowHeaders     string `mapstructure:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" yaml:"expose-headers"`
	AllowCredentials bool   `mapstructure:"allow-credentials" yaml:"allow-credentials"`
}

type CORSMiddleware struct {
	Mode string
	List map[string]*CORSHeader
}

const (
	AllowOrigin      = "Access-Control-Allow-Origin"
	AllowHeaders     = "Access-Control-Allow-Headers"
	AllowMethods     = "Access-Control-Allow-Methods"
	ExposeHeaders    = "Access-Control-Expose-Headers"
	AllowCredentials = "Access-Control-Allow-Credentials"
	Vary             = "Vary"
)

func NewCORSMiddleware(mode string, list map[string]*CORSHeader) *CORSMiddleware {
	return &CORSMiddleware{Mode: mode, List: list}
}

func (m *CORSMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set(Vary, "Origin") // https://textslashplain.com/2018/08/02/cors-and-vary/

		switch m.Mode {
		case "allow-all":
			// 放行全部请求
			w.Header().Set(AllowOrigin, origin)
			w.Header().Set(AllowHeaders, "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-Sign-Id")
			w.Header().Set(AllowMethods, "POST, GET, PATCH, OPTIONS, DELETE, PUT")
			w.Header().Set(ExposeHeaders, "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			w.Header().Set(AllowCredentials, "true")
		case "whitelist":
			// 白名单模式
			if v, ok := m.List[origin]; ok {
				w.Header().Set(AllowHeaders, v.AllowHeaders)
				w.Header().Set(AllowMethods, v.AllowMethods)
				w.Header().Set(ExposeHeaders, v.ExposeHeaders)
				if v.AllowCredentials {
					w.Header().Set(AllowCredentials, "true")
				}
			}

		case "blacklist":
			// 黑名单模式
			if _, ok := m.List[origin]; ok {
				return
			}
		}

		next(w, r)
	}
}
