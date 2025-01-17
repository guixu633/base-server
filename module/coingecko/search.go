package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CurrentInfoResponse struct {
	MarketData   MarketDataResponse `json:"market_data"`
	Localization map[string]string  `json:"localization"`
}

type MarketDataResponse struct {
	CurrentPrice                 map[string]float64 `json:"current_price"`
	AllTimeHigh                  map[string]float64 `json:"ath"`
	AllTimeHighDate              map[string]string  `json:"ath_date"`
	AllTimeLow                   map[string]float64 `json:"atl"`
	AllTimeLowDate               map[string]string  `json:"atl_date"`
	MarketCap                    map[string]float64 `json:"market_cap"`
	TotalVolume                  map[string]float64 `json:"total_volume"`
	High24h                      map[string]float64 `json:"high_24h"`
	Low24h                       map[string]float64 `json:"low_24h"`
	PriceChange24h               float64            `json:"price_change_24h"`
	PriceChangePercentage24h     float64            `json:"price_change_percentage_24h"`
	PriceChangePercentage7d      float64            `json:"price_change_percentage_7d"`
	PriceChangePercentage14d     float64            `json:"price_change_percentage_14d"`
	PriceChangePercentage30d     float64            `json:"price_change_percentage_30d"`
	PriceChangePercentage60d     float64            `json:"price_change_percentage_60d"`
	PriceChangePercentage200d    float64            `json:"price_change_percentage_200d"`
	PriceChangePercentage1y      float64            `json:"price_change_percentage_1y"`
	MarketCapChange24h           float64            `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64            `json:"market_cap_change_percentage_24h"`
}

type MarketDate struct {
	Name                         string  `json:"coin_name"`
	CurrentPrice                 float64 `json:"current_price"`
	AllTimeHigh                  float64 `json:"all_time_high"`
	AllTimeHighDate              string  `json:"all_time_high_date"`
	AllTimeLow                   float64 `json:"all_time_low"`
	AllTimeLowDate               string  `json:"all_time_low_date"`
	MarketCap                    float64 `json:"market_cap"`
	TotalVolume                  float64 `json:"total_volume"`
	High24h                      float64 `json:"high_24h"`
	Low24h                       float64 `json:"low_24h"`
	PriceChange24h               float64 `json:"price_change_24h"`
	PriceChangePercentage24h     float64 `json:"price_change_percentage_24h"`
	PriceChangePercentage7d      float64 `json:"price_change_percentage_7d"`
	PriceChangePercentage14d     float64 `json:"price_change_percentage_14d"`
	PriceChangePercentage30d     float64 `json:"price_change_percentage_30d"`
	PriceChangePercentage60d     float64 `json:"price_change_percentage_60d"`
	PriceChangePercentage200d    float64 `json:"price_change_percentage_200d"`
	PriceChangePercentage1y      float64 `json:"price_change_percentage_1y"`
	MarketCapChange24h           float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
}

func (i CurrentInfoResponse) ToMarketDate() *MarketDate {
	market := i.MarketData
	return &MarketDate{
		Name:                         i.Localization["zh"],
		CurrentPrice:                 market.CurrentPrice["usd"],
		AllTimeHigh:                  market.AllTimeHigh["usd"],
		AllTimeHighDate:              market.AllTimeHighDate["usd"],
		AllTimeLow:                   market.AllTimeLow["usd"],
		AllTimeLowDate:               market.AllTimeLowDate["usd"],
		MarketCap:                    market.MarketCap["usd"],
		TotalVolume:                  market.TotalVolume["usd"],
		High24h:                      market.High24h["usd"],
		Low24h:                       market.Low24h["usd"],
		PriceChange24h:               market.PriceChange24h,
		PriceChangePercentage24h:     market.PriceChangePercentage24h,
		PriceChangePercentage7d:      market.PriceChangePercentage7d,
		PriceChangePercentage14d:     market.PriceChangePercentage14d,
		PriceChangePercentage30d:     market.PriceChangePercentage30d,
		PriceChangePercentage60d:     market.PriceChangePercentage60d,
		PriceChangePercentage200d:    market.PriceChangePercentage200d,
		PriceChangePercentage1y:      market.PriceChangePercentage1y,
		MarketCapChange24h:           market.MarketCapChange24h,
		MarketCapChangePercentage24h: market.MarketCapChangePercentage24h,
	}
}

func (c *Coingecko) Search(ctx context.Context, coinID string) (*MarketDate, error) {
	// 构建请求 URL
	url := fmt.Sprintf("%s/coins/%s?developer_data=false&community_data=false&tickers=false",
		c.cfg.Url, coinID)

	// 创建新的 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加 API key 到请求头
	req.Header.Set("x-cg-pro-api-key", c.cfg.ApiKey)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var response CurrentInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.ToMarketDate(), nil
}
