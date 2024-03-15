package upload

import (
	"github.com/viktorkomarov/picman/internal/domain"
)

type saver interface {
	SaveImage(image domain.Image) error
}

type UploadImageUseCase struct {
	saver        saver
	imageBuilder *domain.ImageBuilder
}

func NewUploadImageUseCase(saver saver) *UploadImageUseCase {
	return &UploadImageUseCase{
		saver:        saver,
		imageBuilder: domain.NewImageBuilder(),
	}
}
