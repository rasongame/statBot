package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

var (
	BotStarted = time.Now()
)

func SendAdminList(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var args []string
	msg := "Admin list\n"
	fmt.Println(message.CommandArguments())
	var chatID = message.Chat.ID
	if message.CommandArguments() != "" {
		fmt.Println(len(strings.Split(message.CommandArguments(), " ")))
		args = strings.Split(message.CommandArguments(), " ")

	}
	if len(args) >= 1 {
		customChatID, err := strconv.ParseInt(args[0], 10, 64)
		if err == nil {
			chatID = customChatID
		} else {
			utils.PanicErr(err)
		}
	}

	msg = fmt.Sprintln("ChatID:", chatID, "\n", msg)
	for _, value := range utils.AdminRightsCache[chatID] {

		isCreator := ""
		if value.IsCreator() {
			isCreator = "+"
		} else {
			isCreator = "-"
		}

		msg = fmt.Sprintln(msg, value.User.ID, value.User.FirstName, value.User.LastName,
			value.User.UserName, value.CanDeleteMessages, isCreator)
	}
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, msg))
	utils.PanicErr(err)
}
func SendBotHealth(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var mem runtime.MemStats

	text :=
		`
Updates Processed: %d
BotUptime: %s 
Hostname: %s
Go Version: %s
Platform: %s,
Architecture: %s, 
Alloc: %v MiB
Total Alloc: %v MiB
Heap Alloc: %v MiB
Heap Use: %v MiB
Sys: %v MiB
GC Calls: %v
NumCPU: %d

`
	uptime := time.Now().Sub(BotStarted)
	info := utils.GetAboutInfo()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	runtime.ReadMemStats(&mem)
	msg.Text = fmt.Sprintf(msg.Text, utils.UpdatesProcessed, uptime, info.Hostname,
		info.GoVersion, info.Platform, info.Architecture,
		utils.BToMb(mem.Alloc),
		utils.BToMb(mem.TotalAlloc),
		utils.BToMb(mem.HeapAlloc),
		utils.BToMb(mem.HeapInuse),
		utils.BToMb(mem.Sys),
		mem.NumGC,
		runtime.NumCPU(),
	)
	runtime.GC()
	msg.ReplyToMessageID = message.MessageID
	_, err := bot.Send(msg)

	if err != nil {
		fmt.Println(err)
	}
}
