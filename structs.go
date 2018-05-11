package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

// Structure for extending
type Update struct {
	*tgbotapi.Update
}

type User struct {
	ID    int `storm:"unique"`
	Token string
}
