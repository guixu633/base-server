package service

import (
	"context"

	"github.com/guixu633/base-server/module/qdrant"
)

func (s *Service) GetCryptoArticle(ctx context.Context, query string, topk int, scoreThreshold float64) ([]*qdrant.CryptoArticle, error) {
	return s.qdrant.SearchCryptoArticle(ctx, query, topk, scoreThreshold)
}

func (s *Service) UpsertCryptoArticle(ctx context.Context, article *qdrant.CryptoArticle) error {
	return s.qdrant.UpsertCryptoArticle(ctx, article)
}

func (s *Service) Retrieval(ctx context.Context, req qdrant.RetrievalRequest) (qdrant.RetrievalResponse, error) {
	return s.qdrant.Retrieval(ctx, req)
}
