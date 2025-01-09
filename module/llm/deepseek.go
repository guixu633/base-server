package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeepSeek struct {
	apiKey string
	auth   string
	client *http.Client
}

func NewDeepSeek(apiKey string) *DeepSeek {
	return &DeepSeek{apiKey: apiKey, auth: "Bearer " + apiKey, client: &http.Client{}}
}

func (d *DeepSeek) Chat(request *ChatRequest) (*ChatResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(requestBody))

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", d.auth)

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应
	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return &chatResponse, nil
}
