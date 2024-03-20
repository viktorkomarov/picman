package main

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram/bot"
	"github.com/viktorkomarov/picman/internal/api/telegram/executor"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/provider"
	"github.com/viktorkomarov/picman/internal/fetcher"
	"github.com/viktorkomarov/picman/internal/fs"
)

func main() {
	tgbot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		panic(err)
	}
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 30

	fetcher := fetcher.NewTelegramImageFetcher(fetcher.Config{
		Timeout: time.Second * 10,
	}, tgbot)

	dir, err := fs.NewImageRepository("/tmp", "/home/viktor/picman")
	if err != nil {
		panic(err)
	}
	fsmProvider := provider.NewFSMBuilder(tgbot, dir, fetcher)
	userHub := executor.NewUserHub(fsmProvider)

	bot.RunBot(context.Background(), tgbot, bot.RunBotConfig{
		OffsetInitial:       0,
		PollIntervalSeconds: 30,
		MessageHandler:      userHub,
	})
}
