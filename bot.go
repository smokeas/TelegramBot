package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//mkdir tg_todolist_bot
//cd tg_todolist_bot
//go mod tidy

type Bot struct {
	api     *tgbotapi.BotAPI
	dataDir string
	store   *Store
}

func NewBot(token, dataDir string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = false

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	s := NewStore(filepath.Join(dataDir, "users"))

	return &Bot{api: api, dataDir: dataDir, store: s}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	for upd := range updates {
		if upd.Message == nil {
			continue
		}
		go b.handleMessage(upd.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	text := strings.TrimSpace(msg.Text)

	if text == "/start" || text == "/help" {
		b.reply(msg.Chat.ID, helpText)
		return
	}

	// URL saving (оригинальный функционал)
	if isURL(text) {
		if err := b.store.AddPage(msg.From.ID, text); err != nil {
			b.reply(msg.Chat.ID, "Error saving page: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Saved! 👌")
		return
	}

	parts := strings.Fields(text)
	if len(parts) == 0 {
		b.reply(msg.Chat.ID, "Unknown command")
		return
	}

	switch parts[0] {
	case "/rnd":
		p, err := b.store.PickRandomPage(msg.From.ID)
		if err != nil {
			b.reply(msg.Chat.ID, "You have no saved pages 🙊")
			return
		}
		b.reply(msg.Chat.ID, p)
	case "/todo":
		b.handleTodo(msg)
	case "/note":
		b.handleNote(msg)
	case "/finance":
		b.handleFinance(msg)
	default:
		b.reply(msg.Chat.ID, msgUnknownCommand)
	}
}

func (b *Bot) reply(chatID int64, text string) {
	m := tgbotapi.NewMessage(chatID, text)
	m.ParseMode = "Markdown"
	if _, err := b.api.Send(m); err != nil {
		log.Printf("send error: %v", err)
	}
}

func isURL(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
