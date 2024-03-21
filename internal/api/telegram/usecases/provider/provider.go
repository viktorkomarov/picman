package provider

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	getimage "github.com/viktorkomarov/picman/internal/api/telegram/usecases/get_image"
	listimages "github.com/viktorkomarov/picman/internal/api/telegram/usecases/list_images"
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

func (f *FSMBuilder) GetFSMByCommandType(_type string) *telegram.FSM {
	switch _type {
	case "/upload_image":
		return upload.NewUploadImageFSM(f.repo, f.bot, f.fetcher)
	case "/get_by_name":
		return getimage.NewGetImageFSM(f.repo, f.bot)
	case "/list_images_name":
		return listimages.NewListImagesFSM(f.repo, f.bot)
	default:
		panic("todo")
	}
}
