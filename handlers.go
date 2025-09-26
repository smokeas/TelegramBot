package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const helpText = `Доступные команды:
/rnd - случайная сохранённая страница
/todo_add - добавить задачу
/todo_list - список задач
/todo_done - отметить задачу выполненной
/todo_del - удалить задачу
/note_add - добавить заметку
/note_list - список заметок
/note_del - удалить заметку
/finance_add - добавить доход или расход
/finance_balance - показать баланс
/finance_list - показать список операций
`

const msgUnknownCommand = "Неизвестная команда. Напиши /help для списка команд."

// ---- TODO ----
func (b *Bot) handleTodoAdd(msg *tgbotapi.Message) {
	task := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/todo_add"))
	if task == "" {
		b.reply(msg.Chat.ID, "Укажите текст задачи")
		return
	}
	b.store.AddTodo(msg.From.ID, task)
	b.reply(msg.Chat.ID, "Задача добавлена ✅")
}

func (b *Bot) handleTodoList(msg *tgbotapi.Message) {
	tasks := b.store.GetTodos(msg.From.ID)
	if len(tasks) == 0 {
		b.reply(msg.Chat.ID, "Список задач пуст")
		return
	}
	resp := "Задачи:\n"
	for i, t := range tasks {
		status := "❌"
		if t.Done {
			status = "✅"
		}
		resp += fmt.Sprintf("%d. %s %s\n", i+1, status, t.Text)
	}
	b.reply(msg.Chat.ID, resp)
}

func (b *Bot) handleTodoDone(msg *tgbotapi.Message) {
	arg := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/todo_done"))
	i, err := strconv.Atoi(arg)
	if err != nil {
		b.reply(msg.Chat.ID, "Укажите номер задачи")
		return
	}
	if err := b.store.MarkTodoDone(msg.From.ID, i-1); err != nil {
		b.reply(msg.Chat.ID, "Нет такой задачи")
		return
	}
	b.reply(msg.Chat.ID, fmt.Sprintf("Задача %d отмечена выполненной ✅", i))
}

func (b *Bot) handleTodoDel(msg *tgbotapi.Message) {
	arg := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/todo_del"))
	i, err := strconv.Atoi(arg)
	if err != nil {
		b.reply(msg.Chat.ID, "Укажите номер задачи")
		return
	}
	if err := b.store.DeleteTodo(msg.From.ID, i-1); err != nil {
		b.reply(msg.Chat.ID, "Нет такой задачи")
		return
	}
	b.reply(msg.Chat.ID, fmt.Sprintf("Задача %d удалена 🗑️", i))
}

// ---- NOTES ----
func (b *Bot) handleNoteAdd(msg *tgbotapi.Message) {
	note := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/note_add"))
	if note == "" {
		b.reply(msg.Chat.ID, "Укажите текст заметки")
		return
	}
	b.store.AddNote(msg.From.ID, note)
	b.reply(msg.Chat.ID, "Заметка добавлена 📝")
}

func (b *Bot) handleNoteList(msg *tgbotapi.Message) {
	notes := b.store.GetNotes(msg.From.ID)
	if len(notes) == 0 {
		b.reply(msg.Chat.ID, "Список заметок пуст")
		return
	}
	resp := "Заметки:\n"
	for i, n := range notes {
		resp += fmt.Sprintf("%d. %s\n", i+1, n)
	}
	b.reply(msg.Chat.ID, resp)
}

func (b *Bot) handleNoteDel(msg *tgbotapi.Message) {
	arg := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/note_del"))
	i, err := strconv.Atoi(arg)
	if err != nil {
		b.reply(msg.Chat.ID, "Укажите номер заметки")
		return
	}
	if err := b.store.DeleteNote(msg.From.ID, i-1); err != nil {
		b.reply(msg.Chat.ID, "Нет такой заметки")
		return
	}
	b.reply(msg.Chat.ID, fmt.Sprintf("Заметка %d удалена 🗑️", i))
}

// ---- FINANCE ----
func (b *Bot) handleFinanceAdd(msg *tgbotapi.Message) {
	op := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/finance_add"))
	if op == "" {
		b.reply(msg.Chat.ID, "Укажите сумму операции (например: +100 или -50)")
		return
	}

	// Определяем тип
	kind := "income"
	amount, err := strconv.ParseFloat(op, 64)
	if err != nil {
		b.reply(msg.Chat.ID, "Неверный формат. Используй числа, например: 100 или -50")
		return
	}
	if amount < 0 {
		kind = "expense"
		amount = -amount
	}

	if err := b.store.AddFinance(msg.From.ID, kind, amount, ""); err != nil {
		b.reply(msg.Chat.ID, "Ошибка: "+err.Error())
		return
	}
	b.reply(msg.Chat.ID, "Операция добавлена 💰")
}

func (b *Bot) handleFinanceBalance(msg *tgbotapi.Message) {
	income, expense := b.store.FinanceBalance(msg.From.ID)
	balance := income - expense
	resp := fmt.Sprintf("Доход: %.2f\nРасход: %.2f\nБаланс: %.2f 💵", income, expense, balance)
	b.reply(msg.Chat.ID, resp)
}

func (b *Bot) handleFinanceList(msg *tgbotapi.Message) {
	ops := b.store.GetFinance(msg.From.ID)
	if len(ops) == 0 {
		b.reply(msg.Chat.ID, "Операций нет")
		return
	}
	resp := "Операции:\n"
	for i, op := range ops {
		sign := "+"
		if op.Kind == "expense" {
			sign = "-"
		}
		resp += fmt.Sprintf("%d. %s%.2f %s\n", i+1, sign, op.Amount, op.Desc)
	}
	b.reply(msg.Chat.ID, resp)
}
