package llm

type MessageRole string

const (
	SystemMessageRole    MessageRole = "system"
	UserMessageRole      MessageRole = "user"
	AssistantMessageRole MessageRole = "assistant"
	ToolMessageRole      MessageRole = "tool"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func SystemMessage(prompt string) *Message {
	return &Message{
		Role:    SystemMessageRole,
		Content: prompt,
	}
}

func UserMessage(message string) *Message {
	return &Message{
		Role:    UserMessageRole,
		Content: message,
	}
}

func AssistantMessage(message string) *Message {
	return &Message{
		Role:    AssistantMessageRole,
		Content: message,
	}
}
