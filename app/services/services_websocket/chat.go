package services_websocket

import (
	"fmt"
	"net/http"
	"shalabing-gin/app/common/request/request_websocket"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Act      string `json:"act"`
	Type     int    `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

func HandleChatConnections(c *gin.Context, form request_websocket.ChatRequest) {
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
		return
	}

	// 设置最大消息大小
	conn.SetReadLimit(512 * 1024)
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

	client := &Client{conn: conn}

	// 监听客户端消息
	// defer func() {
	// 	mutex.Lock()
	// 	delete(clients, form.ID)
	// 	mutex.Unlock()
	// 	client.conn.Close()
	// }()

	mutex.Lock()
	clients[form.ID] = client
	mutex.Unlock()

	// 发送连接成功信息
	// openMsg := fmt.Sprintf(`{"act":"open","error": "用户{\"user_id\": \"%s\"} 连接成功"}`, form.ID)
	// conn.WriteMessage(websocket.TextMessage, []byte(openMsg))

	openMsg := struct {
		Act   string `json:"act"`
		Error string `json:"error"`
	}{
		Act:   "open",
		Error: fmt.Sprintf(`用户{"user_id": "%s"} 连接成功`, form.ID),
	}
	// 使用 WriteJSON 发送 JSON 数据
	err = conn.WriteJSON(openMsg)
	if err != nil {
		fmt.Println("发送 JSON 失败:", err)
	}

	// 时刻等待接收客户端的消息
	for {
		var msg Message
		err := client.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("用户 %s 发送消息失败: %v\n", form.ID, err)
			break
		}

		sendMessage(msg)
	}
}

// SendMessage 发送消息给指定用户
func sendMessage(msg Message) {
	// mutex.Lock()
	senderClient, senderExists := clients[msg.Sender]
	receiverClient, receiverExists := clients[msg.Receiver]
	// mutex.Unlock()

	if !senderExists {
		return
	}

	if receiverExists {
		err := receiverClient.conn.WriteJSON(msg)
		if err != nil {
			fmt.Printf(`{"act":"message","error": "用户{%s} 接收消息失败"}`, msg.Receiver)
			return
		}
	} else {
		errorMsg := struct {
			Act   string `json:"act"`
			Error string `json:"error"`
		}{
			Act:   "message",
			Error: fmt.Sprintf(`"用户{%s} 不在线，消息未送达"}`, msg.Receiver),
		}
		// 使用 WriteJSON 发送 JSON 数据
		err := senderClient.conn.WriteJSON(errorMsg)
		if err != nil {
			fmt.Println("发送 JSON 失败:", err)
		}

		// senderClient.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"act":"message","error": "用户{%s} 不在线，消息未送达"}`, msg.Receiver)))
	}
}
