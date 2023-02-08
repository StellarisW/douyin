package pprof

import (
	"douyin/app/common/log"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	Host string
	Port int
}

func (ps Server) Start() {
	log.Logger.Info("starting server...", zap.String("host", ps.Host), zap.Int("port", ps.Port))

	addr := fmt.Sprintf("%s:%d", ps.Host, ps.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Logger.Fatal("initialize pprof server failed.", zap.Error(err))
	}
}

func (ps Server) Stop() {
	log.Logger.Info("pprof server stop")
}
