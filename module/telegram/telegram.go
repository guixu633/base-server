package telegram

import (
	"net/http"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/workflow"
)

type TGBot struct {
	cfg      *config.Telegram
	bot      *tgbotapi.BotAPI
	workflow *workflow.Workflow
	users    map[int64]string
}

// https://go-telegram-bot-api.dev/
// https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5
func GetBot(cfg *config.Telegram, workflow *workflow.Workflow, client *http.Client) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.ApiToken, tgbotapi.APIEndpoint, client)
	if err != nil {
		return nil, err
	}
	return &TGBot{
		cfg:      cfg,
		bot:      bot,
		workflow: workflow,
		users:    make(map[int64]string),
	}, nil
}

func formatTelegramMessage(text string) string {
	// 1. 处理标题的 '#' 符号
	text = strings.ReplaceAll(text, "#", "\\#")

	// 2. 处理其他特殊字符
	special := []string{"_", "[", "]", "(", ")", "~", "`", ">", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range special {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	// 3. 处理粗体
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "*$1*")

	return text
}
