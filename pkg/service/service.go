package service

import (
	"net/http"
	"time"

	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/oss"
	"github.com/guixu633/base-server/module/telegram"
	"github.com/guixu633/base-server/module/workflow"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg      *config.Config
	client   *http.Client
	oss      *oss.Oss
	workflow *workflow.Workflow
	telegram *telegram.TGBot
}

func NewService(cfg *config.Config) *Service {
	client := &http.Client{
		Timeout: time.Minute,
	}

	oss, err := oss.NewOss(&cfg.Oss)
	if err != nil {
		logrus.WithField("err", err).Error("初始化oss失败")
		panic(err)
	}

	workflow := workflow.NewWorkflow(&cfg.Workflow, client)
	bot, err := telegram.GetBot(&cfg.Telegram, workflow, true)
	if err != nil {
		logrus.WithField("err", err).Error("初始化telegram失败")
		panic(err)
	}

	svc := &Service{
		client:   client,
		oss:      oss,
		cfg:      cfg,
		workflow: workflow,
		telegram: bot,
	}

	go bot.Response()

	return svc
}
