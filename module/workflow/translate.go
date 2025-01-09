package workflow

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Workflow) Translate(ctx context.Context, content string) (string, error) {
	inputs := map[string]string{"content": content}
	resp, err := w.CallWorkflowBlock(ctx, "app-4iCR6YVi7bS1sEmvFN06fvJj", inputs)
	if err != nil {
		return "", err
	}

	result, ok := resp["content"].(string)
	if !ok {
		return "", errors.New("content is not a string")
	}
	logrus.WithField("content", content).WithField("result", result).Info(ctx, "translate")
	return result, nil
}

func (w *Workflow) TranslateRetry(ctx context.Context, content string) (string, error) {
	for i := 0; i < 3; i++ {
		result, err := w.Translate(ctx, content)
		if err == nil {
			return result, nil
		}
		time.Sleep(time.Duration(2+time.Now().UnixNano()%4) * time.Second)
	}
	return "", errors.New("failed to get result")
}
