package dir

import (
	"fmt"
	"os"

	"github.com/google/renameio/v2"
)

type File struct {
	Payload  []byte
	Name     string
	FullPath string
}

type Files struct {
	tmpDir  string
	baseDir string
}

func NewFiles(tmpDir, baseDir string) *Files {
	return &Files{
		tmpDir:  tmpDir,
		baseDir: baseDir,
	}
}

func (f *Files) NewFile(name string, payload []byte) File {
	return File{
		Payload:  payload,
		Name:     name,
		FullPath: fullPath(f.baseDir, name),
	}
}

func fullPath(baseDir, name string) string {
	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, name)
}

func (f *Files) AddFile(file File) error {
	return renameio.WriteFile(
		file.FullPath, file.Payload, os.ModePerm, renameio.WithTempDir(f.tmpDir),
	)
}

func (f *Files) readFile(name string) (File, error) {
	payload, err := os.ReadFile(fullPath(f.baseDir, name))
	if err != nil {
		return File{}, fmt.Errorf("os.ReadFile: %w", err)
	}
	return f.NewFile(name, payload), nil
}

func (f *Files) ListFiles() ([]File, error) {
	dirEntities, err := os.ReadDir(f.baseDir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}

	out := make([]File, 0, len(dirEntities))
	for _, dir := range dirEntities {
		file, err := f.readFile(dir.Name())
		if err != nil {
			return nil, fmt.Errorf("readFile %s: %w", dir.Name(), err)
		}

		out = append(out, file)
	}

	return out, nil
}

func (f *Files) DeleteByName(name string) error {
	return os.RemoveAll(fullPath(f.baseDir, name))
}

func (f *Files) GetByName(name string) (File, error) {
	return f.readFile(name)
}
