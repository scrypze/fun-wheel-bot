package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "/start" {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Открыть колесо фортуны", "https://scrypze.ru"),
				),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приступим!")
			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		}
	}
}
