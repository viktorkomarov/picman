package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/viktorkomarov/picman/internal/domain"
	"github.com/viktorkomarov/picman/internal/fs"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		panic(err)
	}
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 30

	for update := range bot.GetUpdatesChan(cfg) {
		if update.Message == nil {
			continue
		}

		fmt.Println(update.Message.Photo)
	}
}

func download() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		panic(err)
	}

	dir, err := fs.NewImageRepository("/tmp", "/home/viktor/picman")
	if err != nil {
		panic(err)
	}

	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: "",
	})
	if err != nil {
		panic(err)
	}

	fileURL := "" + file.FilePath
	fmt.Println(fileURL)

	resp, err := http.Get(fileURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Header)

	data := bytes.NewBuffer(make([]byte, 0))

	_, err = io.Copy(data, resp.Body)
	if err != nil {
		panic(err)
	}

	rawData, err := io.ReadAll(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(rawData)

	builder := domain.NewImageBuilder()
	if err := builder.SetPayload(rawData); err != nil {
		panic(err)
	}
	if err := builder.SetName("test_3.jpg"); err != nil {
		panic(err)
	}

	if err := dir.SaveImage(builder.Image()); err != nil {
		panic(err)
	}
	fmt.Println("Ok!")
}
