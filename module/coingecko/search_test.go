package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	coin := getCoingecko(t)
	result, err := coin.Search(context.Background(), "bitcoin")
	assert.NoError(t, err)
	json, err := json.MarshalIndent(result, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(json))
}
