package main

import (
	"context"
	pb "fun-wheel-bot/grpc"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/grpc"
)

const token = "7733625826:AAHYWZzBRUNUTFfUGUaGdCK0iGXFScNmy-s"

func main() {
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

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFunWheelServiceClient(conn)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)

		// Команда /start
		if update.Message.Text == "/start" {
			log.Println("Received /start command")

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonURL("Перейти на веб-страницу с колесом фортуны", "http://localhost:8080/wheel"),
				},
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите на кнопку ниже, чтобы перейти к колесу фортуны!")
			msg.ReplyMarkup = inlineKeyboard

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

		if update.Message.Text == "/createwheel" {
			resp, err := client.CreateWheel(context.Background(), &pb.CreateWheelRequest{
				ChatId: update.Message.Chat.ID,
			})
			if err != nil {
				log.Printf("Error creating wheel: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при создании колеса"))
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.GetMessage()))
		}

		if update.Message.Text == "/addoption" {
			resp, err := client.AddOption(context.Background(), &pb.AddOptionRequest{
				ChatId: update.Message.Chat.ID,
				Option: "option1",
			})
			if err != nil {
				log.Printf("Error adding option: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при добавлении опции"))
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp.GetMessage()))
		}

		if update.Message.Text == "/spinwheel" {
			resp, err := client.SpinWheel(context.Background(), &pb.SpinWheelRequest{
				ChatId: update.Message.Chat.ID,
			})
			if err != nil {
				log.Printf("Error spinning wheel: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при вращении колеса"))
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Результат: "+resp.GetResult()))
		}
	}
}
