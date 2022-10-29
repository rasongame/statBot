package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"os"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

func adminPrintStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var chatID = utils.LinFloodID
	var convErr error
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"
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

		default:
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"

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

	logFile, err := os.ReadFile(fmt.Sprintf("%d.log", chatID))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return
	}
	var users []utils.SomePlaceholder
	// CalcUserMessagesLegacy нужен чтобы избежать кэширования ненужной бяки
	users = CalcUserMessagesLegacy(logFile, fromTime, message.Chat.ID)

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
