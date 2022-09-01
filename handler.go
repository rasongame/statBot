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
func SuperuserFilter(b *tgbotapi.BotAPI, m *tgbotapi.Message) bool {
	if m.From.ID == Superuser {
		return true
	} else {
		return false
	}
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
