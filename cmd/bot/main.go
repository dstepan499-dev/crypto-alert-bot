package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/dstepan499-dev/crypto-alert-bot/internal/coingecko"
)

func main() {
	// Loading variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env-file was not found, read environmental variables")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN was not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Bot initialization fault", err)
	}

	bot.Debug = true
	log.Printf("Authorized as bot: @%s", bot.Self.UserName)

	cryptoClient := coingecko.NewClient(10 * time.Second)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Listening to incoming messages
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID,
					"Hello!\n\n"+"I can show you cryptocurrency prices.\n"+
						"Use command:\n"+
						"`/price btc` or `/price eth usd`")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			case "price":
				args := strings.Fields(update.Message.CommandArguments())
				if len(args) == 0 {
					msg := tgbotapi.NewMessage(chatID, "Select coin. Example: `/price btc` or `/price sol usd`")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					continue
				}

				coin := args[0]
				currency := "usd"
				if len(args) > 1 {
					currency = args[1]
				}

				price, err := cryptoClient.GetPrice(coin, currency)
				if err != nil {
					log.Printf("Price receiving error: %v", err)
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
					bot.Send(msg)
					continue
				}

				replyText := fmt.Sprintf("Price *%s*: **$%.2f** (%s)", strings.ToUpper(coin), price, strings.ToUpper(currency))
				msg := tgbotapi.NewMessage(chatID, replyText)
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Simple echo-answer for "/start"
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! Bot started and ready to work")

			bot.Send(msg)
		}
	}
}
