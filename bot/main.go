package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Println("Токен бота не найден. Убедитесь, что TELEGRAM_BOT_TOKEN указан в .env")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := "Привет! Не можешь определиться с выбором? [Тогда жми сюда](http://scrypze.ru)"
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = "Markdown"

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}
