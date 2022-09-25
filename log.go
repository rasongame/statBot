package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"statBot/utils"
	"strings"
)

func PrintLogToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	fmt.Println(len(args))
	if len(args) < 2 {
		entries, err := os.ReadDir(".")
		if err != nil {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "")
		for _, entry := range entries {
			msg.Text = msg.Text + entry.Name() + "\n"
		}
		bot.Send(msg)
		return
	}
	file := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FilePath(args[1]))
	bot.Send(file)

}
func WriteToLog(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	f, err := os.OpenFile(fmt.Sprintf("%d.log", message.Chat.ID), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		msg := tgbotapi.NewMessage(utils.ReportChat, err.Error())
		_, botErr := bot.Send(msg)
		if botErr != nil {
			fmt.Errorf("%s", botErr.Error())
		}
		return
	}
	encoded, _ := json.Marshal(message)

	_, err = f.WriteString(fmt.Sprintf("%s\n", encoded))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
	}
	defer f.Close()
}
