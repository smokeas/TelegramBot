package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не задан")
	}

	bot, err := NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
