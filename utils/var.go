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
	//AllowedChats = map[int64]bool{
	//	559723688:      true, // rasongame
	//	-1001549183364: true, // Linux Food
	//	-749918079:     true, // 123
	//	-1001373811109: true, // Linux Flood
	//	-1001558727831: true, // 123
	//	-1001740354030: true,
	//	-1001053617676: true,
	//}
)

func init() {
	AdminRightsCache = map[int64]map[int64]tgbotapi.ChatMember{}
}
