package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	Handlers     map[string]Handler
	ReportChat   int64 = 559723688
	Superuser    int64 = 559723688
	AllowedChats       = map[int64]bool{
		559723688:      true, // rasongame
		-1001549183364: true, // Linux Food
		-749918079:     true, // 123
		-1001373811109: true,
	}
	helpText = strings.TrimSpace(`
/whoami - отправляет id юзера
/pop - отправляет самые использумые слова (только в чатах)
/stat - отправляет круговую диаграмму по теме кто больше нафлудил (только в чатах)
`)
)

// 1 Day = 86400 sec
func init() {
	Handlers = make(map[string]Handler)
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
	users := CalcUserMessages(logFile, fromTime)
	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
	photo := tgbotapi.FilePath(fileName)
	secondMsg := tgbotapi.NewPhoto(message.Chat.ID, photo)
	_, err = bot.Send(secondMsg)
	if err != nil {
		fmt.Errorf(err.Error())
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

func main() {
	token, err := os.ReadFile("token")
	if err != nil {
		panic(err.Error())
	}
	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	for i, value := range AllowedChats {
		log.Printf("AllowedChat %d: %t", i, value)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	AddHandler("astat", adminPrintStatToChat, SuperuserFilter)
	AddHandler("stat", printStatToChat, ChatOnly)
	AddHandler("test", testCmd, FalseFilter)
	AddHandler("whoami", idCmd, TrueFilter)
	AddHandler("pop", printPopularWords, ChatOnly)
	AddHandler("help", helpCmd, TrueFilter)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if handle, ok := Handlers[update.Message.Command()]; ok {
					if ok && handle.Filter(bot, update.Message) {
						handle.Handler(bot, update.Message)
					}
				}
			}
			if AllowedChats[update.Message.Chat.ID] {
				fmt.Println("write to log ", update.Message.Chat.ID)
				WriteToLog(bot, update.Message)
			}
		}
	}
}
