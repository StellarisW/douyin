package listen

import (
	"douyin/app/common/log"
	"douyin/app/service/chat/api/internal/ws/internal/service/listen/internal/handler"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Service struct {
	channel chan os.Signal
	logger  *zap.Logger

	wsHandler  func(response http.ResponseWriter, r *http.Request)
	listenPort string
}

func (m *Service) Start() {
	http.HandleFunc("/", m.wsHandler)

	go func() {
		if err := http.ListenAndServe(":"+m.listenPort, nil); err != http.ErrServerClosed {
			m.logger.Fatal("http server err", zap.Error(err))
		}
	}()

	signal.Notify(m.channel, syscall.SIGINT)
	<-m.channel
}

func (m *Service) Stop() {
	close(m.channel)
}

func NewService(listenPort string) (*Service, error) {
	_, err := strconv.ParseInt(listenPort, 10, 32)
	if err != nil {
		return nil, err
	}

	return &Service{
		channel:    make(chan os.Signal),
		logger:     log.Logger,
		wsHandler:  handler.WsPage,
		listenPort: listenPort,
	}, nil
}
