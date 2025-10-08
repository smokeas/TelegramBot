package main

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	b.states[userID] = "" // сбрасываем предыдущее состояние
	args := update.Message.CommandArguments()
	msg := tgbotapi.NewMessage(userID, "")
	switch update.Message.Command() {
	case "todo_add":
		if args != "" {
			id, err := b.store.AddTask(userID, args)
			if err != nil {
				msg.Text = "Ошибка при добавлении задачи"
			} else {
				msg.Text = fmt.Sprintf("Задача %d добавлена.", id)
			}
		} else {
			msg.Text = "Укажи текст задачи"
			b.states[userID] = "todo_add"
		}
	case "todo_done":
		if args != "" {
			num, err := strconv.Atoi(args)
			if err != nil {
				msg.Text = "Номер задачи должен быть числом"
			} else {
				err := b.store.MarkDone(userID, num)
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = fmt.Sprintf("Задача %d отмечена как выполненная.", num)
				}
			}
		} else {
			msg.Text = "Укажи номер задачи для пометки выполненной"
			b.states[userID] = "todo_done"
		}
	case "todo_del":
		if args != "" {
			num, err := strconv.Atoi(args)
			if err != nil {
				msg.Text = "Номер задачи должен быть числом"
			} else {
				err := b.store.DeleteTask(userID, num)
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = fmt.Sprintf("Задача %d удалена.", num)
				}
			}
		} else {
			msg.Text = "Укажи номер задачи для удаления"
			b.states[userID] = "todo_del"
		}
	default:
		msg.Text = "Неизвестная команда"
	}
	b.api.Send(msg)
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := b.states[userID]
	text := update.Message.Text
	msg := tgbotapi.NewMessage(userID, "")
	switch state {
	case "todo_add":
		id, err := b.store.AddTask(userID, text)
		if err != nil {
			msg.Text = "Ошибка при добавлении задачи"
		} else {
			msg.Text = fmt.Sprintf("Задача %d добавлена.", id)
		}
	case "todo_done":
		num, err := strconv.Atoi(text)
		if err != nil {
			msg.Text = "Номер задачи должен быть числом"
		} else {
			err := b.store.MarkDone(userID, num)
			if err != nil {
				msg.Text = err.Error()
			} else {
				msg.Text = fmt.Sprintf("Задача %d отмечена как выполненной.", num)
			}
		}
	case "todo_del":
		num, err := strconv.Atoi(text)
		if err != nil {
			msg.Text = "Номер задачи должен быть числом"
		} else {
			err := b.store.DeleteTask(userID, num)
			if err != nil {
				msg.Text = err.Error()
			} else {
				msg.Text = fmt.Sprintf("Задача %d удалена.", num)
			}
		}
	default:
		return // если не в состоянии ожидания, игнорируем сообщение
	}
	b.states[userID] = "" // сбрасываем состояние после обработки
	b.api.Send(msg)
}
