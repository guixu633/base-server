package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestCoinAnalysis(t *testing.T) {
	db := getDB(t)
	analysis, err := db.GetLatestCoinAnalysis(context.Background(), "BTC")
	assert.NoError(t, err)
	fmt.Println(analysis)
}

func TestGetAllCoinAnalysis(t *testing.T) {
	db := getDB(t)
	analysis, err := db.GetAllCoinAnalysis(context.Background())
	assert.NoError(t, err)
	for _, a := range analysis {
		fmt.Println(a.CoinID)
		fmt.Println(a.CreatedAt)
		fmt.Println()
	}

	fmt.Println(len(analysis))
}
