package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func testCmd(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	_, err := b.Send(tgbotapi.NewMessage(m.Chat.ID, "hello world"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
func idCmd(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	_, err := b.Send(tgbotapi.NewMessage(m.Chat.ID, strconv.FormatInt(m.From.ID, 10)))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func helpCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, helpText))
}

func printStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	ChatID := message.Chat.ID
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
	users := CalcUserMessages(logFile, fromTime)
	go func() {
		for _, v := range users {
			UpdateCache(&v, DB)
		}
	}()
	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	msg.Caption = fmt.Sprintf("Написано сообщений за %s", fromTimeText)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Println(err.Error())
	}

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
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("10 самых популярных слов за %s\n", fromTimeText))
	smallestNumber := int(math.Min(10, float64(len(wordsFreq))))

	for i, v := range wordsFreq[:smallestNumber] {
		msg.Text = msg.Text + fmt.Sprintf("%d| %s: %d\n", i, v.word, v.freq)
	}
	bot.Send(msg)
}
