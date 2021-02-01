package handlers

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler interface {
	Handle(update tgApi.Update) error
}

