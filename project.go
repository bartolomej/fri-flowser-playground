package main

import "fmt"

type Project struct {
	emulator   *FlowEmulator
	repository *GitRepository
}

func New() *Project {
	repository := &GitRepository{}
	emulator := &FlowEmulator{}

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
