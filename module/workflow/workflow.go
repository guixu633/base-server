package workflow

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/guixu633/base-server/module/config"
	"github.com/sirupsen/logrus"
)

type Workflow struct {
	cfg    *config.Workflow
	client *http.Client
}

func NewWorkflow(cfg *config.Workflow, client *http.Client) *Workflow {
	return &Workflow{
		cfg:    cfg,
		client: client,
	}
}

type WorkflowRequest struct {
	Inputs       map[string]string `json:"inputs"`
	ResponseMode string            `json:"response_mode"`
	User         string            `json:"user"`
}

type WorkflowResponse struct {
	Data struct {
		Error   string                 `json:"error"`
		Outputs map[string]interface{} `json:"outputs"`
	} `json:"data"`
}

func (w *Workflow) CallWorkflowStream(ctx context.Context, token string, inputs map[string]string) error {
	// 构造请求体
	requestBody := &WorkflowRequest{
		Inputs:       inputs,
		ResponseMode: "streaming",
		User:         "guixu633",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	// 构造请求
	req, err := http.NewRequestWithContext(ctx, "POST", w.cfg.Url+"/v1/workflows/run", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "text/event-stream") // 添加这行以请求流式响应

	// 发送请求
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取流式响应
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取响应失败: %v", err)
		}

		// 处理每一行数据
		fmt.Printf("收到数据: %s", string(line))
	}

	return nil
}

func (w *Workflow) CallWorkflowBlock(ctx context.Context, token string, inputs map[string]string) (map[string]interface{}, error) {
	// 构造请求体
	requestBody := &WorkflowRequest{
		Inputs:       inputs,
		ResponseMode: "blocking",
		User:         "guixu633",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %v", err)
	}

	logrus.WithField("body", string(jsonBody)).Info(ctx, "workflow call block with inputs")

	// 构造请求
	url := fmt.Sprintf("%s/v1/workflows/run", w.cfg.Url)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 发送请求
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	logrus.WithField("inputs", inputs).WithField("body", string(body)).Info(ctx, "workflow call block success")
	var response WorkflowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if response.Data.Error != "" {
		return nil, fmt.Errorf("workflow call error: %s", response.Data.Error)
	}
	return response.Data.Outputs, nil
}

func parseStrList(data map[string]interface{}, key string) ([]string, error) {
	content, ok := data[key].([]interface{})
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	result := make([]string, len(content))
	for i, v := range content {
		result[i] = v.(string)
	}
	return result, nil
}
