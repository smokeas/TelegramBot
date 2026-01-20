package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	store *Store
	state map[int64]string
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:   api,
		store: NewStore("data"),
		state: make(map[int64]string),
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(m *tgbotapi.Message) {
	userID := m.Chat.ID
	text := strings.TrimSpace(m.Text)

	// Если пользователь в состоянии — ждём текст
	if state, ok := b.state[userID]; ok {
		b.handleUserState(userID, state, text)
		delete(b.state, userID)
		return
	}

	// Команды
	if !m.IsCommand() {
		return
	}

	switch m.Command() {
	case "start", "help":
		b.send(userID, b.helpText())

	case "todo_add":
		b.state[userID] = "todo_add"
		b.send(userID, "Введите задачу:")

	case "todo_list":
		list := b.store.ListTodos(userID)
		if list == "" {
			b.send(userID, "Список задач пуст.")
		} else {
			b.send(userID, list)
		}

	case "todo_done":
		b.state[userID] = "todo_done"
		b.send(userID, "Введите номер выполненной задачи:")

	case "todo_del":
		b.state[userID] = "todo_del"
		b.send(userID, "Введите номер задачи для удаления:")

	case "note_add":
		b.state[userID] = "note_add"
		b.send(userID, "Введите заметку:")

	case "note_list":
		list := b.store.ListNotes(userID)
		if list == "" {
			b.send(userID, "Список заметок пуст.")
		} else {
			b.send(userID, list)
		}

	case "note_del":
		b.state[userID] = "note_del"
		b.send(userID, "Введите номер заметки для удаления:")

	case "finance_add":
		b.state[userID] = "finance_add"
		b.send(userID, "Введите сумму и комментарий (например: +500 зарплата или -200 еда):")

	case "finance_list":
		b.send(userID, b.store.ListFinance(userID))

	case "finance_balance":
		b.send(userID, b.store.Balance(userID))

	case "rnd":
		b.send(userID, b.store.Random(userID))

	default:
		b.send(userID, "Неизвестная команда. Введите /help")
	}
}

func (b *Bot) handleUserState(userID int64, state, text string) {
	switch state {
	case "todo_add":
		b.store.AddTodo(userID, text)
		b.send(userID, "Задача добавлена ✅")

	case "todo_done":
		b.store.DoneTodo(userID, text)
		b.send(userID, "Задача отмечена как выполненная ✅")

	case "todo_del":
		b.store.DeleteTodo(userID, text)
		b.send(userID, "Задача удалена 🗑️")

	case "note_add":
		b.store.AddNote(userID, text)
		b.send(userID, "Заметка сохранена 📝")

	case "note_del":
		b.store.DeleteNote(userID, text)
		b.send(userID, "Заметка удалена 🗑️")

	case "finance_add":
		b.store.AddFinance(userID, text)
		b.send(userID, "Финансовая запись добавлена 💰")
	}
}

func (b *Bot) helpText() string {
	return strings.Join([]string{
		"/todo_add",
		"/todo_list",
		"/todo_done",
		"/todo_del",
		"/note_add",
		"/note_list",
		"/note_del",
		"/finance_add",
		"/finance_list",
		"/finance_balance",
		"/rnd",
	}, "\n")
}

func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
