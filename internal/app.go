package internal

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Andrianov/psychoHelpBOT/internal/config"
	"github.com/Andrianov/psychoHelpBOT/internal/router"
)

type App struct {
	cfg    *config.Config
	bot    *tgApi.BotAPI
	router router.Router
}

func New(cfg *config.Config, bot *tgApi.BotAPI, router router.Router) *App {
	return &App{cfg, bot, router}
}

func (a *App) Start() {
	u := tgApi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		a.notify(err)
		return
	}

	for update := range updates {
		if err := a.handleUpdate(update); err != nil {
			go a.notify(err)
		}
	}
}

func (a *App) handleUpdate(update tgApi.Update) error {
	defer func() {
		if err := recover(); err != nil {
			go a.notify(err)
		}
	}()

	if update.Message != nil && update.Message.IsCommand() {
		if handler, exists := a.router.GetCommandHandler(update.Message.Command()); exists {
			return handler.Handle(update)
		}

		// some text with "i dont know what's going on"
		return nil
	}

	return a.router.GetMessageHandler().Handle(update)
}

func (a *App) notify(err interface{}) {
	var text string
	switch v := err.(type) {
	case string:
		text = v
	case error:
		text = v.Error()
	default:
		return
	}

	msg := tgApi.NewMessage(a.cfg.TechChatID, text)
	_, _ = a.bot.Send(msg)
}
