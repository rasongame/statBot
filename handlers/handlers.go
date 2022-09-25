package handlers

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func Test(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	_, err := b.Send(tgbotapi.NewMessage(m.Chat.ID, "hello world"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
func Id(b *tgbotapi.BotAPI, m *tgbotapi.Message) {
	_, err := b.Send(tgbotapi.NewMessage(m.Chat.ID, strconv.FormatInt(m.From.ID, 10)))
	if err != nil {
		fmt.Println(err.Error())
	}
}

const (
	windowStickerFileId = "CAACAgIAAxkBAAER_M5jHd5cpDyOcIbq3QmpHR5mnmySBgACjwADfI5YFfhv_xslSwqzKQQ"
)

var (
	alphabet = make(map[string]string, 96)
)

func SendOpenedWindow(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	sticker := tgbotapi.FileID(windowStickerFileId)
	msg := tgbotapi.NewSticker(message.Chat.ID, sticker)
	bot.Send(msg)
}
func SendDecodedMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		return
	}
	decodedText := message.ReplyToMessage.Text
	for sourceLetter, letter := range alphabet {
		decodedText = strings.ReplaceAll(decodedText, sourceLetter, letter)
	}
	fmt.Println(decodedText)
	msg := tgbotapi.NewMessage(message.Chat.ID, decodedText)
	bot.Send(msg)
}
func SendDecodedBase64Message(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(message.ReplyToMessage.Text)
	if err != nil {
		decoded = []byte(err.Error())
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, string(decoded))
	bot.Send(msg)
}
func init() {
	go func() {
		alphabet["Q"] = "Й"
		alphabet["W"] = "Ц"
		alphabet["E"] = "У"
		alphabet["R"] = "К"
		alphabet["T"] = "Е"
		alphabet["Y"] = "Н"
		alphabet["U"] = "Г"
		alphabet["I"] = "Ш"
		alphabet["O"] = "Щ"
		alphabet["P"] = "З"
		alphabet["{"] = "Х"
		alphabet["}"] = "Ъ"
		alphabet["A"] = "ф"
		alphabet["S"] = "Ы"
		alphabet["D"] = "В"
		alphabet["F"] = "А"
		alphabet["G"] = "П"
		alphabet["H"] = "Р"
		alphabet["J"] = "О"
		alphabet["K"] = "Л"
		alphabet["L"] = "Д"
		alphabet[":"] = "Ж"
		alphabet["|"] = "Э"
		alphabet["Z"] = "Я"
		alphabet["X"] = "Ч"
		alphabet["C"] = "С"
		alphabet["V"] = "М"
		alphabet["B"] = "И"
		alphabet["N"] = "Т"
		alphabet["M"] = "Ь"
		alphabet["<"] = "Б"
		alphabet[">"] = "Ю"
		//-----------------
		//-----------------
		alphabet["q"] = "й"
		alphabet["w"] = "ц"
		alphabet["e"] = "у"
		alphabet["r"] = "к"
		alphabet["t"] = "е"
		alphabet["y"] = "н"
		alphabet["u"] = "г"
		alphabet["i"] = "ш"
		alphabet["o"] = "щ"
		alphabet["p"] = "з"
		alphabet["["] = "х"
		alphabet["]"] = "ъ"
		// ---
		alphabet["a"] = "ф"
		alphabet["s"] = "ы"
		alphabet["d"] = "в"
		alphabet["f"] = "а"
		alphabet["g"] = "п"
		alphabet["h"] = "р"
		alphabet["j"] = "о"
		alphabet["k"] = "л"
		alphabet["l"] = "д"
		alphabet[";"] = "ж"
		alphabet["\\"] = "э"
		// --
		alphabet["z"] = "я"
		alphabet["x"] = "ч"
		alphabet["c"] = "с"
		alphabet["v"] = "м"
		alphabet["b"] = "и"
		alphabet["n"] = "т"
		alphabet["m"] = "ь"
		alphabet[","] = "б"
		alphabet["."] = "ю"
		alphabet["/"] = "."
		alphabet["?"] = ","
		alphabet["&"] = "?"
		alphabet["@"] = "\""
	}() // Generate var alphabet

}
