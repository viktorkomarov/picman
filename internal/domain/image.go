package domain

type Image struct {
	Name    string
	Payload []byte
}

type ImageRepository interface {
	SaveImage(image Image) error
	GetByName(name string) Image
	DeleteByName(name string) error
	ListImages() []Image
}
