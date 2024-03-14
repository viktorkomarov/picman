package domain

import "errors"

type Image struct {
	Name    string
	Payload []byte
}

var ErrImageAlreadyExists error = errors.New("image already exists")

type ImageRepository interface {
	SaveImage(image Image) error
	GetByName(name string) Image
	DeleteByName(name string) error
	ListImages() []Image
}
