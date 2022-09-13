package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

type CacheUser struct {
	User     User
	LifeTime int64
}

const (
	LinFloodID int64 = -1001373811109
	ReportChat int64 = 559723688
	Superuser  int64 = 559723688
)

var (
	DB           *gorm.DB
	Handlers     map[string]Handler
	CachedUsers  map[int64]CacheUser
	AllowedChats = map[int64]bool{
		559723688:      true, // rasongame
		-1001549183364: true, // Linux Food
		-749918079:     true, // 123
		-1001373811109: true,
	}
	helpText = strings.TrimSpace(`
/whoami - отправляет id юзера
/pop - отправляет самые использумые слова (только в чатах)
/stat - отправляет круговую диаграмму по теме кто больше нафлудил (только в чатах)
Важное замечание: 24 часа это значит 24 часа. Значения в статистике это и значат
`)
)

// 1 Day = 86400 sec
func init() {
	Handlers = make(map[string]Handler)
	CachedUsers = make(map[int64]CacheUser)
}

func InitBot() *tgbotapi.BotAPI {
	token, ok := os.LookupEnv("rtoken")
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: env variable \"rtoken\" is not set")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}
	statsCmd := tgbotapi.BotCommand{
		Command:     "stats",
		Description: "Стата по количеству сообщений от юзеров (/stats day/week/month)",
	}
	popCmd := tgbotapi.BotCommand{
		Command:     "pop",
		Description: "Стата по популярным словам (day/week/month)",
	}

	cmds := tgbotapi.NewSetMyCommands(statsCmd, popCmd)
	bot.Send(cmds)
	AddHandler("astats", adminPrintStatToChat, IsAdminFilter)
	AddHandler("stats", printStatToChat, ChatOnly)
	AddHandler("test", testCmd, FalseFilter)
	AddHandler("whoami", idCmd, TrueFilter)
	AddHandler("pop", printPopularWords, ChatOnly)
	AddHandler("health", adminSendBotHealth, IsAdminFilter)
	AddHandler("help", func(api *tgbotapi.BotAPI, message *tgbotapi.Message) {
		msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
		bot.Send(msg)

	}, TrueFilter)

	return bot
}

func main() {
	bot := InitBot()
	var err error
	DB, err = gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	DB.AutoMigrate(&Chat{}, &User{})
	LoadCache(DB)
	for i, value := range AllowedChats {
		log.Printf("AllowedChat %d: %t", i, value)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if handle, ok := Handlers[update.Message.Command()]; ok {
					if ok && handle.Filter(bot, update.Message) {
						go handle.Handler(bot, update.Message)
					}
				}
			}
			if strings.ToLower(update.Message.Text) == "стало душно" {
				sendOpenedWindow(bot, update.Message)
			}
			if AllowedChats[update.Message.Chat.ID] {
				fmt.Println("write to log ", update.Message.Chat.ID)
				go ProcessDB(update)
				WriteToLog(bot, update.Message)
			}
		}
	}

}
func ProcessDB(update tgbotapi.Update) {
	var ch *Chat
	var user *User
	DB.Where(&Chat{Id: update.Message.Chat.ID}).Find(&ch)
	DB.Where(&User{Id: update.Message.From.ID}).Find(&user)
	if ch.Id == 0 {
		fmt.Println("Add to DB Chat ID ", update.Message.Chat.ID)
		DB.Create(&Chat{
			Id:    update.Message.Chat.ID,
			Type:  update.Message.Chat.Type,
			Title: update.Message.Chat.Title,
		})
	}
	if user.Id == 0 {
		fmt.Println("Add to DB User ID ", update.Message.From.ID)
		DB.Create(&User{
			Id:           update.Message.From.ID,
			FirstName:    update.Message.From.FirstName,
			LastName:     update.Message.From.LastName,
			Username:     update.Message.From.UserName,
			LanguageCode: update.Message.From.LanguageCode,
		})
	}
}
