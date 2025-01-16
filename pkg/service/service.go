package service

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/guixu633/base-server/module/coingecko"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/embedding"
	"github.com/guixu633/base-server/module/oss"
	"github.com/guixu633/base-server/module/qdrant"
	"github.com/guixu633/base-server/module/telegram"
	"github.com/guixu633/base-server/module/workflow"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg       *config.Config
	client    *http.Client
	oss       *oss.Oss
	workflow  *workflow.Workflow
	telegram  *telegram.TGBot
	qdrant    *qdrant.Qdrant
	coingecko *coingecko.Coingecko
}

func NewService(cfg *config.Config) (*Service, error) {
	var client *http.Client
	if cfg.Meta.Env == "local" || cfg.Meta.Env == "" {
		proxyUrl, _ := url.Parse("http://127.0.0.1:7890") // 根据你的代理情况修改
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: time.Minute,
		}
	} else {
		client = &http.Client{}
	}

	oss, err := oss.NewOss(&cfg.Oss)
	if err != nil {
		logrus.WithField("err", err).Error("初始化oss失败")
		return nil, err
	}

	workflow := workflow.NewWorkflow(&cfg.Workflow, client)
	bot, err := telegram.GetBot(&cfg.Telegram, workflow, client)
	if err != nil {
		logrus.WithField("err", err).Error("初始化telegram失败")
		return nil, err
	}

	embedEngine := embedding.NewEmbedEngine(&cfg.Embedding, client)

	qdrant, err := qdrant.NewClient(&cfg.Qdrant, &cfg.CryptoArticle, embedEngine)
	if err != nil {
		logrus.WithField("err", err).Error("初始化qdrant失败")
		return nil, err
	}

	err = qdrant.InitVdb(context.Background())
	if err != nil {
		logrus.WithField("err", err).Error("初始化vdb失败")
		return nil, err
	}

	coingecko := coingecko.NewCoingecko(&cfg.Coingecko, client)

	svc := &Service{
		client:    client,
		oss:       oss,
		cfg:       cfg,
		workflow:  workflow,
		telegram:  bot,
		qdrant:    qdrant,
		coingecko: coingecko,
	}

	go bot.Response()

	return svc, nil
}
