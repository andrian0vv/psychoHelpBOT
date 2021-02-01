package start

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"dev/tgbot/cmd/service/internal/config"
	"dev/tgbot/cmd/service/internal/handlers"
	"dev/tgbot/cmd/service/internal/questionnaire"
	"dev/tgbot/cmd/service/internal/storage"
)

type Handler struct {
	cfg     *config.Config
	bot     *tgApi.BotAPI
	storage storage.Storage
}

func New(cfg *config.Config, bot *tgApi.BotAPI, storage storage.Storage) handlers.Handler {
	return &Handler{cfg, bot, storage}
}

func (h *Handler) Handle(update tgApi.Update) error {
	return questionnaire.New(h.cfg, h.bot, h.storage).Intro(update)
}
