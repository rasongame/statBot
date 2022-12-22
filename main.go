package main

import (
	"flag"
	"fmt"
	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"statBot/handlers"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

var (
	dbPath           string
	BlacklistedUsers = map[int64]bool{5449020876: true}
	helpText         = strings.TrimSpace(`
/whoami - отправляет id юзера
/pop - отправляет самые использумые слова (только в чатах)
/stat - отправляет круговую диаграмму по теме кто больше нафлудил (только в чатах)
Важное замечание: 24 часа это значит 24 часа. Значения в статистике это и значат
`)
)

// 1 Day = 86400 sec
func init() {
	utils.Handlers = make(map[string]utils.Handler)
}
func parseFlags() {
	var allowedChats, blacklistedUsers string
	defaultAllowedChats := "559723688,-1001549183364,-749918079,-1001373811109,-1001558727831,-1001740354030,-1001053617676"
	defaultBlackListed := "5449020876"

	flag.StringVar(&dbPath, "db", "bot.db", "bot.db path")
	flag.StringVar(&allowedChats, "allowedChats", defaultAllowedChats, "allowedChats")
	flag.StringVar(&blacklistedUsers, "blacklistedUsers", defaultBlackListed, "blacklistedUsers")

	flag.Parse()

	splittedChats := strings.Split(allowedChats, ",")
	splittedUsers := strings.Split(blacklistedUsers, ",")

	for i, str := range splittedChats {
		fmt.Println(i, "Adding", str)
		chatId, err := strconv.ParseInt(str, 10, 64)
		utils.PanicErr(err)
		utils.AllowedChats[chatId] = true
	}
	for i, str := range splittedUsers {
		fmt.Println(i, "Blacklisting", str)
		userId, err := strconv.ParseInt(str, 10, 64)
		utils.PanicErr(err)
		BlacklistedUsers[userId] = true
	}
}

func main() {
	parseFlags()
	bot := InitBot()
	var err error
	utils.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	DB := utils.DB
	x, err := DB.DB()
	x.SetMaxOpenConns(1)
	utils.PanicErr(err)
	err = DB.AutoMigrate(&utils.Chat{}, &utils.User{}, &utils.ChatMessage{}, &utils.ChatAudio{})
	utils.PanicErr(err)
	for i, value := range utils.AllowedChats {
		log.Printf("AllowedChat %d: %t", i, value)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			utils.UpdatesProcessed++

			CallHandler(bot, update)
			if strings.ToLower(update.Message.Text) == "стало душно" {
				handlers.SendOpenedWindow(bot, update.Message)
			}
			if utils.AllowedChats[update.Message.Chat.ID] {
				fmt.Println("write to log ", update.Message.Chat.ID)

				go func() {
					if update.EditedMessage != nil {
						msg := update.EditedMessage
						if msg.Audio != nil {
							tx := DB.Create(&utils.ChatAudio{
								MessageId:     int64(msg.MessageID),
								UniqueFileId:  msg.Audio.FileID,
								ChatId:        msg.Chat.ID,
								FromId:        msg.From.ID,
								UserFirstName: msg.From.FirstName,
								UserLastName:  msg.From.LastName,
								UserUsername:  msg.From.UserName,
							})
							tx.Commit()
						}
					}
					if update.Message != nil {
						tx := DB.Begin()
						msg := update.Message
						tx.Create(&utils.ChatMessage{
							ChatId:        msg.Chat.ID,
							MessageId:     int64(update.Message.MessageID),
							Text:          msg.Text,
							UserId:        msg.From.ID,
							Date:          msg.Date,
							UserFirstName: msg.From.FirstName,
							UserLastName:  msg.From.LastName,
							UserUsername:  msg.From.UserName,
						})
						if update.Message.Audio != nil {
							tx.Create(&utils.ChatAudio{
								UniqueFileId:  msg.Audio.FileID,
								MessageId:     int64(msg.MessageID),
								ChatId:        msg.Chat.ID,
								FromId:        msg.From.ID,
								UserFirstName: msg.From.FirstName,
								UserLastName:  msg.From.LastName,
								UserUsername:  msg.From.UserName,
							})
						}
						tx.Commit()
					}

				}()
				//WriteToLog(bot, update.Message)

			}
		}
		if update.CallbackQuery != nil {
			go CallbackQueryHandler(bot, update)
		}
	}
}
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
		utils.PanicErr(err)
		chatId, err = strconv.ParseInt(separated[2], 10, 64)
		utils.PanicErr(err)

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
		utils.PanicErr(err)
	}
}
func CallHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		handlerName := utils.RightCommandExtractor(update.Message, bot.Self.UserName)

		if handle, ok := utils.Handlers[handlerName]; ok {
			if ok && handle.Filter(bot, update.Message) {
				go func() {
					timeStart := time.Now()
					handle.Handler(bot, update.Message)
					fmt.Println(time.Now().Sub(timeStart))
				}()
			}
		}
	}
}
