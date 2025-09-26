package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не установлен")
	}

	bot, err := NewBot(token, "data")
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	log.Printf("Authorized on account %s", bot.api.Self.UserName)
	bot.Run()
}
