package services_websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader 将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:12345"
		return true
	},
}

type Client struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

var clients = make(map[string]*Client)

// 定义最大连接数
const maxClients = 2

// 定义计数器和互斥锁
var (
	connectionCount int
	mutex           sync.Mutex
)

func (c *Client) ReadJSON(v any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.ReadJSON(v)
}
