package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram/bot"
	"github.com/viktorkomarov/picman/internal/api/telegram/executor"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/provider"
	"github.com/viktorkomarov/picman/internal/fetcher"
	"github.com/viktorkomarov/picman/internal/fs"
)

func main() {
	appConfig := mustParseConfig()

	tgbot, err := tgbotapi.NewBotAPI(appConfig.TgBot.Token)
	if err != nil {
		log.Fatalf("failed to create bot: %s", err.Error())
	}
	fetcher := fetcher.NewTelegramImageFetcher(appConfig.ImageFetcher, tgbot)

	dir, err := fs.NewImageRepository(appConfig.FS.TMPPath, appConfig.FS.BasePath)
	if err != nil {
		log.Fatalf("failed to image repo: %s", err.Error())
	}
	fsmProvider := provider.NewFSMBuilder(tgbot, dir, fetcher)
	userHub := executor.NewUserHub(fsmProvider)

	bot.RunBot(context.Background(), appConfig.TgBot.UpdateConfig, tgbot, userHub)
}
