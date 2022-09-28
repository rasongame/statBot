package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"statBot/utils"
)

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
	if utils.AdminRightsCache[message.Chat.ID] == nil {
		chatCfg := tgbotapi.ChatConfig{ChatID: message.Chat.ID}
		res, _ := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{chatCfg})
		if utils.AdminRightsCache[message.Chat.ID] == nil {
			utils.AdminRightsCache[message.Chat.ID] = map[int64]tgbotapi.ChatMember{}
		}
		for _, admin := range res {
			utils.AdminRightsCache[message.Chat.ID][admin.User.ID] = admin
		}
		go utils.AdminRightChatUpdater(bot)
	}
	isAdmin := utils.AdminRightsCache[message.Chat.ID][message.From.ID].CanDeleteMessages

	if isAdmin || message.From.ID == utils.Superuser {
		return true
	}
	return false
}
