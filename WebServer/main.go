package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Status string `json:"status"`
}

type Response struct {
	Text    string `json:"text"`
	Message string `json:"message"` // ✅ Исправлено написание
}

var messages []Message // ✅ Правильное объявление переменной

func GetHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &messages)
}

func PostHandler(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Text:    "Error",
			Message: "Could not add the message", // ✅ Добавлено двоеточие
		})
	}
	messages = append(messages, message) // ✅ Используем правильное имя переменной
	return c.JSON(http.StatusOK, Response{
		Text:    "OK",
		Message: "Message added successfully",
	})
}

func main() {
	e := echo.New()
	e.GET("/get", GetHandler)
	e.POST("/post", PostHandler)

	e.Start(":8080")
}
