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

func TestGetName(t *testing.T) {
	coin := getCoingecko(t)
	idList := []string{
		"bitcoin",
		"ethereum",
		"ripple",
		"binancecoin",
		"solana",
		"dogecoin",
		"cardano",
		"tron",
		"avalanche-2",
		"sui",
		"the-open-network",
		"chainlink",
		"stellar",
		"shiba-inu",
		"hedera-hashgraph",
		"polkadot",
		"bitcoin-cash",
		"uniswap",
		"pepe",
		"litecoin",
		"near",
		"internet-computer",
		"aave",
		"ethereum-classic",
		"vechain",
		"fetch-ai",
		"arbitrum",
		"filecoin",
		"algorand",
		"cosmos",
		"bonk",
		"injective-protocol",
		"theta-token",
		"the-graph",
		"fantom",
		"dogwifcoin",
		"sei-network",
		"floki",
		"gala",
		"the-sandbox",
		"thorchain",
		"jupiter",
		"tezos",
		"eos",
		"iota",
		"flow",
		"bittorrent",
		"axie-infinity",
		"decentraland",
		"elrond-erd-2",
		"apecoin",
		"zcash",
		"eigenlayer",
		"notcoin",
		"neiro",
		"terra-luna-2",
		"og-fan-token",
	}
	name := make(map[string]string)
	for _, id := range idList {
		result, err := coin.Search(context.Background(), id)
		assert.NoError(t, err)
		name[id] = result.Name
	}
	json, err := json.MarshalIndent(name, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(json))
}
