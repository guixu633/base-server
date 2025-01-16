package coingecko

import (
	"net/http"

	"github.com/guixu633/base-server/module/config"
)

type Coingecko struct {
	cfg    *config.Coingecko
	client *http.Client
}

func NewCoingecko(cfg *config.Coingecko, client *http.Client) *Coingecko {
	return &Coingecko{
		cfg:    cfg,
		client: client,
	}
}
