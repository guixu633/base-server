package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guixu633/base-server/module/config"
	"github.com/sirupsen/logrus"
)

// api: https://platform.openai.com/docs/api-reference/embeddings
// price: https://openai.com/api/pricing/

type OpenaiEmbedEngine struct {
	cfg    *config.Embedding
	client *http.Client
}

type OpenaiEmbedRequest struct {
	// text-embedding-3-small/ text-embedding-3-large
	Model string `json:"model"`
	// 512, 1024, 2048, 4096
	Dimensions int      `json:"dimensions"`
	Input      []string `json:"input"`
	// base64 or float
	EncodingFormat string `json:"encoding_format"`
}

type OpenaiEmbedResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (e *OpenaiEmbedEngine) Embed(text string, query bool) ([]float32, error) {
	// 构建请求体
	reqBody := OpenaiEmbedRequest{
		Model:          e.cfg.OpenaiModel,
		Dimensions:     e.cfg.OpenaiDimension,
		Input:          []string{text},
		EncodingFormat: "float",
	}

	// 创建请求
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", e.cfg.OpenaiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.cfg.OpenaiApiKey)

	// 发送请求
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != http.StatusOK {
		logrus.WithField("status_code", statusCode).Error("request failed")
		return nil, fmt.Errorf("request failed with status code: %d", statusCode)
	}
	// 读取响应
	var embedResp OpenaiEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if len(embedResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	return embedResp.Data[0].Embedding, nil
}

func (e *OpenaiEmbedEngine) EmbedBatch(texts []string, query bool) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		embedding, err := e.Embed(text, query)
		if err != nil {
			return nil, fmt.Errorf("embed text[%d] failed: %w", i, err)
		}
		embeddings[i] = embedding
	}
	return embeddings, nil
}
