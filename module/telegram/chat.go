package telegram

import (
	"context"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (b *TGBot) Response() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		convID := b.users[chatID]
		convID, answer, err := b.workflow.DemoChat(context.Background(), update.Message.Text, convID)
		if err != nil {
			logrus.Errorf("get answer error: %v", err)
			continue
		}
		b.users[chatID] = convID

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatTelegramMessage(answer))
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		// msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.bot.Send(msg); err != nil {
			logrus.Errorf("send message error: %v", err)
			continue
		}
	}
}

// func formatTelegramMessage(text string) string {
// 	// 1. 将标题转换为加粗文本
// 	// text = regexp.MustCompile(`### (.+)`).ReplaceAllString(text, "*$1*")

// 	// 2. 处理粗体 ** 到单个 *
// 	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "*$1*")

// 	// 3. 处理特殊字符
// 	special := []string{"_", "[", "]", "(", ")", "~", "`", ">", "+", "-", "=", "|", "{", "}", ".", "!"}
// 	for _, char := range special {
// 		text = strings.ReplaceAll(text, char, "\\"+char)
// 	}

// 	// 4. 处理百分比符号
// 	text = strings.ReplaceAll(text, "%", "\\%")

// 	return text
// }

func formatTelegramMessage(text string) string {
	// 1. 处理标题的 '#' 符号
	text = strings.ReplaceAll(text, "#", "\\#")

	// 2. 处理其他特殊字符
	special := []string{"_", "[", "]", "(", ")", "~", "`", ">", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range special {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	// 3. 处理粗体
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "*$1*")

	return text
}
