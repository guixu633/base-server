package workflow

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatflow(t *testing.T) {
	w := getWorkflow(t)
	conversationId, answer, err := w.DemoChat(context.Background(), "我孩子今年多大了？", "56155327-fbb1-4ba3-8746-d22f449ab55d")
	assert.NoError(t, err)
	fmt.Println(conversationId)
	fmt.Println(answer)
}
