package embedding

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func TestEmbed(t *testing.T) {
	engine := getEmbedEngine(t)
	data := []string{
		"人工智能技术正在快速发展，深度学习模型能够处理越来越复杂的任务。从图像识别到自然语言处理，AI 系统展现出惊人的能力。特别是在医疗诊断、自动驾驶和智能助手等领域，AI 已经开始改变我们的生活方式。",
		"随着技术的进步，人工智能的应用范围不断扩大。深度学习算法现在可以完成许多复杂任务。AI 在图像处理和语言理解方面取得了重大突破，并在医疗、交通和个人助理等领域发挥重要作用。",
	}
	embeddings, err := engine.EmbedBatch(data, false)
	assert.NoError(t, err)
	fmt.Println(embeddings)
	similarity := ConsineSimilarity(embeddings[0], embeddings[1])
	fmt.Println(similarity)
}

func getEmbedEngine(t *testing.T) *EmbedEngine {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	client := &http.Client{}
	return NewEmbedEngine(&cfg.Embedding, client)
}
