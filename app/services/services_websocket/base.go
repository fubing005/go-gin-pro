package services_websocket

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader 将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:12345"
		return true
	},
}

type Client struct {
	Conn  *websocket.Conn
	mutex sync.Mutex
}

var Clients = make(map[string]*Client)

// 定义最大连接数
const maxClients = 2

// 定义计数器和互斥锁
var (
	connectionCount int
	mutex           sync.Mutex
)

func GetWsConn(c *gin.Context, clientId string) (client *Client, err error) {
	// 限制最大连接数，加锁以确保线程安全
	// mutex.Lock()
	// defer mutex.Unlock()
	// // 检查连接数是否超过最大限制
	// if connectionCount > maxClients {
	// 	http.Error(c.Writer, "达到最大连接数", http.StatusServiceUnavailable)
	// 	return
	// }
	// // 增加连接数
	// connectionCount++
	// defer func() {
	// 	connectionCount--
	// }()

	//// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket 连接失败"})
		return nil, err
	}

	// 设置最大消息大小
	// conn.SetReadLimit(512 * 1024)
	// // 设置读取超时时间
	// conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	// // 设置写入超时时间
	// conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	// // 设置压缩
	// conn.EnableWriteCompression(true)
	// conn.SetCompressionLevel(gzip.BestSpeed)
	// // 设置 ping handler
	// conn.SetPingHandler(func(appData string) error {
	// 	// 收到 ping 后更新读取超时
	// 	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	// 	// 回复 pong
	// 	return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
	// })
	// // 设置 pong handler
	// conn.SetPongHandler(func(appData string) error {
	// 	// 收到 pong 后更新读取超时
	// 	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	// 	return nil
	// })

	client = &Client{Conn: conn}

	// 读取到客户端消息后，关闭连接，并且从 clients 中删除
	// defer func() {
	// 	mutex.Lock()
	// 	delete(clients, form.ID)
	// 	mutex.Unlock()
	// 	client.conn.Close()
	// }()

	mutex.Lock()
	Clients[clientId] = client
	mutex.Unlock()
	return client, nil
}

func (c *Client) ReadJSON(v any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Conn.ReadJSON(v)
}
