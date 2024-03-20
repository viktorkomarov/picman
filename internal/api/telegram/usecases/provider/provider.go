package provider

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/upload"
	"github.com/viktorkomarov/picman/internal/domain"
	"github.com/viktorkomarov/picman/internal/fetcher"
)

type FSMBuilder struct {
	bot     *tgbotapi.BotAPI
	repo    domain.ImageRepository
	fetcher *fetcher.TelegramImageFetcher
}

func NewFSMBuilder(bot *tgbotapi.BotAPI, repo domain.ImageRepository, fetcher *fetcher.TelegramImageFetcher) *FSMBuilder {
	return &FSMBuilder{
		bot:     bot,
		repo:    repo,
		fetcher: fetcher,
	}
}

func (f *FSMBuilder) GetFSMByCommandType(_type telegram.FSMType) *telegram.FSM {
	switch _type {
	case telegram.FSMTypeUpload:
		return upload.NewUploadImageFSM(f.repo, f.bot, f.fetcher)
	default:
		panic("todo")
	}
}
