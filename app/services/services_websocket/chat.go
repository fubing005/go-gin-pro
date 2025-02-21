package services_websocket

import (
	"fmt"
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

func HandleChatMessage(c *gin.Context, form request_websocket.ChatRequest) {
	client, err := GetWsConn(c, form.ID)
	if err != nil {
		fmt.Println("获取 WebSocket 连接失败:", err)
	}

	openMsg := struct {
		Act   string `json:"act"`
		Error string `json:"error"`
	}{
		Act:   "open",
		Error: fmt.Sprintf(`用户{"user_id": "%s"} 连接成功`, form.ID),
	}
	// 使用 WriteJSON 发送 JSON 数据
	err = client.Conn.WriteJSON(openMsg)
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

		// 通过 Ollama 生成 AI 回复
		if msg.Type == 1 { //普通聊天
			sendMessage(msg)
		} else { //ai对话
			if msg.Content == "" {
				fmt.Println("请求中缺少 'content' 字段")
				continue
			}
			if err := streamOllamaResponse(msg); err != nil {
				fmt.Println("Ollama 处理失败:", err)
				break
			}
		}
	}
}

// SendMessage 发送消息给指定用户
func sendMessage(msg Message) {
	mutex.Lock()
	senderClient, senderExists := Clients[msg.Sender]
	receiverClient, receiverExists := Clients[msg.Receiver]
	mutex.Unlock()

	if !senderExists {
		return
	}

	if receiverExists {
		err := receiverClient.Conn.WriteJSON(msg)
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
		err := senderClient.Conn.WriteJSON(errorMsg)
		if err != nil {
			fmt.Println("发送 JSON 失败:", err)
		}
	}
}
