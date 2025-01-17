package telegram

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/workflow"
	"github.com/stretchr/testify/assert"
)

func TestDemoChat(t *testing.T) {
	bot := getBot(t)
	bot.Response()
}

func TestTelegramV2(t *testing.T) {
	rawAnswer := `### 1. 信息概要
- **比特币价格突破10万美元**：近期，比特币价格多次突破10万美元大关，最高达到102,000美元。市场情绪有所回暖，但山寨币整体表现疲软。
- **CPI数据影响**：美国12月CPI数据低于预期，缓解了通胀担忧，推动比特币价格上涨。核心CPI年率为3.2%，低于预期的3.3%。
- **期权市场和ETF资金流入**：短期隐含波动率上升，做多力量增强。比特币现货ETF净流入7.5479亿美元，显示机构需求强劲。
- **特朗普即将上任**：市场对特朗普的新政策充满期待，尤其是关于加密货币的相关政策。预计新政府可能推出利好加密货币的行政令。
- **市场风险**：尽管市场情绪好转，但投资者仍需关注通胀压力的持续性以及即将到来的美国零售销售数据可能带来的波动。

### 2. 问题分析
- **比特币价格反弹原因**：
  - **CPI数据低于预期**：缓解了市场对通胀的担忧，提振了投资者信心。
  - **特朗普政策预期**：市场对特朗普即将上任后可能推出的加密货币利好政策充满期待。
  - **机构资金流入**：比特币现货ETF的大量净流入表明机构投资者对市场的积极态度。
- **潜在风险和不确定性**：
  - **通胀压力持续性**：尽管CPI数据低于预期，但仍需关注未来几个月的通胀情况。
  - **政策变化**：特朗普的具体政策细节尚未公布，存在不确定性。
  - **市场波动**：即将到来的美国零售销售数据可能引发市场波动。

### 3. 建议与观点
- **保持警惕**：尽管市场情绪回暖，但投资者应密切关注宏观经济数据和政策变化，特别是美国零售销售数据和特朗普的具体政策。
- **适度配置短期期权**：可以考虑适度配置短期期权进行战术性交易，以应对市场波动。
- **多元化投资**：在当前市场环境下，建议投资者不要过度集中于单一资产，而是进行多元化投资，分散风险。
- **长期视角**：对于长期投资者来说，比特币的基本面仍然稳固，可以继续持有并关注长期增长潜力。

请注意，以上分析基于现有资讯，市场情况可能会迅速变化，投资者应谨慎决策，并做好风险管理。`

	answer := formatTelegramMessage(rawAnswer)
	fmt.Println(answer)

	bot := getBot(t)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}
	}
}

func TestTelegram(t *testing.T) {
	answer := `## 1. 信息概要\n- **比特币价格突破10万美元**：近期，比特币价格多次突破10万美元大关，最高达到102,000美元。市场情绪有所回暖，但山寨币整体表现疲软。\n- **CPI数据影响**：美国12月CPI数据低于预期，缓解了通胀担忧，推动比特币价格上涨。核心CPI年率为3.2%，低于预期的3.3%。\n- **期权市场和ETF资金流入**：短期隐含波动率上升，做多力量增强。比特币现货ETF净流入7.5479亿美元，显示机构需求强劲。\n- **特朗普即将上任**：市场对特朗普的新政策充满期待，尤其是关于加密货币的相关政策。预计新政府可能推出利好加密货币的行政令。\n- **市场风险**：尽管市场情绪好转，但投资者仍需关注通胀压力的持续性以及即将到来的美国零售销售数据可能带来的波动。\n\n### 2. 问题分析\n- **比特币价格反弹原因**：\n  - **CPI数据低于预期**：缓解了市场对通胀的担忧，提振了投资者信心。\n  - **特朗普政策预期**：市场对特朗普即将上任后可能推出的加密货币利好政策充满期待。\n  - **机构资金流入**：比特币现货ETF的大量净流入表明机构投资者对市场的积极态度。\n- **潜在风险和不确定性**：\n  - **通胀压力持续性**：尽管CPI数据低于预期，但仍需关注未来几个月的通胀情况。\n  - **政策变化**：特朗普的具体政策细节尚未公布，存在不确定性。\n  - **市场波动**：即将到来的美国零售销售数据可能引发市场波动。\n\n### 3. 建议与观点\n- **保持警惕**：尽管市场情绪回暖，但投资者应密切关注宏观经济数据和政策变化，特别是美国零售销售数据和特朗普的具体政策。\n- **适度配置短期期权**：可以考虑适度配置短期期权进行战术性交易，以应对市场波动。\n- **多元化投资**：在当前市场环境下，建议投资者不要过度集中于单一资产，而是进行多元化投资，分散风险。\n- **长期视角**：对于长期投资者来说，比特币的基本面仍然稳固，可以继续持有并关注长期增长潜力。\n\n请注意，以上分析基于现有资讯，市场情况可能会迅速变化，投资者应谨慎决策，并做好风险管理。`
	answer = formatTelegramMessage(answer)
	bot := getBot(t)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		fmt.Println("get a update")
		fmt.Println(update)
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}

func getBot(t *testing.T) *TGBot {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	workflow := workflow.NewWorkflow(&cfg.Workflow, &http.Client{})

	proxyUrl, _ := url.Parse("http://127.0.0.1:7890") // 根据你的代理情况修改
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: time.Second * 10,
	}

	bot, err := GetBot(&cfg.Telegram, workflow, client)
	assert.NoError(t, err)
	bot.bot.Debug = true
	return bot
}
