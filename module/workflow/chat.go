package workflow

import "context"

func (w *Workflow) DemoChat(ctx context.Context, query, conversationId string) (string, string, error) {
	token := "app-98ehTdEHEYn45dgKMuIOF6bf"
	inputs := map[string]string{}
	return w.CallChatflowBlock(ctx, token, query, conversationId, inputs)
}
