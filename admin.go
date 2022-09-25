package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"statBot/utils"
	"strings"
	"time"
)

func adminPrintStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	logFile, err := os.ReadFile(fmt.Sprintf("%d.log", utils.LinFloodID))
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
	var users []utils.SomePlaceholder
	// CalcUserMessagesLegacy нужен чтобы избежать кэширования ненужной бяки
	if caseNumber >= 1 {
		users = CalcUserMessagesLegacy(logFile, fromTime)
	} else {
		users = CalcUserMessages(logFile, fromTime)
	}

	for _, v := range users {
		utils.UpdateCache(&v, DB, CachedUsers)
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
