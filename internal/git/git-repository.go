package git

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/rs/zerolog"
	"os"
	"path"
)

type Repository struct {
	logger     *zerolog.Logger
	repository *git.Repository
	storage    *memory.Storage
	fs         billy.Filesystem
}

type RepositoryFile struct {
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
}

func New(logger *zerolog.Logger) *Repository {
	return &Repository{
		logger: logger,
	}
}

func (r *Repository) Clone(url string) error {
	fs := memfs.New()
	storage := memory.NewStorage()
	repository, err := git.Clone(storage, fs, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		return err
	}

	r.repository = repository
	r.fs = fs
	r.storage = storage

	return nil
}

func (r *Repository) Stat(path string) (os.FileInfo, error) {
	return r.fs.Stat(path)
}

func (r *Repository) MkdirAll(path string, perm os.FileMode) error {
	panic("Unimplemented")
}

func (r *Repository) WriteFile(filename string, data []byte, perm os.FileMode) error {
	panic("Unimplemented")
}

func (r *Repository) ReadFile(path string) ([]byte, error) {
	f, err := r.fs.OpenFile(path, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer func(f billy.File) {
		err := f.Close()
		if err != nil {
			r.logger.Error().Msg(fmt.Sprintf("Error closing file %s: %s\n", path, err))
		}
	}(f)

	stat, err := r.fs.Stat(path)
	if err != nil {
		return nil, err
	}
	content := make([]byte, stat.Size())

	_, err = f.Read(content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *Repository) Files() ([]RepositoryFile, error) {
	return r.recursiveFiles(r.fs.Root())
}

func (r *Repository) recursiveFiles(dir string) ([]RepositoryFile, error) {
	files, err := r.fs.ReadDir(dir)
	repoFiles := make([]RepositoryFile, 0, len(files))

	for _, file := range files {
		fullPath := path.Join(dir, file.Name())
		if file.Mode().IsDir() {
			nestedFiles, nestedErr := r.recursiveFiles(fullPath)

			if nestedErr != nil {
				return nil, nestedErr
			}

			repoFiles = append(repoFiles, RepositoryFile{
				Path:        fullPath,
				IsDirectory: true,
			})

			repoFiles = append(repoFiles, nestedFiles...)
		} else {
			repoFiles = append(repoFiles, RepositoryFile{
				Path:        fullPath,
				IsDirectory: false,
			})
		}
	}

	return repoFiles, err
}

func (r *Repository) Commit() {
	panic("Unimplemented")
}

func (r *Repository) Push() {
	panic("Unimplemented")
}
