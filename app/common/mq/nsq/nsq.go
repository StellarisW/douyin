package nsq

import (
	apollo "douyin/app/common/config"
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/reiver/go-telnet"
	"math/rand"
	"strings"
	"time"
)

func getConfig() (*nsq.Config, error) {
	config := nsq.NewConfig()

	v, err := apollo.Common().GetViper("NSQ")
	if err != nil {
		return nil, err
	}

	config.AuthSecret = v.GetString("nsq-auth.secret")

	return config, nil

}

func mustGetNSQDAddr() string {
	addr, err := getNSQDAddr()
	if err != nil {
		return ""
	}
	return addr
}

func getNSQDAddr() (string, error) {
	v, err := apollo.Common().GetViper("NSQ")
	if err != nil {
		return "", err
	}

	addrsSrc := v.GetString("nsq-nsqd.tcp")
	addrs := strings.Split(addrsSrc, ";")
	addrNum := len(addrs)

	// 简单的负载均衡
	rand.Seed(time.Now().Unix())
	serverNum := rand.Intn(addrNum - 1)

	cnt := 0
	for i := serverNum; cnt != addrNum-1; i++ {
		addr := addrs[i]
		// 检测端口连通性
		_, err := telnet.DialTo(addr)
		if err == nil {
			return addr, nil
		}

		if i == addrNum {
			i = 0
		}
		cnt++
	}
	return "", errors.New("no servers available")
}

func mustGetNSQLookupAddrs() []string {
	addrs, err := getNSQLookupdAddrs()
	if err != nil {
		return nil
	}
	return addrs
}

func getNSQLookupdAddrs() ([]string, error) {
	v, err := apollo.Common().GetViper("NSQ")
	if err != nil {
		return nil, err
	}

	addrsSrc := v.GetString("nsq-lookupd.http")
	addrs := strings.Split(addrsSrc, ";")

	return addrs, nil
}
