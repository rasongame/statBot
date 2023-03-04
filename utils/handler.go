package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

func CallbackQueryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callbackData := update.CallbackData()
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	var eventType = ""
	separated := strings.Split(callbackData, ";")
	if len(separated) >= 1 {
		eventType = separated[0]
	} else {
		return
	}
	if eventType == "deleteStats" {
		var chatId int64 = 0
		var userId int64 = 0
		userId, err := strconv.ParseInt(separated[1], 10, 64)
		PanicErr(err)
		chatId, err = strconv.ParseInt(separated[2], 10, 64)
		PanicErr(err)

		if chatId == update.CallbackQuery.Message.Chat.ID && userId == update.CallbackQuery.From.ID {
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			}
			bot.Send(deleteConfig)
		} else {
			callback.Text = "ты кто такой, давай до свиданья, вася"
		}

	}

	if _, err := bot.Request(callback); err != nil {
		PanicErr(err)
	}
}
func CallHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		handlerName := RightCommandExtractor(update.Message, bot.Self.UserName)

		if handle, ok := Handlers[handlerName]; ok {
			if ok && handle.Filter(bot, update.Message) {
				timeStart := time.Now()
				handle.Handler(bot, update.Message)
				fmt.Println(time.Now().Sub(timeStart))

			}
		}
	}
}
