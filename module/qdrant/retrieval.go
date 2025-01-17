package qdrant

import (
	"context"
	"fmt"
)

type RetrievalRequest struct {
	KnowledgeID string            `json:"knowledge_id"`
	Query       string            `json:"query"`
	Settings    RetrievalSettings `json:"settings"`
}

type RetrievalSettings struct {
	TopK           int     `json:"top_k"`
	ScoreThreshold float64 `json:"score_threshold"`
}

type RetrievalResponse struct {
	Records []RetrievalRecord `json:"records"`
}

type RetrievalRecord struct {
	Content  string         `json:"content"`
	Score    float64        `json:"score"`
	Title    string         `json:"title"`
	Metadata map[string]any `json:"metadata"`
}

func (q *Qdrant) Retrieval(ctx context.Context, req RetrievalRequest) (RetrievalResponse, error) {
	switch req.KnowledgeID {
	case "crypto_article":
		return q.RetrievalCryptoArticle(ctx, req)
	default:
		return RetrievalResponse{}, fmt.Errorf("unknown knowledge id: %s", req.KnowledgeID)
	}
}
