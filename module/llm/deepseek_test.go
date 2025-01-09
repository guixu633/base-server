package llm

import (
	"fmt"
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func TestDeepSeek(t *testing.T) {
	deepseek := getDeepSeek(t)
	request := &ChatRequest{
		Model: "deepseek-chat",
		Messages: []*Message{
			UserMessage("Hello, who are you?"),
		},
	}
	response, err := deepseek.Chat(request)
	assert.NoError(t, err)
	fmt.Println(response)
}

func getDeepSeek(t *testing.T) *DeepSeek {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	deepseek := NewDeepSeek(cfg.LLM.DeepseekApiKey)
	return deepseek
}
