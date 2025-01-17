package telegram

import (
	"context"
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

		// 检查消息是否是群聊中的 @ 消息或私聊消息
		if !b.shouldProcessMessage(update.Message) {
			continue
		}

		chatID := update.Message.Chat.ID

		// 清理消息文本（去除 @ 部分）
		messageText := b.cleanMessageText(update.Message)

		// 检查清理后的消息是否是命令
		if strings.HasPrefix(messageText, "/") {
			// 创建一个新的 Message 对象，包含清理后的文本
			cmdMessage := *update.Message
			cmdMessage.Text = messageText
			b.handleCommand(&cmdMessage)
			continue
		}

		convID := b.users[chatID]
		convID, answer, err := b.workflow.DemoChat(context.Background(), messageText, convID)
		if err != nil {
			logrus.Errorf("get answer error: %v", err)
			continue
		}
		b.users[chatID] = convID

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatTelegramMessage(answer))
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		msg.ReplyToMessageID = update.Message.MessageID

		// msg.ReplyToMessageID = update.Message.MessageID

		if _, err := b.bot.Send(msg); err != nil {
			logrus.Errorf("send message error: %v", err)
			continue
		}
	}
}

// 添加新的命令处理函数
func (b *TGBot) handleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Text {
	case "/start":
		msg.Text = "欢迎使用本机器人！"
	case "/help":
		msg.Text = "这是帮助信息..."
	default:
		msg.Text = "未知命令，请使用 /help 查看可用命令"
	}
	msg.ReplyToMessageID = message.MessageID
	if _, err := b.bot.Send(msg); err != nil {
		logrus.Errorf("发送命令响应消息错误: %v", err)
	}
}

// 添加新的辅助方法
func (b *TGBot) shouldProcessMessage(message *tgbotapi.Message) bool {
	// 如果是私聊消息，直接处理
	if message.Chat.Type == "private" {
		return true
	}

	// 如果是群聊消息，检查是否 @ 了机器人
	if message.Chat.Type == "group" || message.Chat.Type == "supergroup" {
		return message.Text != "" && b.isBotMentioned(message)
	}

	return false
}

func (b *TGBot) isBotMentioned(message *tgbotapi.Message) bool {
	if message.Entities == nil {
		return false
	}

	botUsername := b.bot.Self.UserName
	for _, entity := range message.Entities {
		if entity.Type == "mention" && message.Text[entity.Offset+1:entity.Offset+entity.Length] == botUsername {
			return true
		}
	}
	return false
}

func (b *TGBot) cleanMessageText(message *tgbotapi.Message) string {
	if message.Chat.Type == "private" {
		return message.Text
	}

	// 移除 @ 机器人的用户名部分
	text := message.Text
	if message.Entities != nil {
		for _, entity := range message.Entities {
			if entity.Type == "mention" {
				text = strings.TrimSpace(text[entity.Offset+entity.Length:])
				break
			}
		}
	}
	return text
}
