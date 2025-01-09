package telegram

import (
	"fmt"
	"net/http"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/module/workflow"
	"github.com/stretchr/testify/assert"
)

func TestDemoChat(t *testing.T) {
	bot := getBot(t)
	bot.Response()
}

func TestTelegram(t *testing.T) {
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID

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
	bot, err := GetBot(&cfg.Telegram, workflow, true)
	assert.NoError(t, err)
	bot.bot.Debug = true
	return bot
}
