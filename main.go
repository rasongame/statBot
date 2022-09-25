package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"statBot/handlers"
	"statBot/utils"
	"strings"
	"time"
)

var (
	BotStarted = time.Now()
	DB         *gorm.DB

	AllowedChats = map[int64]bool{
		559723688:      true, // rasongame
		-1001549183364: true, // Linux Food
		-749918079:     true, // 123
		-1001373811109: true,
		-1001558727831: true, // 123
	}
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

func strContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isExplicitAllowedCommand(value string) bool {
	value = strings.ToLower(value[1:])
	return strContains([]string{"stats", "astats", "health", "pop"}, value)
}

func rightCommandExtractor(m *tgbotapi.Message, botNickName string) string {
	if !m.IsCommand() {
		return ""
	}

	// Explicit bot call in public chats
	if m.Chat.Type == "supergroup" {
		splittedCommand := strings.Split(m.Text, "@")
		if len(splittedCommand) > 1 && strings.ToLower(splittedCommand[1]) == strings.ToLower(botNickName) {
			return splittedCommand[0][1:]
		}
	} else {
		return m.Command()
	}

	if isExplicitAllowedCommand(m.Text) {
		return m.Command()
	}

	return ""
}

func main() {
	bot := InitBot()
	var err error
	DB, err = gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
	utils.PanicErr(err)
	err = DB.AutoMigrate(&utils.Chat{}, &utils.User{})
	utils.PanicErr(err)
	utils.LoadCache(DB, utils.CachedUsers)
	for i, value := range AllowedChats {
		log.Printf("AllowedChat %d: %t", i, value)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			CallHandler(bot, update)
			if strings.ToLower(update.Message.Text) == "стало душно" {
				handlers.SendOpenedWindow(bot, update.Message)
			}
			if AllowedChats[update.Message.Chat.ID] {
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
						userCacheFinal.Messages++

					}

				}()
				WriteToLog(bot, update.Message)

			}
		}
	}
}

func CallHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		handler_name := rightCommandExtractor(update.Message, bot.Self.UserName)

		if handle, ok := utils.Handlers[handler_name]; ok {
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
