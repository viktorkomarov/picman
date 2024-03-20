package domain

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	formatGIF  = "gif"
	formatPNG  = "png"
	formatJPEG = "jpg"
)

var allFormats = []string{
	formatGIF, formatPNG, formatJPEG,
}

var ErrIncorrectName error = fmt.Errorf("name format must be: __.[%s]", strings.Join(allFormats, ","))

type Image struct {
	Name    string
	Payload []byte
}

type ImageBuilder struct {
	image Image
}

func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{}
}

func (i *ImageBuilder) SetName(name string) error {
	count := 0
	for _, suffix := range allFormats {
		if strings.HasSuffix(name, "."+suffix) {
			count++
		}
	}
	if count == 1 {
		i.image.Name = name
		return nil
	}
	return ErrIncorrectName
}

func (i *ImageBuilder) SetPayload(payload []byte) error {
	_, _, err := image.Decode(bytes.NewReader(payload))
	if err != nil {
		return err
	}
	i.image.Payload = payload
	return nil
}

func (i *ImageBuilder) Image() Image {
	return i.image
}

var ErrImageAlreadyExists error = errors.New("image already exists")

type ImageRepository interface {
	SaveImage(image Image) error
	GetByName(name string) (Image, error)
	DeleteByName(name string) error
	ListImages() ([]Image, error)
}
