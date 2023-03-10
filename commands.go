package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"os"
	"statBot/utils"
	"strconv"
	"strings"
	"time"
)

func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
func printTable(table [][]string) string {
	myStr := ""
	// get number of columns from the first table row
	columnLengths := make([]int, len(table[0]))
	for _, line := range table {
		for i, val := range line {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}

	var lineLength int
	for _, c := range columnLengths {
		lineLength += c + 3 // +3 for 3 additional characters before and after each field: "| %s "
	}
	lineLength += 1 // +1 for the last "|" in the line

	for i, line := range table {
		if i == 0 { // table header
			myStr = myStr + fmt.Sprintf("+%s+\n", strings.Repeat("-", lineLength-2)) // lineLength-2 because of "+" as first and last character
		}
		for j, val := range line {
			myStr = myStr + fmt.Sprintf("| %-*s ", columnLengths[j], val)
			if j == len(line)-1 {
				myStr = myStr + fmt.Sprintf("|\n")
			}
		}
		if i == 0 || i == len(table)-1 { // table header or last line
			myStr = myStr + fmt.Sprintf("+%s+\n", strings.Repeat("-", lineLength-2)) // lineLength-2 because of "+" as first and last character
		}
	}
	return myStr
}
func GenerateTextStats(elements []utils.SomePlaceholder, limit int, fromTimeText string, totalMessages int) string {
	base := fmt.Sprintf("Статистика за %s.\nВсего сообщений: %d\n", fromTimeText, totalMessages)
	var HoursActivity [4]int
	for i, el := range elements[:limit] {

		for i := 0; i <= 23; i++ {
			switch {
			case i <= 7 && i > 1:
				HoursActivity[0] += el.MessagesAt[i]
			case i <= 12 && i > 8:
				HoursActivity[1] += el.MessagesAt[i]
			case i <= 18 && i > 12:
				HoursActivity[2] += el.MessagesAt[i]
			case i <= 23 && i > 18:
				HoursActivity[3] += el.MessagesAt[i]
			}

		}
		base = base + fmt.Sprintf(
			"%d. %s %s: %d\n",
			1+i,
			el.User.FirstName,
			el.User.LastName,
			el.Messages,
		)
	}
	shit := [][]string{
		{"0-6", "6-12", "12-18", "18-24"},
		{strconv.Itoa(HoursActivity[3]), strconv.Itoa(HoursActivity[0]), strconv.Itoa(HoursActivity[1]), strconv.Itoa(HoursActivity[2])},
	}
	avgMessagesInHour := printTable(shit)
	//avgMessagesInHour := fmt.Sprintf("+%s+\n", strings.Repeat("-", 45))
	//avgMessagesInHour = avgMessagesInHour + fmt.Sprintf("| %-6s | %-6s | %-6s | %-6s |\n", "0-6", "6-12", "12-18", "18-24")
	//avgMessagesInHour = avgMessagesInHour + fmt.Sprintf("|%s|\n", strings.Repeat("-", 45))
	//avgMessagesInHour = avgMessagesInHour + fmt.Sprintf("| %-6d | %-6d | %-6d | %-6d |\n", HoursActivity[3], HoursActivity[0], HoursActivity[1], HoursActivity[2])
	//avgMessagesInHour = avgMessagesInHour + fmt.Sprintf("+%s+\n", strings.Repeat("-", 45))
	return fmt.Sprintf("%s\n%s", base, avgMessagesInHour)

}
func GenerateDeleteKeyboard(chatId int64, userId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("deleteStats;%d;%d", userId, chatId))))
}
func printStatToChat(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startTime := time.Now()
	ChatID := message.Chat.ID
	cmdArgs := message.CommandArguments()
	fromTime := time.Now().AddDate(0, 0, -1)
	fromTimeText := "последние 24 часа"
	var dayIsSelected bool

	var err error
	if cmdArgs != "" {
		args := strings.Split(cmdArgs, " ")
		switch args[0] {
		case "year":
			fromTime = time.Now().AddDate(-1, 0, 0)
			fromTimeText = "последний год"
		case "alltime":
			fromTime = time.Unix(0, 0)
			fromTimeText = "всё время существования бота здесь"
		case "month":
			fromTime = time.Now().AddDate(0, -1, 0)
			fromTimeText = "последний месяц"

		case "week":
			fromTime = time.Now().AddDate(0, 0, -7)
			fromTimeText = "последнюю неделю"

		case "day":
			fromTime = time.Now().AddDate(0, 0, -1)
			fromTimeText = "последние 24 часа"

		default:
			dayIsSelected = true
			pattern := "02.01.2006"
			fromTime, err = time.Parse(pattern, args[0])
			fromTimeText = args[0]
			if err != nil {
				dayIsSelected = false
				fromTime = time.Now().AddDate(0, 0, -1)
				fromTimeText = "последние 24 часа"
			}

		}
	}
	to := time.Now()
	if dayIsSelected {
		to = fromTime.AddDate(0, 0, 1)
	}
	totalMessages, users := CalcUserMessages(fromTime, to, ChatID)

	fileName := fmt.Sprintf("%d-activeStat.png", message.Chat.ID)
	if totalMessages <= 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No messages? o_O")
		bot.Send(msg)
		return
	}
	if true {
		text := GenerateTextStats(users, int(math.Min(15, float64(len(users)))), fromTimeText, totalMessages)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ReplyMarkup = GenerateDeleteKeyboard(message.Chat.ID, message.From.ID)
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	} else {
		RenderActiveUsers(users, fmt.Sprintf(fileName), int(math.Min(15, float64(len(users)))), fromTimeText)
		photo := tgbotapi.FilePath(fileName)
		msg := tgbotapi.NewPhoto(message.Chat.ID, photo)
		msg.Caption = fmt.Sprintf("Написано сообщений за %s \nВсего сообщений: %d\n%v", fromTimeText, totalMessages, time.Now().Sub(startTime))
		msg.ReplyToMessageID = message.MessageID

		msg.ReplyMarkup = GenerateDeleteKeyboard(message.Chat.ID, message.From.ID)
		bot.Send(msg)
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
	smallestNumber := int(math.Min(10, float64(len(wordsFreq))))
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%d самых популярных слов за %s\n", smallestNumber, fromTimeText))

	for i, v := range wordsFreq[:smallestNumber] {
		msg.Text = msg.Text + fmt.Sprintf("%d| %s: %d\n", i, v.Word, v.Freq)
	}
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
