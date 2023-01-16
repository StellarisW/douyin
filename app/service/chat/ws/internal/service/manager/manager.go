package manager

import (
	"douyin/app/common/log"
	"douyin/app/service/chat/ws/internal/client"
	"go.uber.org/zap"
	"time"
)

var (
	Manager = NewClientManager()
)

func (m *ClientManager) start() {
	for {
		select {
		case conn := <-m.RegisterChan:
			m.Register(conn)

		case login := <-m.LoginChan:
			m.Login(login)

		case conn := <-m.UnregisterChan:
			m.Unregister(conn)

		case msg := <-m.BroadcastChan:
			clients := m.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

func (m *ClientManager) CheckClientExists(client *client.Client) bool {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()

	_, ok := m.Clients[client]

	return ok
}

func (m *ClientManager) GetClients() map[*client.Client]struct{} {
	clients := make(map[*client.Client]struct{})

	m.DeepCopyClients(func(client *client.Client) {
		clients[client] = struct{}{}
	})

	return clients
}

func (m *ClientManager) AddClient(client *client.Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()

	m.Clients[client] = struct{}{}
}

func (m *ClientManager) DelClient(client *client.Client) {
	m.ClientsLock.Lock()
	defer m.ClientsLock.Unlock()

	if _, ok := m.Clients[client]; ok {
		delete(m.Clients, client)
	}
}

func (m *ClientManager) DeepCopyClients(f func(client *client.Client)) {
	m.ClientsLock.RLock()
	defer m.ClientsLock.RUnlock()

	for key := range m.Clients {
		f(key)
	}
}

func (m *ClientManager) GetClientNum() int {
	return len(m.Clients)
}

func (m *ClientManager) GetUserClient(userId int64) *client.Client {
	m.UserLock.RLock()
	defer m.UserLock.RUnlock()

	if v, ok := m.Users[userId]; ok {
		return v
	}

	return nil
}

func (m *ClientManager) GetUserClientNum() int {
	return len(m.Users)
}

func (m *ClientManager) AddUser(userId int64, client *client.Client) {
	m.UserLock.Lock()
	defer m.UserLock.Unlock()

	m.Users[userId] = client
}

func (m *ClientManager) DelUser(client *client.Client) bool {
	m.UserLock.Lock()
	defer m.UserLock.Unlock()

	if v, ok := m.Users[client.UserId]; ok {
		if v.Addr != client.Addr {
			return false
		}

		delete(m.Users, client.UserId)
	}

	return true
}

func (m *ClientManager) GetUserList() []int64 {
	userList := make([]int64, 0)

	m.UserLock.RLock()
	defer m.UserLock.RUnlock()

	for _, v := range m.Users {
		userList = append(userList, v.UserId)
	}

	return userList
}

func (m *ClientManager) GetUserClients() []*client.Client {
	clients := make([]*client.Client, 0)

	m.UserLock.RLock()
	defer m.UserLock.RUnlock()

	for _, v := range m.Users {
		clients = append(clients, v)
	}

	return clients
}

func (m *ClientManager) SendMessageAll(selfClient *client.Client, message []byte) {
	clients := m.GetUserClients()

	for _, v := range clients {
		if v != selfClient {
			v.SendMessage(message)
		}
	}
}

func (m *ClientManager) Register(client *client.Client) {
	m.AddClient(client)
}

func (m *ClientManager) Login(login *client.Login) {
	c := login.Client

	if m.CheckClientExists(c) {
		m.AddUser(login.UserId, login.Client)
	}

	log.Logger.Debug("user login in", zap.Int64("user_id", login.UserId), zap.String("addr", c.Addr))
}

func (m *ClientManager) Unregister(client *client.Client) {
	m.DelClient(client)

	res := m.DelUser(client)
	if !res {
		return
	}

	log.Logger.Debug("user disconnected", zap.Int64("user_id", client.UserId), zap.String("addr", client.Addr))
}

func GetUserClient(userId int64) (client *client.Client) {
	return Manager.GetUserClient(userId)
}

func ClearTimeoutConnections() {
	now := time.Now().Unix()

	clients := Manager.GetClients()
	for c := range clients {
		if c.IsHeartBeatTimeout(now) {
			log.Logger.Debug("heartbeat timeout, disconnect",
				zap.Int64("user_id", c.UserId),
				zap.String("addr", c.Addr),
				zap.Int64("login_time", c.LoginTime),
				zap.Int64("heartbeat_time", c.HeartbeatTime),
			)

			_ = c.Socket.Close()
		}
	}
}

func GetUserList() []int64 {
	return Manager.GetUserList()
}

func SendMessageAll(userId int64, data string) {
	selfClient := Manager.GetUserClient(userId)

	Manager.SendMessageAll(selfClient, []byte(data))
}
