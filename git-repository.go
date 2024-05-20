package main

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"path"
)

type GitRepository struct {
	repository *git.Repository
	storage    *memory.Storage
	fs         billy.Filesystem
}

type GitRepositoryFile struct {
	Path        string
	IsDirectory bool
}

func (r *GitRepository) Clone(url string) error {
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

func (r *GitRepository) Files() ([]GitRepositoryFile, error) {
	return r.recursiveFiles(r.fs.Root())
}

func (r *GitRepository) recursiveFiles(dir string) ([]GitRepositoryFile, error) {
	files, err := r.fs.ReadDir(dir)
	repoFiles := make([]GitRepositoryFile, 0, len(files))

	for _, file := range files {
		fullPath := path.Join(dir, file.Name())
		if file.Mode().IsDir() {
			nestedFiles, nestedErr := r.recursiveFiles(fullPath)

			if nestedErr != nil {
				return nil, nestedErr
			}

			repoFiles = append(repoFiles, GitRepositoryFile{
				Path:        fullPath,
				IsDirectory: true,
			})

			repoFiles = append(repoFiles, nestedFiles...)
		} else {
			repoFiles = append(repoFiles, GitRepositoryFile{
				Path:        fullPath,
				IsDirectory: false,
			})
		}
	}

	return repoFiles, err
}

func (r *GitRepository) Commit() {
	fmt.Println("Commiting")
}

func (r *GitRepository) Push() {
	fmt.Println("Pushing")
}
