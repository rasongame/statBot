package main

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
	word string
	freq int
}

func (p WordFreq) String() string {
	return fmt.Sprintf("%s %d", p.word, p.freq)
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
