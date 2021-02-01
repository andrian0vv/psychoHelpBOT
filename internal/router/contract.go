package router

import (
	"fmt"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"dev/tgbot/cmd/service/internal/handlers"
)

type Router interface {
	RegisterMessageHandler(handler Handler)
	GetMessageHandler() Handler
	RegisterCommandHandler(command string, handler Handler)
	GetCommandHandler(command string) (Handler, bool)
}

type Handler = handlers.Handler

type DefaultHandler struct {}

func (h *DefaultHandler) Handle(update tgApi.Update) error {
	fmt.Println("default handler on", update.Message.From, update.Message.Text)
	return nil
}
