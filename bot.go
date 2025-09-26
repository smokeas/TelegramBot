package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	// --- сохранение ссылок ---
	if isURL(text) {
		if err := b.store.AddPage(msg.From.ID, text); err != nil {
			b.reply(msg.Chat.ID, "Ошибка при сохранении страницы: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Сохранено! 👌")
		return
	}

	parts := strings.Fields(text)
	if len(parts) == 0 {
		b.reply(msg.Chat.ID, msgUnknownCommand)
		return
	}

	switch parts[0] {
	// --- Pages ---
	case "/rnd":
		p, err := b.store.PickRandomPage(msg.From.ID)
		if err != nil {
			b.reply(msg.Chat.ID, "У тебя нет сохранённых страниц 🙊")
			return
		}
		b.reply(msg.Chat.ID, p)

	// --- TODO ---
	case "/todo_add":
		b.handleTodoAdd(msg)
	case "/todo_list":
		b.handleTodoList(msg)
	case "/todo_done":
		b.handleTodoDone(msg)
	case "/todo_del":
		b.handleTodoDel(msg)

	// --- NOTES ---
	case "/note_add":
		b.handleNoteAdd(msg)
	case "/note_list":
		b.handleNoteList(msg)
	case "/note_del":
		b.handleNoteDel(msg)

	// --- FINANCE ---
	case "/finance_add":
		b.handleFinanceAdd(msg)
	case "/finance_balance":
		b.handleFinanceBalance(msg)
	case "/finance_list":
		b.handleFinanceList(msg)

	default:
		b.reply(msg.Chat.ID, msgUnknownCommand)
	}
}

func (b *Bot) reply(chatID int64, text string) {
	m := tgbotapi.NewMessage(chatID, text)
	// не используем Markdown — отправляем plain text, чтобы "_" не интерпретировались
	// m.ParseMode = "Markdown"  // <- удалить или закомментировать
	if _, err := b.api.Send(m); err != nil {
		log.Printf("Ошибка при отправке: %v", err)
	}
}

func isURL(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
