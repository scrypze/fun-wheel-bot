package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			wheelURL := fmt.Sprintf("http://scrypze.ru/wheel/%d", update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, 
				fmt.Sprintf("Создайте своё колесо фортуны: %s", wheelURL))
			bot.Send(msg)
		}
	}
}
