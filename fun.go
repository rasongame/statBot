package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	windowStickerFileId = "CAACAgIAAxkBAAER_M5jHd5cpDyOcIbq3QmpHR5mnmySBgACjwADfI5YFfhv_xslSwqzKQQ"
)

func sendOpenedWindow(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	sticker := tgbotapi.FileID(windowStickerFileId)
	msg := tgbotapi.NewSticker(message.Chat.ID, sticker)
	bot.Send(msg)
}
