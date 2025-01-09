package telegram

import (
	"net/http"
	"net/url"
	"time"

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

func GetBot(cfg *config.Telegram, workflow *workflow.Workflow, proxy bool) (*TGBot, error) {
	var client *http.Client
	if proxy {
		proxyUrl, _ := url.Parse("http://127.0.0.1:7890") // 根据你的代理情况修改
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: time.Second * 10,
		}
	} else {
		client = &http.Client{}
	}

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
