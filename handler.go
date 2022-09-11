package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func AddHandler(command string, handler HandlerFunc, filter FilterFunc) (Handler, bool) {
	if h, ok := Handlers[command]; ok {
		if !ok {
			return h, false
		}
	}
	h := Handler{handler, filter}
	Handlers[command] = h
	return h, true
}

func TrueFilter(b *tgbotapi.BotAPI, m *tgbotapi.Message) bool {
	return true
}
func FalseFilter(b *tgbotapi.BotAPI, m *tgbotapi.Message) bool {
	return false
}
func ChatOnly(b *tgbotapi.BotAPI, m *tgbotapi.Message) bool {
	return m.Chat.IsGroup() || m.Chat.IsSuperGroup()
}

func IsAdminFilter(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	cfg := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: LinFloodID,
			UserID: message.From.ID,
		},
	}
	res, _ := bot.GetChatMember(cfg)
	isAdmin := res.CanDeleteMessages
	if isAdmin || message.From.ID == Superuser {
		return true
	}
	return false
}
