package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramImageFetcher struct {
	client  http.Client
	baseUrl string
	bot     *tgbotapi.BotAPI
}

type Config struct {
	Timeout time.Duration
}

func NewTelegramImageFetcher(cfg Config, bot *tgbotapi.BotAPI) *TelegramImageFetcher {
	return &TelegramImageFetcher{
		client: http.Client{
			Timeout: cfg.Timeout,
		},
		baseUrl: fmt.Sprintf("https://api.telegram.org/file/bot%s/", bot.Token),
		bot:     bot,
	}
}

func (t *TelegramImageFetcher) Fetch(ctx context.Context, fileID string) ([]byte, error) {
	fileInfo, err := t.bot.GetFile(tgbotapi.FileConfig{
		FileID: fileID,
	})
	if err != nil {
		return nil, fmt.Errorf("bot.GetFile: %w", err)
	}
	fileURL := t.baseUrl + fileInfo.FilePath

	resp, err := t.client.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}
	defer resp.Body.Close()

	bufferReader := bytes.NewBuffer(make([]byte, 0))
	if _, err := io.Copy(bufferReader, resp.Body); err != nil {
		return nil, fmt.Errorf("io.Copy: %w", err)
	}

	return bufferReader.Bytes(), nil
}
