package llm

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	openai := getOpenAI(t)
	start := time.Now()
	request := &ChatRequest{
		Model: "gpt-4o",
		Messages: []*Message{
			UserMessage("Hello, who are you?"),
		},
	}
	response, err := openai.Chat(request)
	assert.NoError(t, err)
	fmt.Println(response)
	data, err := json.Marshal(response)
	assert.NoError(t, err)
	fmt.Println(string(data))
	fmt.Println(time.Since(start))
}

func getOpenAI(t *testing.T) *OpenAI {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	openai := NewOpenAI(cfg.LLM.OpenaiApiKey)
	return openai
}
