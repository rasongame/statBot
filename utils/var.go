package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	UpdatesProcessed       int64
	Handlers               map[string]Handler
	AdminRightsCache       map[int64]map[int64]tgbotapi.ChatMember
	AdminRightUpdateTicker = time.NewTicker(15 * time.Minute)
	Aliases                = map[string]int64{
		"flood": LinFloodID,
		"help":  -1001053617676,
	}
	AllowedChats = map[int64]bool{}
)

func init() {
	AdminRightsCache = map[int64]map[int64]tgbotapi.ChatMember{}
}
