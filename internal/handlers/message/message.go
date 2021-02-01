package message

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Andrianov/psychoHelpBOT/internal/config"
	"github.com/Andrianov/psychoHelpBOT/internal/handlers"
	"github.com/Andrianov/psychoHelpBOT/internal/questionnaire"
	"github.com/Andrianov/psychoHelpBOT/internal/storage"
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
	return questionnaire.New(h.cfg, h.bot, h.storage).Continue(update)
}
