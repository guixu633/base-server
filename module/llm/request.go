package llm

// https://platform.openai.com/docs/api-reference/chat/create
// https://openai.com/api/pricing/
type ChatRequest struct {
	Model               string          `json:"model"`
	Messages            []*Message      `json:"messages"`
	Temperature         *float64        `json:"temperature,omitempty"`       // 0 to 2
	FrequencyPenalty    *float64        `json:"frequency_penalty,omitempty"` // -2.0 to 2.0
	PresencePenalty     *float64        `json:"presence_penalty,omitempty"`  // -2.0 to 2.0
	MaxCompletionTokens *int            `json:"max_completion_tokens,omitempty"`
	N                   *int            `json:"n,omitempty"`
	ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
	Stop                *[]string       `json:"stop,omitempty"`
	Stream              *bool           `json:"stream,omitempty"`

	// https://platform.openai.com/tokenizer
	LogitBias map[int]int `json:"logit_bias,omitempty"` // -100 to 100
	// 显示 next token 的概率
	Logprobs *bool `json:"logprobs,omitempty"`
	// 显示 topn 的 token
	TopLogprobs *int `json:"top_logprobs,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

func BuildChatRequest(model string) *ChatRequest {
	return &ChatRequest{
		Model: model,
	}
}

func (r *ChatRequest) SetTemperature(temperature float64) *ChatRequest {
	if temperature < 0 || temperature > 2 {
		return r
	}
	r.Temperature = &temperature
	return r
}

func (r *ChatRequest) SetMaxCompletionTokens(maxTokens int) *ChatRequest {
	if maxTokens < 1 {
		return r
	}
	r.MaxCompletionTokens = &maxTokens
	return r
}

func (r *ChatRequest) SetResponseJsonFormat() *ChatRequest {
	r.ResponseFormat = &ResponseFormat{
		Type: "json_object",
	}
	return r
}

func (r *ChatRequest) SetN(n int) *ChatRequest {
	if n > 10 {
		n = 10
	}
	if n < 1 {
		return r
	}
	r.N = &n
	return r
}

func (r *ChatRequest) AddStop(stop string) *ChatRequest {
	if r.Stop == nil {
		r.Stop = &[]string{}
	}
	*r.Stop = append(*r.Stop, stop)
	return r
}

func (r *ChatRequest) SetStream() *ChatRequest {
	stream := true
	r.Stream = &stream
	return r
}

func (r *ChatRequest) AddMessage(message *Message) *ChatRequest {
	r.Messages = append(r.Messages, message)
	return r
}

func (r *ChatRequest) AddSystem(message string) *ChatRequest {
	r.Messages = append(r.Messages, SystemMessage(message))
	return r
}

func (r *ChatRequest) AddUser(message string) *ChatRequest {
	r.Messages = append(r.Messages, UserMessage(message))
	return r
}

func (r *ChatRequest) AddAssistant(message string) *ChatRequest {
	r.Messages = append(r.Messages, AssistantMessage(message))
	return r
}
