package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/samber/lo"
	"github.com/viktorkomarov/picman/internal/domain"
	"github.com/viktorkomarov/picman/internal/fs/dir"
	"github.com/viktorkomarov/picman/internal/utils/keylock"
	"github.com/wneessen/go-fileperm"
)

type ImageRepository struct {
	locker keylock.HashedMutex
	files  *dir.Files
}

func validateDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s must be a direction", dir)
	}

	perm, err := fileperm.New(dir)
	if err != nil {
		return fmt.Errorf("fileperm.New: %w", err)
	}
	if !perm.UserReadable() || !perm.UserWritable() {
		return fmt.Errorf("current user must have read and write rights")
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
		files:  dir.NewFiles(tmpDir, baseDir),
		locker: keylock.NewHashedMutex(100),
	}, nil
}

func (i *ImageRepository) withLock(key string, fn func() error) error {
	i.locker.Lock(key)
	defer i.locker.Unlock(key)
	return fn()
}

func (i *ImageRepository) SaveImage(image domain.Image) error {
	return i.withLock(image.Name, func() error {
		_, err := i.files.GetByName(image.Name)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return i.files.AddFile(i.files.NewFile(image.Name, image.Payload))
		case err == nil:
			return fmt.Errorf("%w: %s", domain.ErrImageAlreadyExists, image.Name)
		default:
			return fmt.Errorf("files.GetByName[%s]: %w", image.Name, err)
		}
	})
}

func (i *ImageRepository) DeleteByName(name string) error {
	return i.withLock(name, func() error {
		_, err := i.files.GetByName(name)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return nil
		case err == nil:
			return i.files.DeleteByName(name)
		default:
			return fmt.Errorf("files.GetByName: %s", name)
		}
	})
}

func (i *ImageRepository) ListImages() ([]domain.Image, error) {
	files, err := i.files.ListFiles()
	if err != nil {
		return nil, fmt.Errorf("ListFiles: %w", err)
	}

	return lo.Map(files, func(file dir.File, _ int) domain.Image {
		return fileToImage(file)
	}), nil
}

func (i *ImageRepository) GetByName(name string) (domain.Image, error) {
	file, err := i.files.GetByName(name)
	if err != nil {
		return domain.Image{}, fmt.Errorf("files.GetByName: %w", err)
	}
	return fileToImage(file), nil
}

func fileToImage(file dir.File) domain.Image {
	return domain.Image{
		Payload: file.Payload,
		Name:    file.Name,
	}
}
