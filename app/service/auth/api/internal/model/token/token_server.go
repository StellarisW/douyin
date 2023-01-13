package token

import (
	apollo "douyin/app/common/config"
	"github.com/spf13/cast"
)

var (
	ClientSecrets map[string]string // 客户端密钥
)

// InitClientSecret 初始化客户端密钥
func InitClientSecret() (err error) {
	ClientSecrets = make(map[string]string)

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
		for detailName, detailValue := range details {
			if detailName == "secret" {
				ClientSecrets[k] = cast.ToString(detailValue)
				break
			}
		}
	}
	return nil
}
