package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func adminSendBotHealth(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
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
	info := GetAboutInfo()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	runtime.ReadMemStats(&mem)
	msg.Text = fmt.Sprintf(msg.Text, uptime, info.Hostname,
		info.GoVersion, info.Platform, info.Architecture,
		bToMb(mem.Alloc), bToMb(mem.TotalAlloc), bToMb(mem.Sys), mem.NumGC, runtime.NumCPU())
	runtime.GC()
	bot.Send(msg)
}
func adminPrintStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	logFile, err := os.ReadFile(fmt.Sprintf("%d.log", LinFloodID))
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"
	caseNumber := 0
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	if cmdArgs != "" {
		args := strings.Split(cmdArgs, " ")

		switch args[0] {
		case "month":
			fromTime = time.Now().AddDate(0, 0, -30)
			fromTimeText = "последний месяц"
			caseNumber = 1
		case "week":
			fromTime = time.Now().AddDate(0, 0, -7)
			fromTimeText = "последнюю неделю"
			caseNumber = 2

		case "day":
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"
			caseNumber = 3
		}
	}
	var users []SomePlaceholder
	// CalcUserMessagesLegacy нужен чтобы избежать кэширования ненужной бяки
	if caseNumber >= 1 {
		users = CalcUserMessagesLegacy(logFile, fromTime)
	} else {
		users = CalcUserMessages(logFile, fromTime)
	}

	for _, v := range users {
		UpdateCache(&v, DB)
	}

	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	msg.Caption = fmt.Sprintf("Написано сообщений за %s", fromTimeText)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Errorf(err.Error())
	}

}
