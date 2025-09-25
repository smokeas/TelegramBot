package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const helpText = `Привет! 👋 Я твой персональный бот.

Доступные команды:
/rnd — случайная сохранённая страница
/todo add <текст> — добавить задачу
/todo list — список задач
/todo done <номер> — отметить задачу выполненной
/todo del <номер> — удалить задачу
/note add <текст> — добавить заметку
/note list — список заметок
/note del <номер> — удалить заметку
/finance add income|expense <сумма> <описание> — добавить доход или расход
/finance balance — показать баланс
/finance list — показать список операций
`

const msgUnknownCommand = "Я не понял 🤔 попробуй /help"

func (b *Bot) handleTodo(msg *tgbotapi.Message) {
	parts := strings.Fields(msg.Text)
	if len(parts) < 2 {
		b.reply(msg.Chat.ID, "Формат: /todo [add|list|done|del]")
		return
	}

	cmd := parts[1]
	userID := msg.From.ID

	switch cmd {
	case "add":
		if len(parts) < 3 {
			b.reply(msg.Chat.ID, "Напиши задачу после команды")
			return
		}
		task := strings.Join(parts[2:], " ")
		b.store.AddTodo(userID, task)
		b.reply(msg.Chat.ID, "Добавлено ✅")
	case "list":
		todos := b.store.GetTodos(userID)
		if len(todos) == 0 {
			b.reply(msg.Chat.ID, "Список пуст 👌")
			return
		}
		out := "📝 *Список дел:*\n"
		for i, t := range todos {
			status := "❌"
			if t.Done {
				status = "✅"
			}
			out += fmt.Sprintf("%d. %s %s\n", i+1, status, t.Text)
		}
		b.reply(msg.Chat.ID, out)
	case "done":
		if len(parts) < 3 {
			b.reply(msg.Chat.ID, "Формат: /todo done <номер>")
			return
		}
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			b.reply(msg.Chat.ID, "Некорректный номер")
			return
		}
		if err := b.store.MarkTodoDone(userID, idx-1); err != nil {
			b.reply(msg.Chat.ID, "Ошибка: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Задача выполнена ✅")
	case "del":
		if len(parts) < 3 {
			b.reply(msg.Chat.ID, "Формат: /todo del <номер>")
			return
		}
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			b.reply(msg.Chat.ID, "Некорректный номер")
			return
		}
		if err := b.store.DeleteTodo(userID, idx-1); err != nil {
			b.reply(msg.Chat.ID, "Ошибка: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Удалено 🗑")
	default:
		b.reply(msg.Chat.ID, "Неизвестная команда")
	}
}

func (b *Bot) handleNote(msg *tgbotapi.Message) {
	parts := strings.Fields(msg.Text)
	if len(parts) < 2 {
		b.reply(msg.Chat.ID, "Формат: /note [add|list|del]")
		return
	}

	cmd := parts[1]
	userID := msg.From.ID

	switch cmd {
	case "add":
		if len(parts) < 3 {
			b.reply(msg.Chat.ID, "Напиши заметку после команды")
			return
		}
		note := strings.Join(parts[2:], " ")
		b.store.AddNote(userID, note)
		b.reply(msg.Chat.ID, "Сохранено 📝")
	case "list":
		notes := b.store.GetNotes(userID)
		if len(notes) == 0 {
			b.reply(msg.Chat.ID, "Заметок нет 👌")
			return
		}
		out := "🗒 *Заметки:*\n"
		for i, n := range notes {
			out += fmt.Sprintf("%d. %s\n", i+1, n)
		}
		b.reply(msg.Chat.ID, out)
	case "del":
		if len(parts) < 3 {
			b.reply(msg.Chat.ID, "Формат: /note del <номер>")
			return
		}
		idx, err := strconv.Atoi(parts[2])
		if err != nil {
			b.reply(msg.Chat.ID, "Некорректный номер")
			return
		}
		if err := b.store.DeleteNote(userID, idx-1); err != nil {
			b.reply(msg.Chat.ID, "Ошибка: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Удалено 🗑")
	default:
		b.reply(msg.Chat.ID, "Неизвестная команда")
	}
}

func (b *Bot) handleFinance(msg *tgbotapi.Message) {
	parts := strings.Fields(msg.Text)
	if len(parts) < 2 {
		b.reply(msg.Chat.ID, "Формат: /finance [add|balance|list]")
		return
	}

	cmd := parts[1]
	userID := msg.From.ID

	switch cmd {
	case "add":
		if len(parts) < 4 {
			b.reply(msg.Chat.ID, "Формат: /finance add income|expense <сумма> <описание>")
			return
		}
		kind := parts[2]
		amount, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			b.reply(msg.Chat.ID, "Некорректная сумма")
			return
		}
		desc := ""
		if len(parts) > 4 {
			desc = strings.Join(parts[4:], " ")
		}
		if err := b.store.AddFinance(userID, kind, amount, desc); err != nil {
			b.reply(msg.Chat.ID, "Ошибка: "+err.Error())
			return
		}
		b.reply(msg.Chat.ID, "Операция сохранена 💰")
	case "balance":
		inc, exp := b.store.FinanceBalance(userID)
		bal := inc - exp
		out := fmt.Sprintf("Доходы: %.2f\nРасходы: %.2f\nБаланс: %.2f", inc, exp, bal)
		b.reply(msg.Chat.ID, out)
	case "list":
		tx := b.store.GetFinance(userID)
		if len(tx) == 0 {
			b.reply(msg.Chat.ID, "Операций нет 👌")
			return
		}
		out := "💰 *Финансы:*\n"
		for i, t := range tx {
			out += fmt.Sprintf("%d. [%s] %.2f — %s\n", i+1, t.Kind, t.Amount, t.Desc)
		}
		b.reply(msg.Chat.ID, out)
	default:
		b.reply(msg.Chat.ID, "Неизвестная команда")
	}
}
