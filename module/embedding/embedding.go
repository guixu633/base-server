package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/guixu633/base-server/module/config"
)

type EmbedEngine struct {
	cfg    *config.Embedding
	client *http.Client
}

type EmbedRequest struct {
	Model string `json:"model"`
	Input struct {
		Texts []string `json:"texts"`
	} `json:"input"`
	Parameters EmbedRequestParams `json:"parameters"`
}

type EmbedRequestParams struct {
	// 1024, 768, and 512, default is 1024
	Dimension int `json:"dimension,omitempty"`
	// query and document, default is document
	TextType string `json:"text_type,omitempty"`
}

type EmbedResponse struct {
	Output struct {
		Embeddings []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
}

func NewEmbedEngine(cfg *config.Embedding, client *http.Client) *EmbedEngine {
	return &EmbedEngine{
		cfg:    cfg,
		client: client,
	}
}

func (e *EmbedEngine) Embed(text string, query bool) ([]float32, error) {
	embeddings, err := e.EmbedBatch([]string{text}, query)
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

func (e *EmbedEngine) EmbedBatch(texts []string, query bool) ([][]float32, error) {
	// 构建请求体
	reqBody := EmbedRequest{
		Model: e.cfg.Model,
		Parameters: EmbedRequestParams{
			Dimension: e.cfg.Dimension,
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

	req, err := http.NewRequest("POST", e.cfg.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.cfg.ApiKey)

	// 发送请求
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 提取嵌入向量
	embeddings := make([][]float32, len(embedResp.Output.Embeddings))
	for i, emb := range embedResp.Output.Embeddings {
		embeddings[i] = emb.Embedding
	}

	return embeddings, nil
}

func ConsineSimilarity(embedding1, embedding2 []float32) float32 {
	dotProduct := float32(0)
	for i := range embedding1 {
		dotProduct += embedding1[i] * embedding2[i]
	}
	return dotProduct / (L2Norm(embedding1) * L2Norm(embedding2))
}

func L2Norm(embedding []float32) float32 {
	sum := float32(0)
	for _, value := range embedding {
		sum += value * value
	}
	return float32(math.Sqrt(float64(sum)))
}
