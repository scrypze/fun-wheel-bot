package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/start" {
			log.Println("Received /start command")

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonURL("Перейти на веб-страницу с колесом фортуны", "http://scrypze.ru"),
				},
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите на кнопку ниже, чтобы перейти к колесу фортуны!")
			msg.ReplyMarkup = inlineKeyboard

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
	}
}
