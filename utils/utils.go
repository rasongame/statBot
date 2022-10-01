package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func LoadCache(db *gorm.DB, CachedUsers map[int64]CacheUser) {
	var users []User
	db.Find(&users)
	for _, user := range users {
		CachedUsers[user.Id] = CacheUser{
			User:     user,
			LifeTime: time.Now().Unix() + 3600,
		}
	}
}

func UpdateCache(placeholder *SomePlaceholder, db *gorm.DB, CachedUsers map[int64]CacheUser) {
	if CachedUsers[placeholder.User.ID].LifeTime <= time.Now().Unix() {
		u := CacheUser{
			User: User{
				Id:           placeholder.User.ID,
				FirstName:    placeholder.User.FirstName,
				LastName:     placeholder.User.LastName,
				Username:     placeholder.User.UserName,
				LanguageCode: placeholder.User.LanguageCode,
				LastSeen:     placeholder.LastSeenAt,
			},
			LifeTime: time.Now().Unix() + 3600, // 60 min
		}
		CachedUsers[placeholder.User.ID] = u

		db.Save(&u.User)

	}
}
func PanicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func AdminRightChatUpdater(b *tgbotapi.BotAPI) {
	for {
		select {
		case <-AdminRightUpdateTicker.C:
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

func strContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isExplicitAllowedCommand(commandNormalized string) bool {
	commandWithoutParams := strings.Split(commandNormalized, " ")[0]
	commandNormalized = strings.ToLower(commandWithoutParams)[1:]
	return strContains([]string{"stats", "astats", "health", "pop"}, commandNormalized)
}

func RightCommandExtractor(m *tgbotapi.Message, botNickName string) string {
	normalizedBotNickName := strings.ToLower(botNickName)
	if !m.IsCommand() {
		return ""
	}

	// Explicit bot call in public chats
	if m.Chat.Type == "supergroup" {
		splittedCommand := strings.Split(m.Text, "@")

		if len(splittedCommand) > 1 {
			possibleBotNickName := strings.ToLower(strings.Split(splittedCommand[1], " ")[0])
			if possibleBotNickName == normalizedBotNickName {
				return m.Command()
			}
		}
	} else {
		return m.Command()
	}

	if isExplicitAllowedCommand(m.Text) {
		return m.Command()
	}

	return ""
}
