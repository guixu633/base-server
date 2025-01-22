package workflow

import "context"

func (w *Workflow) DemoChat(ctx context.Context, query, conversationId string) (string, string, error) {
	token := "app-cR1gm0hNz2cEgUBTwl9QPyqh"
	inputs := map[string]string{}
	return w.CallChatflowBlock(ctx, token, query, conversationId, inputs)
}
