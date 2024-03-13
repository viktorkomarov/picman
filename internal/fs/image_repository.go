package fs

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/viktorkomarov/picman/internal/domain"
)

type ImageRepository struct {
	tmpDir  string
	baseDir string
}

const rwxForPublicMask = fs.FileMode(007)

func validateDir(baseDir string) error {
	stat, err := os.Stat(baseDir)
	if err != nil {
		return fmt.Errorf("file.State: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s isn't a direction", baseDir)
	}

	if stat.Mode()&rwxForPublicMask != rwxForPublicMask {
		return fmt.Errorf("%s should have rights for others(rwx)", baseDir)
	}
	return nil
}

func NewImageRepository(tmpDir, baseDir string) (*ImageRepository, error) {
	if err := validateDir(baseDir); err != nil {
		return nil, fmt.Errorf("valildate base dir: %w", err)
	}
	if err := validateDir(tmpDir); err != nil {
		return nil, fmt.Errorf("validate tmp dir: %w", err)
	}

	return &ImageRepository{
		tmpDir:  tmpDir,
		baseDir: baseDir,
	}, nil
}

func (i *ImageRepository) normalizeFileName(name string) string {
	return fmt.Sprintf("%s/%s", i.baseDir, name)
}

func (i *ImageRepository) SaveImage(image domain.Image) error {
	tmpFile, err := os.CreateTemp(i.tmpDir, image.Name)
	if err != nil {
		return fmt.Errorf("os.CreateTemp: %w", err)
	}
	if err := os.WriteFile(tmpFile.Name(), image.Payload, os.ModePerm); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}
	// todo::read about write to file atomically
	return os.Rename(tmpFile.Name(), i.normalizeFileName(image.Name))
}

func (i *ImageRepository) DeleteByName(name string) error {
	return os.RemoveAll(i.normalizeFileName(name))
}

func (i *ImageRepository) readImage(name string) (domain.Image, error) {
	payload, err := os.ReadFile(i.normalizeFileName(name))
	if err != nil {
		return domain.Image{}, fmt.Errorf("os.Open: %w", err)
	}
	return domain.Image{
		Name:    name,
		Payload: payload,
	}, nil
}

func (i *ImageRepository) ListImages() ([]domain.Image, error) {
	dirs, err := os.ReadDir(i.baseDir)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	out := make([]domain.Image, 0)
	for _, dir := range dirs {
		image, err := i.readImage(dir.Name())
		if err != nil {
			return nil, fmt.Errorf("readImage %s: %w", dir.Name(), err)
		}

		out = append(out, image)
	}

	return out, err
}

func (i *ImageRepository) GetByName(name string) (domain.Image, error) {
	payload, err := os.ReadFile(i.normalizeFileName(name))
	if err != nil {
		return domain.Image{}, fmt.Errorf("os.ReadFile: %w", err)
	}
	return domain.Image{
		Name:    name,
		Payload: payload,
	}, nil
}
