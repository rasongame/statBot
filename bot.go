package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"statBot/filters"
	"statBot/handlers"
)

func InitBotCommands(bot *tgbotapi.BotAPI) {
	statsCmd := tgbotapi.BotCommand{
		Command:     "stats",
		Description: "Стата по количеству сообщений от юзеров (/stats day/week/month)",
	}
	popCmd := tgbotapi.BotCommand{
		Command:     "pop",
		Description: "Стата по популярным словам (day/week/month)",
	}
	decodeCmd := tgbotapi.BotCommand{
		Command:     "decode",
		Description: "\"ghbdtn vbh\" => \" привет мир\" ",
	}
	decode64Cmd := tgbotapi.BotCommand{
		Command:     "decodebase64",
		Description: "Base64 => \" привет мир\" ",
	}

	cmds := tgbotapi.NewSetMyCommands(statsCmd, popCmd, decodeCmd, decode64Cmd)
	bot.Send(cmds)
}
func InitBotHandlers(bot *tgbotapi.BotAPI) {
	//
	AddHandler("decode", handlers.SendDecodedMessage, filters.TrueFilter)
	AddHandler("decodebase64", handlers.SendDecodedBase64Message, filters.TrueFilter)
	AddHandler("whoami", handlers.Id, filters.TrueFilter)
	//AddHandler("help", nil, filters.TrueFilter)
	//
	AddHandler("health", handlers.SendBotHealth, filters.IsAdminFilter)
	AddHandler("astats", adminPrintStatToChat, filters.IsAdminFilter)
	AddHandler("adminlist", handlers.SendAdminList, filters.IsAdminFilter)
	//
	AddHandler("stats", printStatToChat, filters.ChatOnly)
	AddHandler("pop", printPopularWords, filters.ChatOnly)
	//
	AddHandler("test", handlers.Test, filters.FalseFilter)

}
func InitBot() *tgbotapi.BotAPI {
	token, ok := os.LookupEnv("rtoken")
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: env variable \"rtoken\" is not set")
		os.Exit(1)
	}
	debugMode := false
	debugModeEnv, ok := os.LookupEnv("rdebug")
	if ok && debugModeEnv == "1" {
		debugMode = true
	}

	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = debugMode
	if err != nil {
		log.Panic(err)
	}

	InitBotHandlers(bot)
	InitBotCommands(bot)

	return bot
}
