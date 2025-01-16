package coingecko

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func getCoingecko(t *testing.T) *Coingecko {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	coingecko := NewCoingecko(&cfg.Coingecko, client)
	return coingecko
}
