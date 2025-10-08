package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	tokenFlag := flag.String("tg-bot-token", "", "Telegram Bot Token")
	flag.Parse()

	token := *tokenFlag
	if token == "" {
		token = os.Getenv("TELEGRAM_TOKEN")
	}
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не задан")
	}

	bot, err := NewBot(token)
	if err != nil {
		log.Fatalf("Не удалось создать бота: %v", err)
	}

	log.Println("Бот запущен...")
	bot.Run()
}
