package service

import (
	"context"

	"github.com/guixu633/base-server/module/coingecko"
)

func (s *Service) GetCoinData(ctx context.Context, coinID string) (*coingecko.MarketDate, error) {
	coin, err := s.coingecko.Search(ctx, coinID)
	if err != nil {
		return nil, err
	}
	return coin, nil
}
