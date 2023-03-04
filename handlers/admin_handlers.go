package handlers

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"runtime"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

var (
	BotStarted = time.Now()
)

//	func SendAdminMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//		log.Println("SendAdminMessage fired")
//		var chatIds []int64
//		targetMessage := message.CommandArguments()
//		utils.DB.Raw("select distinct on (chat_id) chat_id from chat_messages ;")
//		utils.DB.Select("distinct on (chat_id) chat_id").Find(&chatIds)
//		fmt.Println("%v \n %s", chatIds, targetMessage)
//
// }
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

<code>Updates Processed: %d
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
NumGoroutine: %d
-- SQL
InUse: %d
Idle: %d
WaitCount: %d
WaitDuration: %d
MaxIdleClosed: %d
MaxIdleTimeClosed: %d
MaxLifetimeClosed: %d
--
etc.
AVG Load: %s
System Uptime: %s 
</code>
`
	sqlDB, _ := utils.DB.DB()
	var stats sql.DBStats
	if sqlDB != nil {
		stats = sqlDB.Stats()
	}
	systemUptime := "not available"
	loadAvg := ""
	if runtime.GOOS != "windows" {
		_systemUptime, err := os.ReadFile("/proc/uptime")
		_loadAvg, avgErr := os.ReadFile("/proc/loadavg")
		systemUptime = string(_systemUptime)
		loadAvg = string(_loadAvg)
		if err != nil {
			log.Println("Error occured:", err)
			systemUptime = "not available"
		}
		if avgErr != nil {
			log.Println("Error occured:", avgErr)
			loadAvg = "not available"
		}
	}

	uptime := time.Now().Sub(BotStarted)
	info := utils.GetAboutInfo()
	msg := tgbotapi.NewMessage(message.Chat.ID, text)

	runtime.ReadMemStats(&mem)
	msg.Text = fmt.Sprintf(msg.Text,
		utils.UpdatesProcessed, uptime, info.Hostname,
		info.GoVersion, info.Platform, info.Architecture,
		utils.BToMb(mem.Alloc),
		utils.BToMb(mem.TotalAlloc),
		utils.BToMb(mem.HeapAlloc),
		utils.BToMb(mem.HeapInuse),
		utils.BToMb(mem.Sys),
		mem.NumGC,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
		stats.MaxIdleClosed,
		stats.MaxIdleTimeClosed,
		stats.MaxLifetimeClosed,
		loadAvg,
		systemUptime,
	)
	runtime.GC()
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "html"
	_, err := bot.Send(msg)

	if err != nil {
		fmt.Println(err)
	}
}
