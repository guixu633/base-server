package qdrant

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCollectionInfo(t *testing.T) {
	qdrant := getQdrant(t)
	info, err := qdrant.GetCryptoArticleInfo()
	assert.NoError(t, err)
	// fmt.Println(info)
	fmt.Println("config:", info.GetConfig())
	fmt.Println("payload schema:", info.GetPayloadSchema())
	fmt.Println("points count:", info.GetPointsCount())
	fmt.Println("segments count:", info.GetSegmentsCount())
}

func TestInitCryptoArticleVDB(t *testing.T) {
	qdrant := getQdrant(t)
	err := qdrant.InitCryptoArticleVDB()
	assert.NoError(t, err)
}

func TestDeleteCryptoArticleVDB(t *testing.T) {
	qdrant := getQdrant(t)
	err := qdrant.DeleteCryptoArticleVDB()
	assert.NoError(t, err)
}

func TestUpsertCryptoArticle(t *testing.T) {
	qdrant := getQdrant(t)
	article := &CryptoArticle{
		Id:          uuid.New().String(),
		Title:       "AIOZ突破1.15美元，24小时涨幅达28.0%",
		Content:     "t金色财经报道，行情显示，AIOZ突破1.15美元，现报价1.11美元，24小时涨幅达28.0%，行情波动较大，请做好风险控制。est",
		PublishTime: "2025-01-06T00:18:01+08:00",
		Src:         "inf_crawl_knowledge_jinsecaijing",
		Url:         "https://www.jinse.cn/lives/435348.html",
	}
	err := qdrant.UpsertCryptoArticle(context.Background(), article)
	assert.NoError(t, err)

	article = &CryptoArticle{
		Id:          uuid.New().String(),
		Title:       "Morgan Stanley Mulls Crypto Trading Launch Through E-Trade Platform: Report",
		Content:     "Morgan Stanley Mulls Crypto Trading Launch Through E-Trade Platform: Report Morgan Stanley Morgan Stanley's addition of crypto trading to E-Trade could expand retail crypto access for its 5.2M clients, positioning E-Trade as a major competitor to Coinbase. Last updated: January 2, 2025 23:32 EST Crypto Reporter Shalini Nagarajan Crypto Reporter Shalini Nagarajan About Author Shalini is a crypto reporter who provides in-depth reports on daily developments and regulatory shifts in the cryptocurrency sector. Author Profile Share Copied Last updated: January 2, 2025 23:32 EST Why Trust Cryptonews Cryptonews has covered the cryptocurrency industry topics since 2017, aiming to provide informative insights to our readers. Our journalists and analysts have extensive experience in market analysis and blockchain technologies. We strive to maintain high editorial standards, focusing on factual accuracy and balanced reporting across all areas - from cryptocurrencies and blockchain projects to industry events, products, and technological developments. Our ongoing presence in the industry reflects our commitment to delivering relevant information in the evolving worl American multinational investment bank Morgan Stanley is reportedly exploring launching crypto trading services through its E-Trade platform. On Jan. 2, The Information reported that the bank’s decision to launch this service is based on the second Trump administration’s promise of a more favorable regulatory environment for digital assets.To support this initiative, the bank is preparing to offer E-Trade customers direct access to crypto trading. It plans to start with major coins like Bitcoin and Ethereum. However, to proceed, the bank must secure regulatory approval from the Federal Reserve. This is required because of its classification as a bank holding company.Morgan Stanley didn’t return Cryptonews’ request for comment by press time.Trump Promises Crypto-Friendly Reforms to Cement US as a Global LeaderMeanwhile, Trump, during his presidential campaign, made bold promises to the crypto industry, with the goal of making the US the “world’s crypto capital.” His approach involves placing pro-crypto individuals in leadership roles at key regulatory agencies.He’s also pledged to foster a more lenient regulatory climate for digital currencies. This strategy includes potentially giving more regulatory control to the CFTC over digital assets, which would mean moving some oversight away from the SEC to the CFTC, known for being more crypto-friendly.E-Trade’s Crypto Trading Ambitions Could Challenge Coinbase’s DominanceMorgan Stanley completed its acquisition of E-Trade in late 2020 through an all-stock deal valued at $13b. This acquisition expanded Morgan Stanley’s wealth management operations by merging with E-Trade’s retail customer base and digital platform.By adding crypto trading to E-Trade’s platform, Morgan Stanley could bring digital assets to a broader audience. This includes E-Trade’s 5.2m clients managing $360b in assets. Such a move could significantly increase retail participation in the crypto market and position E-Trade as a leading traditional brokerage offering crypto services, directly challenging Coinbase.Consequently, this would heighten competition with existing crypto platforms like Coinbase. E-Trade could potentially capture market share by offering competitive pricing or seamlessly integrating crypto with its current financial services.Notably, last year, Morgan Stanley allowed its financial advisors to offer spot Bitcoin ETFs to select clients, becoming a pioneer among major Wall Street banks. The bank reportedly monitored clients’ crypto investments closely to limit exposure to the volatile asset class. This approach helps maintain balanced portfolios. Follow us on Google News Trending News Price PredictionsRecommended Articles Trump Appoints PayPal Veteran David Sacks as ‘White House AI and Crypto Czar’Analyst Predicts XRP Could Hit $8-$20 This Cycle – Could $100 Be Next?Ethereum Price Explosion Imminent, This Historic Price Pattern PredictsMorgan Stanley Mulls Crypto Trading Launch Through E-Trade Platform: ReportDogecoin Price Hits ‘Rock Solid’ Support – $1 Surge Incoming? Altcoin News Dogecoin Price to $3 In 2025? Here’s Why Its Very Possible 2025-01-03 00:14:22, by Joel Frank Price Analysis Top Crypto Coins to Buy in 2025 – VIRTUAL, PEPE, PENGU 2025-01-01 16:05:52, by Alejandro Arrieche Price Analysis Elon Musk’s Pepe Meme Tweet Turns $1,000 Into a Fortune – Could PEPE 100x From Here? 2025-01-01 16:15:54, by Arslan Butt Price Analysis Kekius Maximus Explodes 21,660% After Elon Musk’s X Post 2025-01-01 17:49:04, by Alejandro Arrieche Bitcoin (BTC) Price PredictionEthereum (ETH) Price PredictionRipple (XRP) Price PredictionDogecoin (DOGE) Price PredictionSolana (SOL) Price Prediction Best Crypto WalletsBest Crypto to Buy NowBest Crypto Presales to Invest InBest New Meme Coins to Buy ",
		PublishTime: "2025-01-06T00:29:23+08:00",
		Src:         "inf_crawl_knowledge_cryptonew",
		Url:         "https://cryptonews.com/news/morgan-stanley-mulls-crypto-trading-launch-e-trade/",
	}
	err = qdrant.UpsertCryptoArticle(context.Background(), article)
	assert.NoError(t, err)
}

func TestParseTime(t *testing.T) {
	timestamp := 1736094563
	date := time.Unix(int64(timestamp), 0).Format(time.RFC3339)
	fmt.Println(date)
}

func TestSearchCryptoArticle(t *testing.T) {
	qdrant := getQdrant(t)
	start := time.Now()
	results, err := qdrant.SearchCryptoArticle(context.Background(), "AIOZ突破1.15美元，24小时涨幅达28.0%", 10, 0.5)
	assert.NoError(t, err)
	for _, result := range results {
		json, err := json.MarshalIndent(result, "", "  ")
		assert.NoError(t, err)
		fmt.Println(string(json))
	}
	fmt.Println("total:", len(results))
	fmt.Println("time:", time.Since(start))
}

func TestSearchCryptoArticleStr(t *testing.T) {
	qdrant := getQdrant(t)
	start := time.Now()
	result, err := qdrant.SearchCryptoArticleStr(context.Background(), "AIOZ突破1.15美元，24小时涨幅达28.0%", 10, 0.5)
	assert.NoError(t, err)
	fmt.Println(result)
	fmt.Println("time:", time.Since(start))
}
