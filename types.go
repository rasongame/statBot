package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SomePlaceholder struct {
	User     *tgbotapi.User
	Messages int
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
