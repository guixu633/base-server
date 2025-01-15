package service

import (
	"context"

	"github.com/guixu633/base-server/module/qdrant"
)

func (s *Service) GetCryptoArticle(ctx context.Context, query string) ([]*qdrant.CryptoArticle, error) {
	return s.qdrant.SearchCryptoArticle(ctx, query)
}

func (s *Service) UpsertCryptoArticle(ctx context.Context, article *qdrant.CryptoArticle) error {
	return s.qdrant.UpsertCryptoArticle(ctx, article)
}
