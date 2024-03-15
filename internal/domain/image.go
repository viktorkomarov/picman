package domain

import "errors"

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
	return nil
}

func (i *ImageBuilder) SetPayload(payload []byte) error {
	return nil
}

func (i *ImageBuilder) Image() Image {
	return i.image
}

var ErrImageAlreadyExists error = errors.New("image already exists")

type ImageRepository interface {
	SaveImage(image Image) error
	GetByName(name string) Image
	DeleteByName(name string) error
	ListImages() []Image
}
