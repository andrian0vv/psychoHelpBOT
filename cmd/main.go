package main

import (
	"fmt"
	"log"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Andrianov/psychoHelpBOT/internal"
	"github.com/Andrianov/psychoHelpBOT/internal/config"
	"github.com/Andrianov/psychoHelpBOT/internal/handlers/message"
	"github.com/Andrianov/psychoHelpBOT/internal/handlers/start"
	"github.com/Andrianov/psychoHelpBOT/internal/router"
	"github.com/Andrianov/psychoHelpBOT/internal/storage"
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

	storage := storage.NewMemoryStorage()

	router := router.New()
	router.RegisterCommandHandler("start", start.New(cfg, bot, storage))
	//router.RegisterCommandHandler("go", _go.New(cfg, bot, storage))
	//router.RegisterCommandHandler("cancel", cancel.New(cfg, bot, storage))
	router.RegisterMessageHandler(message.New(cfg, bot, storage))

	app := internal.New(cfg, bot, router)
	app.Start()
}
