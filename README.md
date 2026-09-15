# Crypto & Currency Alert Telegram Bot

A lightweight and efficient Telegram Bot written in **Go (Golang)** for tracking cryptocurrency prices and receiving price threshold alerts.

---

## Features
- **Real-time Prices:** Fetches live cryptocurrency data via CoinGecko REST API (`/price btc`, `/price eth`).
- **Price Alerts:** Set target price notifications (`/alert btc 65000`).
- **Background Worker:** Automated price checking worker powered by Go concurrency (Goroutines & Tickers).
- **Persistent Storage:** SQLite database for storing active alerts and user subscriptions.
- **Docker Ready:** Fully containerized with a lightweight multi-stage Docker build (~20MB image).

---

## Tech Stack
- **Language:** Go 1.23+
- **Database:** SQLite (`modernc.org/sqlite` pure Go driver)
- **API:** Telegram Bot API (`go-telegram-bot-api`), CoinGecko Public API
- **DevOps:** Docker, Docker Compose

---

## Quick Start

1. Clone the repository:
git clone https://github.com/dstepan499-dev/crypto-alert-bot.git
cd crypto-alert-bot

2. Configure Environment:
Create a .env file from .env.example (cp .env.example .env) and set your Telegram Bot token from @BotFather:
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here

3. Run with Docker (Recommended):
docker compose up -d --build

4. Run Locally:
go run cmd/bot/main.go

---

## Commands

- /start - Displays welcome message and basic instructions

- /price <coin> - Returns current price in USD (e.g., /price btc, /price sol)

- /alert <coin> <target_price> - Sets an alert notification (e.g., /alert btc 65000)
