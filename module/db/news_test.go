package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListNews(t *testing.T) {
	db := getDB(t)
	news, err := db.ListNews(context.Background())
	assert.NoError(t, err)
	for _, n := range news {
		fmt.Println(n.CreatedAt)
	}
	fmt.Println(len(news))
}

func TestGetLatestNews(t *testing.T) {
	db := getDB(t)
	news, err := db.GetLatestNews(context.Background())
	assert.NoError(t, err)
	fmt.Println(news)
}
