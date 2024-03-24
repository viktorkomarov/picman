package main

import (
	"flag"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/fetcher"
)

type FSConfig struct {
	TMPPath  string
	BasePath string
}

type TgBotConfig struct {
	Token        string
	UpdateConfig tgbotapi.UpdateConfig
}

type AppConfig struct {
	FS           FSConfig
	ImageFetcher fetcher.Config
	TgBot        TgBotConfig
}

func mustParseConfig() AppConfig {
	imageFetcher := fetcher.Config{
		Timeout: time.Second * 5,
	}

	token, ok := os.LookupEnv("TG_PICMAN_TOKEN")
	if !ok {
		panic("failed to get TG_PICMAN_TOKEN")
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	tgBotConfig := TgBotConfig{
		Token:        token,
		UpdateConfig: updateConfig,
	}

	var (
		tmpPath  string
		basePath string
	)
	flag.StringVar(&tmpPath, "tmp", "", "means tmp dir for atomic operations with files")
	flag.StringVar(&basePath, "base", "", "means main dir where files are stored")

	flag.Parse()

	fsConfig := FSConfig{
		TMPPath:  tmpPath,
		BasePath: basePath,
	}

	return AppConfig{
		FS:           fsConfig,
		TgBot:        tgBotConfig,
		ImageFetcher: imageFetcher,
	}
}
