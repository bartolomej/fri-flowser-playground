package git

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"path"
)

type Repository struct {
	repository *git.Repository
	storage    *memory.Storage
	fs         billy.Filesystem
}

type RepositoryFile struct {
	Path        string
	IsDirectory bool
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
	fmt.Println("Commiting")
}

func (r *Repository) Push() {
	fmt.Println("Pushing")
}