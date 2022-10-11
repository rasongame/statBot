package main

import (
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
	DB *gorm.DB

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
	utils.CachedUsers = make(map[int64]utils.CacheUser)
}

func main() {
	bot := InitBot()
	var err error
	DB, err = gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
	utils.PanicErr(err)
	err = DB.AutoMigrate(&utils.Chat{}, &utils.User{})
	utils.PanicErr(err)
	utils.LoadCache(DB, utils.CachedUsers)
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
				go ProcessDB(update)
				go func() {
					if utils.ChatLogIsLoaded[update.Message.Chat.ID] {
						userCache := utils.ChatLogMessageCache[update.Message.Chat.ID]
						userCacheFinal := userCache[update.Message.From.ID]
						if userCacheFinal == nil {
							userCacheFinal = &utils.SomePlaceholder{
								User:       update.Message.From,
								Messages:   0,
								LastSeenAt: time.Now(),
							}
						}
						if utils.ChatLogIsLoadedTime[update.Message.Chat.ID].Sub(time.Now()) >= time.Hour*24 {
							userCacheFinal.Messages = 0
							utils.ChatLogIsLoadedTime[update.Message.Chat.ID] = time.Now()
							fmt.Println("day ruined... updating utils.ChatLogIsLoadedTime[update.Message.Chat.ID]")
						}

						userCacheFinal.Messages++

					}

				}()
				WriteToLog(bot, update.Message)

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

func ProcessDB(update tgbotapi.Update) {
	var ch *utils.Chat
	var user *utils.User
	DB.Where(&utils.Chat{Id: update.Message.Chat.ID}).Find(&ch)
	DB.Where(&utils.User{Id: update.Message.From.ID}).Find(&user)
	if ch.Id == 0 {
		fmt.Println("Add to DB Chat ID ", update.Message.Chat.ID)
		DB.Create(&utils.Chat{
			Id:    update.Message.Chat.ID,
			Type:  update.Message.Chat.Type,
			Title: update.Message.Chat.Title,
		})
	}
	if user.Id == 0 {
		fmt.Println("Add to DB User ID ", update.Message.From.ID)
		DB.Create(&utils.User{
			Id:           update.Message.From.ID,
			FirstName:    update.Message.From.FirstName,
			LastName:     update.Message.From.LastName,
			Username:     update.Message.From.UserName,
			LanguageCode: update.Message.From.LanguageCode,
		})
	}
}
