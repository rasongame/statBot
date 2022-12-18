package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"strings"
	"time"
)

func GenerateDeleteKeyboard(chatId int64, userId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("deleteStats;%d;%d", userId, chatId))))
}
func printStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startTime := time.Now()
	ChatID := message.Chat.ID
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"

	if cmdArgs != "" {
		args := strings.Split(cmdArgs, " ")
		switch args[0] {
		case "month":
			fromTime = time.Now().AddDate(0, -1, 0)
			fromTimeText = "последний месяц"

		case "week":
			fromTime = time.Now().AddDate(0, 0, -7)
			fromTimeText = "последнюю неделю"

		case "day":
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"

		}
	}

	totalMessages, users := CalcUserMessages(fromTime, ChatID)

	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	msg.Caption = fmt.Sprintf("Написано сообщений за %s \nВсего сообщений: %d\n%v", fromTimeText, totalMessages, time.Now().Sub(startTime))
	msg.ReplyToMessageID = message.MessageID

	msg.ReplyMarkup = GenerateDeleteKeyboard(message.Chat.ID, message.From.ID)
	bot.Send(msg)

}
func printPopularWords(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	ChatID := message.Chat.ID
	logFile, err := os.ReadFile(fmt.Sprintf("%d.log", ChatID))
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"

	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))

		return
	}
	cmdArgs := message.CommandArguments()
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

	wordsFreq := CalcPopularWords(logFile, fromTime)
	smallestNumber := int(math.Min(10, float64(len(wordsFreq))))
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%d самых популярных слов за %s\n", smallestNumber, fromTimeText))

	for i, v := range wordsFreq[:smallestNumber] {
		msg.Text = msg.Text + fmt.Sprintf("%d| %s: %d\n", i, v.Word, v.Freq)
	}
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
