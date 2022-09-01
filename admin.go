package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"strings"
	"time"
)

func adminPrintStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	ChatID := -1001373811109
	logFile, err := os.ReadFile(fmt.Sprintf("%d.log", ChatID))
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"
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
		case "week":
			fromTime = time.Now().AddDate(0, 0, -7)
			fromTimeText = "последнюю неделю"
		case "day":
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"
		}
	}
	users := CalcUserMessages(logFile, fromTime)
	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Errorf(err.Error())
	}

}
