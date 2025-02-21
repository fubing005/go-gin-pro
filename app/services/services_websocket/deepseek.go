package services_websocket

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Ollama API 地址
const ollamaURL = "http://localhost:11434/api/generate"

type DeepseekMessage struct {
	Act      string `json:"act"`
	Type     int    `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

// 通过 Ollama 进行流式响应
func streamOllamaResponse(msg Message) error {
	fmt.Printf("%v\n", msg)

	mutex.Lock()
	senderClient, _ := Clients[msg.Sender]
	mutex.Unlock()

	client := resty.New()

	// 构造请求数据
	requestData := map[string]interface{}{
		"model":  "deepseek-r1:1.5b",
		"prompt": msg.Content,
		"stream": true, // 开启流式输出
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		SetDoNotParseResponse(true). // 避免自动解析响应
		Post(ollamaURL)

	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	// 逐行读取流式数据并发送到 WebSocket
	decoder := json.NewDecoder(resp.RawBody())
	for decoder.More() {
		var result map[string]interface{}
		if err := decoder.Decode(&result); err != nil {
			fmt.Println("解析 Ollama 响应失败:", err)
			break
		}

		// 解析 AI 生成的内容
		if content, ok := result["response"].(string); ok {
			deepseekMessage := DeepseekMessage{
				Act:      "chat_message",
				Type:     msg.Type, //点对点聊天：1，deepseek:2
				Sender:   msg.Sender,
				Receiver: msg.Receiver,
				Content:  content,
			}
			// 使用 WriteJSON 发送数据
			if err := senderClient.Conn.WriteJSON(deepseekMessage); err != nil {
				fmt.Println("发送 WebSocket JSON 失败:", err)
				break
			}
		}
	}

	// 发送完成标识
	errorMsg := struct {
		Act   string `json:"act"`
		Error string `json:"error"`
	}{
		Act:   "message",
		Error: "Deepseek 消息回复完成",
	}
	err = senderClient.Conn.WriteJSON(errorMsg)
	if err != nil {
		return errors.New("Ollamma消息发送失败")
	}
	return nil

}
