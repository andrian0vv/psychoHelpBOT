package main

import (
	"fmt"
	"log"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"dev/tgbot/cmd/service/internal"
	"dev/tgbot/cmd/service/internal/config"
	"dev/tgbot/cmd/service/internal/handlers/cancel"
	_go "dev/tgbot/cmd/service/internal/handlers/go"
	"dev/tgbot/cmd/service/internal/handlers/message"
	"dev/tgbot/cmd/service/internal/handlers/start"
	"dev/tgbot/cmd/service/internal/router"
	"dev/tgbot/cmd/service/internal/storage"
)

func main() {
	fmt.Println("start")
	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgApi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgApi.NewUpdate(0)
	u.Timeout = 60

	storage := storage.NewMemoryStorage()

	router := router.New()
	router.RegisterCommandHandler("start", start.New(cfg, bot, storage))
	router.RegisterCommandHandler("go", _go.New(cfg, bot, storage))
	router.RegisterCommandHandler("cancel", cancel.New(cfg, bot, storage))
	router.RegisterMessageHandler(message.New(cfg, bot, storage))

	app := internal.New(cfg, bot, router)
	app.Start()
}
