package project

import (
	"fmt"
	"fri-flowser-playground/internal/emulator"
	"fri-flowser-playground/internal/git"
	"path"
	"strings"
)

type Project struct {
	blockchain *emulator.Blockchain
	repository *git.Repository
}

func New() *Project {
	repository := &git.Repository{}
	blockchain := &emulator.Blockchain{}

	return &Project{
		repository: repository,
		blockchain: blockchain,
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

	err = p.blockchain.Start()

	if err != nil {
		return err
	}

	contracts, err := p.cadenceContractFiles(files)
	if err != nil {
		return err
	}

	err = p.blockchain.Deploy(contracts)

	if err != nil {
		return err
	}

	return nil
}

func (p *Project) cadenceContractFiles(files []git.RepositoryFile) ([]emulator.ContractDescriptor, error) {
	contracts := make([]emulator.ContractDescriptor, 0)

	for _, file := range files {
		if path.Ext(file.Path) == ".cdc" && strings.Contains(file.Path, "/contracts/") {
			source, err := p.repository.FileContent(file.Path)

			if err != nil {
				return nil, err
			}

			contracts = append(contracts, emulator.ContractDescriptor{
				Source: source,
			})
		}
	}

	return contracts, nil
}
