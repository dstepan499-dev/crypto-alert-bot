package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/dstepan499-dev/crypto-alert-bot/internal/coingecko"
	"github.com/dstepan499-dev/crypto-alert-bot/internal/storage"
)

func main() {
	// Loading variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env was not found")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN was not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Bot initialization fault:", err)
	}

	log.Printf("Authorized as bot: @%s", bot.Self.UserName)

	cryptoClient := coingecko.NewClient(10 * time.Second)
	db, err := storage.NewStorage("alerts.db")
	if err != nil {
		log.Fatalf("SQLite connection error: %v", err)
	}

	go startAlertChecker(bot, db, cryptoClient)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msgText := "Hello!\n\n" +
					"Available commands:\n" +
					"🔹 `/price btc` — find out the current price\n" +
					"🔹 `/alert btc 65000` — set notification if the price decreases/increases to target value"
				msg := tgbotapi.NewMessage(chatID, msgText)
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "price":
				args := strings.Fields(update.Message.CommandArguments())
				if len(args) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Example: `/price btc`"))
					continue
				}
				price, err := cryptoClient.GetPrice(args[0], "usd")
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err)))
					continue
				}
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Price %s: $%.2f", strings.ToUpper(args[0]), price)))

			case "alert":
				args := strings.Fields(update.Message.CommandArguments())
				if len(args) < 2 {
					msg := tgbotapi.NewMessage(chatID, "Format: `/alert <coin> <target_price>`\nExample: `/alert btc 65000`")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					continue
				}

				coin := args[0]
				targetPrice, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Incorrect price format. Enter the number."))
					continue
				}

				if err := db.AddAlert(chatID, coin, targetPrice); err != nil {
					log.Printf("Alert save error: %v", err)
					bot.Send(tgbotapi.NewMessage(chatID, "Didn't manage to save notification"))
					continue
				}

				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Notification was set! I will tell you when %s reaches $%.2f", strings.ToUpper(coin), targetPrice)))
			}
		}
	}
}

func startAlertChecker(bot *tgbotapi.BotAPI, db *storage.Storage, cryptoClient *coingecko.Client) {
	ticker := time.NewTicker(30 * time.Second) // опрос каждые 30 секунд
	for range ticker.C {
		alerts, err := db.GetAlerts()
		if err != nil {
			log.Printf("Worker: alert getting error: %v", err)
			continue
		}

		for _, alert := range alerts {
			currentPrice, err := cryptoClient.GetPrice(alert.Coin, "usd")
			if err != nil {
				continue
			}

			// Main Logic is that if the price is near target price within 1%
			diff := (currentPrice - alert.TargetPrice) / alert.TargetPrice
			if diff < 0 {
				diff = -diff
			}

			if diff <= 0.01 {
				text := fmt.Sprintf("*THE ALERT WAS TRIGGERED!*\n\nCoin *%s* reached target mark!\nCurrent price: **$%.2f** (Target: $%.2f)",
					strings.ToUpper(alert.Coin), currentPrice, alert.TargetPrice)
				msg := tgbotapi.NewMessage(alert.ChatID, text)
				msg.ParseMode = "Markdown"
				bot.Send(msg)

				_ = db.DeleteAlert(alert.ID)
			}
		}
	}
}
