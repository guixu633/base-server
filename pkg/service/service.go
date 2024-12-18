package service

import (
	"net/http"
	"time"

	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/oss"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg    *config.Config
	client *http.Client
	oss    *oss.Oss
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

	svc := &Service{
		client: client,
		oss:    oss,
		cfg:    cfg,
	}
	return svc
}
