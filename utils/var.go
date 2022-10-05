package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	UpdatesProcessed       int64
	Handlers               map[string]Handler
	CachedUsers            map[int64]CacheUser
	ChatLogIsLoaded        map[int64]bool
	ChatLogMessageCache    map[int64]map[int64]*SomePlaceholder
	AdminRightsCache       map[int64]map[int64]tgbotapi.ChatMember
	AdminRightUpdateTicker = time.NewTicker(15 * time.Minute)
	CachedUsersLifeTime    = int64(10) // in seconds
)

func init() {
	AdminRightsCache = map[int64]map[int64]tgbotapi.ChatMember{}
}
