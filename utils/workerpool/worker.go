package workerpool

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"statBot/handlers"
	"statBot/utils"
	"strings"
)

var WorkerChanPool map[int]chan utils.ControlStruct

func init() {
	WorkerChanPool = make(map[int]chan utils.ControlStruct)
}

func UpdateWorker(id int, baseObj *utils.SharedBaseObject, updates tgbotapi.UpdatesChannel, controlChannel <-chan utils.ControlStruct) {
	for {
		select {
		case update := <-updates:
			log.Println(fmt.Sprintf("UpdateWorker %d started work with update id %d", id, update.UpdateID))
			if update.Message != nil {
				utils.UpdatesProcessed++

				utils.CallHandler(baseObj.Bot, update)
				if strings.ToLower(update.Message.Text) == "стало душно" {
					handlers.SendOpenedWindow(baseObj.Bot, update.Message)
				}
				if baseObj.AllowedChatsMode {
					if utils.AllowedChats[update.Message.Chat.ID] {
						fmt.Println("write to log ", update.Message.Chat.ID)
						utils.LogMessage(update, baseObj.DB)
					}
				} else {
					fmt.Println("write to log ", update.Message.Chat.ID)
					utils.LogMessage(update, baseObj.DB)
				}
			}
			if update.CallbackQuery != nil {
				utils.CallbackQueryHandler(baseObj.Bot, update)
			}
			log.Println(fmt.Sprintf("UpdateWorker %d finished work with update id %d", id, update.UpdateID))

		case controlStruct := <-controlChannel:
			log.Println("Parsing control struct")

			if controlStruct.Cmd == "Kill" {
				log.Println("Killing worker ", id)
				return
			}

		}
	}
}
