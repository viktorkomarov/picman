package provider

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/api/telegram"
	deleteimage "github.com/viktorkomarov/picman/internal/api/telegram/usecases/delete_image"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/fallback"
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
		return upload.NewFSM(f.repo, f.bot, f.fetcher)
	case "/get_by_name":
		return getimage.NewFSM(f.repo, f.bot)
	case "/list_images_name":
		return listimages.NewFSM(f.repo, f.bot)
	case "/delete_by_name":
		return deleteimage.NewDeleteImageFSM(f.repo, f.bot)
	default:
		return fallback.NewFSM(f.bot)
	}
}
