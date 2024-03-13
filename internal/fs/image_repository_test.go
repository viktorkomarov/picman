package fs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/viktorkomarov/picman/internal/fs"
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
	i.Suite.NoError(err, "failed to create temp dir")
	i.testDir = testDir

	repo, err := fs.NewImageRepository(os.TempDir(), testDir)
	i.Suite.NoError(err, "failed to create image repository")
	i.repo = repo
}

func (i *ImageRepositoryTestSuite) TearDownSuite() {
	i.Suite.NoError(os.RemoveAll(i.testDir), "failed to create temp dir")
}

func TestImageRepositiry(t *testing.T) {
	suite.Run(t, new(ImageRepositoryTestSuite))
}

func (s *ImageRepositoryTestSuite) SaveImage() {

}
