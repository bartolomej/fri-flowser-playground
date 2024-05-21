package project

import (
	"fmt"
	emulator2 "fri-flowser-playground/internal/emulator"
	"fri-flowser-playground/internal/git"
)

type Project struct {
	emulator   *emulator2.FlowEmulator
	repository *git.Repository
}

func New() *Project {
	repository := &git.Repository{}
	emulator := &emulator2.FlowEmulator{}

	return &Project{
		repository: repository,
		emulator:   emulator,
	}
}

func (p *Project) Open(projectUrl string) error {
	fmt.Printf("Cloning project: %s\n", projectUrl)

	err := p.repository.Clone(projectUrl)

	if err != nil {
		return err
	}

	files, err := p.repository.Files()

	fmt.Printf("Cloned %d files\n", len(files))

	if err != nil {
		return err
	}

	err = p.emulator.Start()

	if err != nil {
		return err
	}

	return nil
}
