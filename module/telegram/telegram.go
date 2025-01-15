package telegram

import (
	"net/http"

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
