package main

import (
	"flag"
	"log"
	"os"
)

//cd D:\MyProjects\TelegramBot\tg_todolist_bot
//dir
//go mod tidy
//go build -o bot.exe
//.\bot.exe -tg-bot-token "7982922881:AAGDbLHmt3csIbianT_j1PJ-83jBzWjX07g"

func main() {
	// флаг для токена
	tgTokenFlag := flag.String("tg-bot-token", "", "Telegram bot token")
	flag.Parse()

	// 1) сначала смотрим флаг, 2) если пустой — переменная окружения
	token := *tgTokenFlag
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}

	if token == "" {
		log.Fatal("bot token not provided: use -tg-bot-token or set TELEGRAM_BOT_TOKEN")
	}

	bot, err := NewBot(token, "data")
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	log.Println("Bot started")
	bot.Run()
}
