package qdrant

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrievalCryptoArticle(t *testing.T) {
	q := getQdrant(t)
	req := RetrievalRequest{
		KnowledgeID: "crypto_article",
		Query:       "比特币",
		Settings:    RetrievalSettings{TopK: 10, ScoreThreshold: 0.5},
	}
	resp, err := q.Retrieval(context.Background(), req)
	assert.NoError(t, err)
	for _, record := range resp.Records {
		data, err := json.MarshalIndent(record, "", "  ")
		assert.NoError(t, err)
		fmt.Println(string(data))
	}
}
