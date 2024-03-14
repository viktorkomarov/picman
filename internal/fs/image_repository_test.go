package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/viktorkomarov/picman/internal/domain"
	"github.com/viktorkomarov/picman/internal/fs"
)

const (
	bananaImagePath = "./testdata/banana.svg"
	catImagePath    = "./testdata/cat.svg"

	bananaFileName = "banana.svg"
	catFileName    = "cat.svg"
)

type ImageRepositoryTestSuite struct {
	suite.Suite

	t       *testing.T
	testDir string
	repo    *fs.ImageRepository
}

func (i *ImageRepositoryTestSuite) SetupSuite() {
	i.t = i.Suite.T()

	testDir, err := os.MkdirTemp(os.TempDir(), "imageTest")
	i.NoError(err, "failed to create temp dir")
	i.testDir = testDir

	repo, err := fs.NewImageRepository(testDir, testDir)
	i.NoError(err, "failed to create image repository")
	i.repo = repo
}

func (i *ImageRepositoryTestSuite) TearDownSuite() {
	i.NoError(os.RemoveAll(i.testDir), "failed to remove temp dir")
}

func TestImageRepositiry(t *testing.T) {
	suite.Run(t, new(ImageRepositoryTestSuite))
}

func (s *ImageRepositoryTestSuite) TestSaveReadDeleteImages() {
	var (
		bananaImageTarget domain.Image
		catImageTarget    domain.Image
	)

	saveImage := func(origPath, name string) domain.Image {
		payload, err := os.ReadFile(origPath)
		s.NoErrorf(err, "failed to read orig file: %s", origPath)

		image := domain.Image{Name: name, Payload: payload}
		s.NoError(s.repo.SaveImage(image), "failed to save %s image", name)
		s.ErrorIs(s.repo.SaveImage(image), domain.ErrImageAlreadyExists, "second save should return error")
		return image
	}

	s.Run("save images", func() {
		bananaImageTarget = saveImage(bananaImagePath, bananaFileName)
		catImageTarget = saveImage(catImagePath, catFileName)
	})

	s.Run("must read two files", func() {
		bananImage, err := s.repo.GetByName(bananaFileName)
		s.NoError(err, "failed to get banana image")
		s.Equalf(bananaImageTarget, bananImage, "banana image is different")

		catImage, err := s.repo.GetByName(catFileName)
		s.NoError(err, "failed to get cat image")
		s.Equalf(catImageTarget, catImage, "cat image is different")

		images, err := s.repo.ListImages()
		s.NoError(err, "failed to read get images")
		s.Len(images, 2)
	})

	s.Run("delete images", func() {
		s.NoError(s.repo.DeleteByName(bananaFileName), "failed to delete banana image")
		s.NoError(s.repo.DeleteByName(catFileName), "failed to delete cat image")
		// second delete shouldn't return err
		s.NoError(s.repo.DeleteByName(bananaFileName), "failed to delete banana image")
		s.NoError(s.repo.DeleteByName(catFileName), "failed to delete cat image")
	})

	s.Run("dir must be empty", func() {
		_, err := s.repo.GetByName(bananaFileName)
		s.Error(err, "banana image mustn't exist")

		_, err = s.repo.GetByName(catFileName)
		s.Error(err, "cat image mustn't exist")

		images, err := s.repo.ListImages()
		s.NoError(err, "failed to read get images")
		s.Len(images, 0)
	})
}
