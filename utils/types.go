package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"runtime"
	"time"
)

type SomePlaceholder struct {
	User       *tgbotapi.User
	Messages   int
	LastSeenAt time.Time
}
type WordFreq struct {
	Word string
	Freq int
}

func (p WordFreq) String() string {
	return fmt.Sprintf("%s %d", p.Word, p.Freq)
}

type HandlerFunc func(api *tgbotapi.BotAPI, message *tgbotapi.Message)
type FilterFunc func(api *tgbotapi.BotAPI, message *tgbotapi.Message) bool
type Handler struct {
	Handler HandlerFunc
	Filter  FilterFunc
}
type User struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	LastSeen     time.Time
}
type Chat struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}
type ChatMessage struct {
	ChatId        int64  `json:"chat_id"`
	MessageId     int64  `json:"message_id"`
	UserId        int64  `json:"user_id"`
	Text          string `json:"text"`
	Date          int    `json:"date"`
	UserFirstName string `json:"user_first_name"`
	UserLastName  string `json:"user_last_name"`
	UserUsername  string `json:"user_username"`
}

func GetAboutInfo() AboutBot {
	info := AboutBot{}
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err.Error())
		info.Hostname = "PotatoPC"
	} else {
		info.Hostname = hostname
	}
	info.GoVersion = runtime.Version()
	info.Platform = runtime.GOOS
	info.Architecture = runtime.GOARCH
	return info
}

type AboutBot struct {
	Hostname     string
	GoVersion    string
	Platform     string
	Architecture string
}
type CacheUser struct {
	User     User
	LifeTime int64
}
