package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime"
	"statBot/utils"
	"time"
)

var (
	BotStarted = time.Now()
)

func SendBotHealth(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var mem runtime.MemStats

	text :=
		`
BotUptime: %s
Hostname: %s
Go Version: %s
Platform: %s,
Architecture: %s, 
Alloc: %v MiB
Total Alloc: %v MiB
Sys: %v MiB
GC Calls: %v
NumCPU: %d
"
`
	uptime := time.Now().Sub(BotStarted)
	info := utils.GetAboutInfo()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	runtime.ReadMemStats(&mem)
	msg.Text = fmt.Sprintf(msg.Text, uptime, info.Hostname,
		info.GoVersion, info.Platform, info.Architecture,
		utils.BToMb(mem.Alloc), utils.BToMb(mem.TotalAlloc), utils.BToMb(mem.Sys), mem.NumGC, runtime.NumCPU())
	runtime.GC()
	bot.Send(msg)
}
