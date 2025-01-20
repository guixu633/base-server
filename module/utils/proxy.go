package utils

import (
	"net/http"
	"net/url"
)

func GetProxyClient() *http.Client {
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	return client
}
