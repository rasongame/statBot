package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"statBot/utils"
	"time"
)

func init() {

}

func CalcUserMessages(from time.Time, to time.Time, chatId int64) (int, []utils.SomePlaceholder) {
	users := make(map[int64]utils.SomePlaceholder)
	var messageListForChats []utils.ChatMessage
	var totalMessages int
	queryString := "chat_id = ? and ? >= date and date >= ?"
	utils.DB.Where(queryString, chatId, to.Unix(), from.Unix()).Find(&messageListForChats)
	for chatMessageKey := range messageListForChats {
		chatMessage := messageListForChats[chatMessageKey]
		totalMessages++
		uzer := users[chatMessage.UserId]
		if uzer.User == nil {
			uzer.User = &tgbotapi.User{
				ID:        chatMessage.UserId,
				IsBot:     false,
				FirstName: chatMessage.UserFirstName,
				LastName:  chatMessage.UserLastName,
				UserName:  chatMessage.UserUsername,
			}

		}
		if chatMessage.Date != 0 {
			tm := time.Unix(int64(chatMessage.Date), 0)
			hour := tm.Hour()
			uzer.MessagesAt[hour]++
		}
		uzer.Messages++
		users[chatMessage.UserId] = uzer

	}
	s := make([]utils.SomePlaceholder, 0, len(users))
	for _, v := range users {
		s = append(s, v)
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Messages > s[j].Messages
	})

	return totalMessages, s
}
func CalcPopularWords(log []byte, fromTime time.Time) []utils.WordFreq {
	var wordFreqs []utils.WordFreq
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].Freq > wordFreqs[j].Freq
	})
	return wordFreqs
}
