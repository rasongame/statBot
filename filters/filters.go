package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"statBot/utils"
	"time"
)

var (
	AdminRightsCache       map[int64]map[int64]tgbotapi.ChatMember
	adminRightUpdateTicker = time.NewTicker(15 * time.Minute)
)

func init() {
	AdminRightsCache = map[int64]map[int64]tgbotapi.ChatMember{}
}
func AdminRightChatUpdater(b *tgbotapi.BotAPI) {
	for {
		select {
		case <-adminRightUpdateTicker.C:
			{
				for chatId, _ := range AdminRightsCache {
					chatCfg := tgbotapi.ChatConfig{ChatID: chatId}
					res, _ := b.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{chatCfg})
					for _, admin := range res {
						AdminRightsCache[chatId][admin.User.ID] = admin
					}
				}
			}
		}
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
func IsAdminFilter(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if AdminRightsCache[message.Chat.ID] == nil {
		chatCfg := tgbotapi.ChatConfig{ChatID: message.Chat.ID}
		res, _ := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{chatCfg})
		if AdminRightsCache[message.Chat.ID] == nil {
			AdminRightsCache[message.Chat.ID] = map[int64]tgbotapi.ChatMember{}
		}
		for _, admin := range res {
			AdminRightsCache[message.Chat.ID][admin.User.ID] = admin
		}
		go AdminRightChatUpdater(bot)
	}
	isAdmin := AdminRightsCache[message.Chat.ID][message.From.ID].CanDeleteMessages

	if isAdmin || message.From.ID == utils.Superuser {
		return true
	}
	return false
}
