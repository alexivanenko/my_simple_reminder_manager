package main

import (
	"log"

	"time"

	"github.com/alexivanenko/my_simple_reminder_manager/config"
	"github.com/alexivanenko/my_simple_reminder_manager/model"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var TBot *tgbotapi.BotAPI

func main() {
	defer model.GetSession().Close()

	bot, err := tgbotapi.NewBotAPI(config.String("bot", "token"))
	if err != nil {
		log.Panic(err)
	}

	TBot = bot

	if config.Is("bot", "debug") {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	doEvery(1*time.Minute, sendNotifications)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func sendNotifications(t time.Time) {
	event := new(model.Event)
	events, err := event.LoadAll()

	if err == nil {
		utc, _ := time.LoadLocation("UTC")
		nowFormatted := time.Now().In(utc).Format("02/01/2006 15:04") + ":00"
		now, _ := time.ParseInLocation("02/01/2006 15:04:05", nowFormatted, utc)

		for _, item := range events {

			loc, locErr := time.LoadLocation(item.TimeZone)

			if locErr != nil {
				continue
			}

			if now.In(loc).Equal(item.Date.In(loc)) {
				msg := tgbotapi.NewMessage(item.ChatID, "Don't forget about "+item.Name)
				TBot.Send(msg)

				//item.Remove()
			}
		}
	}
}
