package main

import (
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"runtime"
	"statBot/utils"
	"statBot/utils/workerpool"
	"strconv"
	"strings"
)

var (
	allowedChatsMode bool
	weNeedToDie      chan bool
	BlacklistedUsers = map[int64]bool{5449020876: true}
)

// 1 Day = 86400 sec
func init() {
	utils.Handlers = make(map[string]utils.Handler)
}
func parseFlags() {
	var allowedChats, blacklistedUsers string
	defaultAllowedChats := "559723688,-1001549183364,-749918079,-1001373811109,-1001558727831,-1001740354030,-1001053617676,-1001386313371"
	defaultBlackListed := "5449020876"

	flag.StringVar(&allowedChats, "allowedChats", defaultAllowedChats, "allowedChats")
	flag.StringVar(&blacklistedUsers, "blacklistedUsers", defaultBlackListed, "blacklistedUsers")
	flag.BoolVar(&allowedChatsMode, "allowedChatsMode", true, "allowedChatsMode")
	flag.Parse()

	splittedChats := strings.Split(allowedChats, ",")
	splittedUsers := strings.Split(blacklistedUsers, ",")

	for i, str := range splittedChats {
		fmt.Println(i, "Adding", str)
		chatId, err := strconv.ParseInt(str, 10, 64)
		utils.PanicErr(err)
		utils.AllowedChats[chatId] = true
	}
	for i, str := range splittedUsers {
		fmt.Println(i, "Blacklisting", str)
		userId, err := strconv.ParseInt(str, 10, 64)
		utils.PanicErr(err)
		BlacklistedUsers[userId] = true
	}
}

func main() {
	parseFlags()
	bot := InitBot()
	var err error
	dsn := "host=localhost user=postgres password=dumpdatabase dbname=rstatbot port=5432 sslmode=disable TimeZone=Europe/Moscow"
	//utils.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	utils.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB := utils.DB
	x, err := DB.DB()
	x.SetMaxOpenConns(runtime.NumCPU() - 1)
	utils.PanicErr(err)
	err = DB.AutoMigrate(&utils.Chat{}, &utils.User{}, &utils.ChatMessage{}, &utils.ChatAudio{})
	utils.PanicErr(err)
	for i, value := range utils.AllowedChats {
		log.Printf("AllowedChat %d: %t", i, value)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	baseObj := utils.SharedBaseObject{
		Bot:              bot,
		DB:               DB,
		AllowedChatsMode: allowedChatsMode,
	}
	log.Println("AllowedChatsMode enabled:", baseObj.AllowedChatsMode)
	for w := 1; w <= runtime.NumCPU()-1; w++ {
		log.Println("starting worker with id", w)
		workerpool.WorkerChanPool[w] = make(chan utils.ControlStruct)
		go workerpool.UpdateWorker(w, &baseObj, updates, workerpool.WorkerChanPool[w])
	}
	//workerpool.WorkerChanPool[1] <- utils.ControlStruct{
	//	Cmd:  "Kill",
	//	Args: "",
	//}
	<-weNeedToDie

}
