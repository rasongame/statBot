package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

func adminPrintStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startTime := time.Now()
	var chatID = utils.LinFloodID
	var convErr error
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"
	var dayIsSelected bool

	var err error
	if cmdArgs != "" {
		args := strings.Split(cmdArgs, " ")

		switch args[0] {
		case "year":
			fromTime = time.Now().AddDate(-1, 0, 0)
			fromTimeText = "последний год"
		case "alltime":
			fromTime = time.Unix(0, 0)
			fromTimeText = "всё время существования бота здесь"
		case "month":
			fromTime = time.Now().AddDate(0, 0, -30)
			fromTimeText = "последний месяц"

		case "week":
			fromTime = time.Now().AddDate(0, 0, -7)
			fromTimeText = "последнюю неделю"

		case "day":
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"

		case "cal":
			dayIsSelected = true
			pattern := "02.01.2006"
			fromTime, err = time.Parse(pattern, args[1])
			fromTimeText = args[1]
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
				bot.Send(msg)
				return
			}

		default:
			dayIsSelected = true
			pattern := "02.01.2006"
			fromTime, err = time.Parse(pattern, args[0])
			fromTimeText = args[0]
			if err != nil {
				dayIsSelected = false
				fromTime = time.Now().AddDate(0, 0, -1)
				fromTimeText = "последние 24 часа"
			}

		}
		if len(args) <= 1 {
			chatID = utils.LinFloodID
		} else {
			if args[1] != "" {
				chatID, convErr = strconv.ParseInt(args[1], 10, 64)
				if convErr != nil {
					log.Println("invalid chatId. using default")
					if val, ok := utils.Aliases[args[1]]; ok {
						chatID = val
					} else {
						chatID = utils.LinFloodID
					}

				}
			}
		}
	}
	to := time.Now()
	if dayIsSelected {
		to = fromTime.AddDate(0, 0, 1)
	}
	totalMessages, users := CalcUserMessages(fromTime, to, chatID)
	if totalMessages <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No messages? o_O")
		bot.Send(msg)
		return
	}
	fileName := fmt.Sprintf("%d-activeStat.png", chatID)
	//RenderActiveUsers(users, fileName, int(math.Min(15, float64(len(users)))), fromTimeText)

	RenderActiveUsers(users, fileName, int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	msg.Caption = fmt.Sprintf("Написано сообщений: %d\nОбработано за %v", totalMessages, time.Now().Sub(startTime))
	bot.Send(msg)

}
