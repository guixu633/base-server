package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ChatflowRequest struct {
	Inputs         map[string]string `json:"inputs"`
	ResponseMode   string            `json:"response_mode"`
	User           string            `json:"user"`
	Query          string            `json:"query"`
	ConversationId string            `json:"conversation_id,omitempty"`
}

type ChatflowResponse struct {
	ConversationId string `json:"conversation_id"`
	Answer         string `json:"answer"`
}

func (w *Workflow) CallChatflowBlock(ctx context.Context, token, query, conversationId string, inputs map[string]string) (string, string, error) {
	// 构造请求体
	requestBody := &ChatflowRequest{
		Inputs:         inputs,
		ResponseMode:   "blocking",
		User:           "guixu633",
		Query:          query,
		ConversationId: conversationId,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", fmt.Errorf("JSON编码失败: %v", err)
	}

	logrus.WithField("body", string(jsonBody)).Info("workflow call block with inputs")

	// 构造请求
	url := fmt.Sprintf("%s/v1/chat-messages", w.cfg.Url)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 发送请求
	resp, err := w.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return "", "", fmt.Errorf("读取响应失败: %v", err)
	}
	logrus.WithField("inputs", inputs).WithField("body", string(body)).Info(ctx, "workflow call block success")
	var response ChatflowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %v", err)
	}

	return response.ConversationId, response.Answer, nil
}
