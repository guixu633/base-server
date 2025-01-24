package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (b *TGBot) Response() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// 处理回调查询
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update)
			continue
		}

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

const (
	StartResponse = `Welcome to TokenSense! 🚀
The most reliable crypto tools before you enter the market.
Join our community:
 📱 Twitter: https://x.com/TokenSenseAI
 💬 Telegram Group: https://t.me/tokensense01
=================`
	HelpResponse = `📖 Help Center - Quick Guide
Find detailed documentation at:
https://github.com/Tokensense-ai/Tokensense
For more assistance:
 • Visit our docs for complete features & tutorials
 • Join our Telegram community for real-time support：https://t.me/tokensense01
 • Follow us on Twitter for updates & tips：: https://x.com/TokenSenseAI`
)

// 修改命令处理函数
func (b *TGBot) handleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Text {
	case "/start":
		msg.Text = StartResponse
		// 创建按钮键盘
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("BTC", "/predict BTC"),
				tgbotapi.NewInlineKeyboardButtonData("ETH", "/predict ETH"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("TRUMP", "/predict TRUMP"),
				tgbotapi.NewInlineKeyboardButtonData("SOL", "/predict SOL"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❓ Help", "/help"),
				tgbotapi.NewInlineKeyboardButtonData("🌍 News", "/news"),
			),
		)
		msg.ReplyMarkup = keyboard
	case "/help":
		msg.Text = HelpResponse
	case "/menu":
		msg.Text = StartResponse
	case "/news":
		news, err := b.db.GetLatestNews(context.Background())
		if err != nil {
			msg.Text = "获取新闻失败"
		}
		msg.Text = news
	default:
		if strings.HasPrefix(message.Text, "/predict") {
			coin := strings.TrimPrefix(message.Text, "/predict ")
			msg.Text = fmt.Sprintf("对 [%s] 的分析即将上线", coin)
			analysis, err := b.db.GetLatestCoinAnalysis(context.Background(), coin)
			if err != nil {
				msg.Text = fmt.Sprintf("没有查询到币种 [%s] 的分析结果", coin)
			}
			msg.Text = analysis
		} else {
			msg.Text = "未知命令，请使用 /help 查看可用命令"

		}
	}
	msg.ReplyToMessageID = message.MessageID
	if _, err := b.bot.Send(msg); err != nil {
		logrus.Errorf("发送命令响应消息错误: %v", err)
	}
}

// 添加处理回调查询的函数
func (b *TGBot) handleCallbackQuery(update tgbotapi.Update) {
	callback := update.CallbackQuery
	// 创建一个新的消息对象，模拟用户发送命令
	cmdMessage := tgbotapi.Message{
		MessageID: callback.Message.MessageID,
		From:      callback.From,
		Chat:      callback.Message.Chat,
		Text:      callback.Data,
	}

	// 处理命令
	b.handleCommand(&cmdMessage)

	// 回应回调查询
	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	b.bot.Request(callbackResponse)
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
