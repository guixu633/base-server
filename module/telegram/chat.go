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

		// 添加命令处理逻辑
		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

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

// 添加新的命令处理函数
func (b *TGBot) handleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Command() {
	case "start":
		msg.Text = "欢迎使用本机器人！"
	case "help":
		msg.Text = "这是帮助信息..."
	default:
		msg.Text = "未知命令，请使用 /help 查看可用命令"
	}

	if _, err := b.bot.Send(msg); err != nil {
		logrus.Errorf("发送命令响应消息错误: %v", err)
	}
}
