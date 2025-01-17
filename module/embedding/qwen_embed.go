package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guixu633/base-server/module/config"
)

type QwenEmbedEngine struct {
	cfg    *config.Embedding
	client *http.Client
}

type QwenEmbedRequest struct {
	Model string `json:"model"`
	Input struct {
		Texts []string `json:"texts"`
	} `json:"input"`
	Parameters QwenEmbedRequestParams `json:"parameters"`
}

type QwenEmbedRequestParams struct {
	// 1024, 768, and 512, default is 1024
	Dimension int `json:"dimension,omitempty"`
	// query and document, default is document
	TextType string `json:"text_type,omitempty"`
}

type QwenEmbedResponse struct {
	Output struct {
		Embeddings []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
}

func (e *QwenEmbedEngine) Embed(text string, query bool) ([]float32, error) {
	embeddings, err := e.EmbedBatch([]string{text}, query)
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return embeddings[0], nil
}

func (e *QwenEmbedEngine) EmbedBatch(texts []string, query bool) ([][]float32, error) {
	// 构建请求体
	reqBody := QwenEmbedRequest{
		Model: e.cfg.QwenModel,
		Parameters: QwenEmbedRequestParams{
			Dimension: e.cfg.QwenDimension,
			TextType:  "document",
		},
	}
	if query {
		reqBody.Parameters.TextType = "query"
	}
	reqBody.Input.Texts = texts

	// 创建请求
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", e.cfg.QwenUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.cfg.QwenApiKey)

	// 发送请求
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var embedResp QwenEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 提取嵌入向量
	embeddings := make([][]float32, len(embedResp.Output.Embeddings))
	for i, emb := range embedResp.Output.Embeddings {
		embeddings[i] = emb.Embedding
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embeddings, nil
}
