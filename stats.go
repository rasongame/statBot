package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	blacklistWords  []string
	whitelistRegExp = regexp.MustCompile("[a-zA-Zа-яА-Я0-9']+")
)

func init() {
	fmt.Println("func init() stats.go start")
	words, err := os.ReadFile("blacklistWords.txt")
	if err != nil {
		panic(err.Error())
	}
	blacklistWords = strings.Split(string(words), "\n")
	fmt.Println("func init() stats.go end")
}
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Key ID; Value = SomePlaceholder
func CalcUserMessages(log []byte, from time.Time) []SomePlaceholder { // map[int64]SomePlaceholder
	splitted := strings.Split(strings.ReplaceAll(string(log), "\r\n", "\n"), "\n")
	users := make(map[int64]SomePlaceholder)
	for k, str := range splitted {
		var unm tgbotapi.Message
		//str := strings.TrimSuffix(str, "\n")
		//fmt.Println(i, str)
		err := json.Unmarshal([]byte(str), &unm)
		if err != nil {
			fmt.Println(k, err)
			continue
		}
		if unm.SenderChat != nil {
			continue
		}
		if len(unm.Text) < 3 {
			continue
		}
		yesterdayTime := from.Unix()
		if int64(unm.Date) >= yesterdayTime {
			//fmt.Println(users[int64(unm.From.ID)])
			uzer := users[int64(unm.From.ID)]
			if uzer.User == nil {
				uzer.User = unm.From
			}
			uzer.Messages++
			users[int64(unm.From.ID)] = uzer

		}
	}
	s := make([]SomePlaceholder, 0, len(users))
	// append all map keys-value pairs to the slice
	for _, v := range users {
		s = append(s, v)
	}
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Messages > s[j].Messages
	})

	return s
}
func CalcPopularWords(log []byte, fromTime time.Time) []WordFreq {
	splitted := strings.Split(strings.ReplaceAll(string(log), "\r\n", "\n"), "\n")
	var text string

	{
		for _, str := range splitted {
			var m tgbotapi.Message
			err := json.Unmarshal([]byte(str), &m)
			if err != nil {
				continue
			}
			yesterdayTime := fromTime.Unix()
			if int64(m.Date) >= yesterdayTime {
				text = text + " " + m.Text
			}
		}
	}
	text = strings.ToLower(text)
	for _, word := range blacklistWords {
		text = strings.ReplaceAll(text, " "+word+" ", "")

	}
	matches := whitelistRegExp.FindAllString(text, -1)
	words := make(map[string]int)
	for _, match := range matches {
		words[match]++
	}
	var wordFreqs []WordFreq
	for k, v := range words {
		wordFreqs = append(wordFreqs, WordFreq{k, v})

	}
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].freq > wordFreqs[j].freq
	})
	return wordFreqs
}