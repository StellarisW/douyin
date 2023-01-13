package auth

import (
	apollo "douyin/app/common/config"
	"fmt"
	"github.com/spf13/cast"
)

type ClientDetail struct {
	ClientId     string `json:"client_id,omitempty"`     // client 的标识
	ClientSecret string `json:"client_secret,omitempty"` // client密钥

	AccessTokenValidityTime  string `json:"access_token_validity_time,omitempty"`  // 访问令牌有效时间,秒
	RefreshTokenValidityTime string `json:"refresh_token_validity_time,omitempty"` // 刷新令牌有效时间,秒

	RegisteredRedirectUri string `json:"registered_redirect_uri,omitempty"` // 重定向地址, 授权码类型中使用
}

var (
	ClientDetails map[string]ClientDetail // 客户端信息
)

// InitClientDetails 初始化客户端信息
func InitClientDetails() (err error) {
	ClientDetails = make(map[string]ClientDetail)

	v, err := apollo.Common().GetViper("auth.yaml")
	if err != nil {
		return err
	}

	mapIface := v.GetStringMap("Client")
	if err != nil {
		return err
	}

	for k, v := range mapIface {
		details := cast.ToStringMap(v)
		var clientDetail ClientDetail
		clientDetail.ClientId = k
		for detailName, detailValue := range details {
			switch detailName {
			case "secret":
				clientDetail.ClientSecret = cast.ToString(detailValue)

			case "accesstokenvaliditytime":
				clientDetail.AccessTokenValidityTime = cast.ToString(detailValue)

			case "refreshtokenvaliditytime":
				clientDetail.RefreshTokenValidityTime = cast.ToString(detailValue)

			case "registeredredirecturi":
				clientDetail.RegisteredRedirectUri = cast.ToString(detailValue)
			}
		}
		ClientDetails[k] = clientDetail
	}
	fmt.Println(ClientDetails)
	return nil
}
