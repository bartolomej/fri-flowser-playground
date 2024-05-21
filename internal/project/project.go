package project

import (
	"fmt"
	"fri-flowser-playground/internal/emulator"
	"fri-flowser-playground/internal/git"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/output"
	"github.com/rs/zerolog"
	"path"
	"strings"
)

type Project struct {
	blockchain *emulator.Blockchain
	repository *git.Repository
	logger     *zerolog.Logger
	kit        *flowkit.Flowkit
}

func New(logger *zerolog.Logger) *Project {
	repository := git.New(logger)
	blockchain := emulator.New(logger)

	return &Project{
		logger:     logger,
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

	kit, err := p.initFlowKit()

	if err != nil {
		return err
	}

	p.kit = kit

	contracts, err := p.cadenceContractFiles(files)
	if err != nil {
		return err
	}

	err = p.blockchain.Deploy(contracts)

	if err != nil {
		return err
	}

	p.logger.Printf("Deployed %d contracts\n", len(contracts))

	return nil
}

func (p *Project) initFlowKit() (*flowkit.Flowkit, error) {
	configFilePaths := []string{
		"flow.json",
	}
	state, err := flowkit.Load(configFilePaths, p.repository)
	if err != nil {
		return nil, err
	}

	network, err := state.Networks().ByName("emulator")
	if err != nil {
		return nil, err
	}

	flowKitLogger := newFlowKitLogger(p.logger)

	return flowkit.NewFlowkit(state, *network, p.blockchain.Gateway(), flowKitLogger), nil
}

func (p *Project) cadenceContractFiles(files []git.RepositoryFile) ([]emulator.ContractDescriptor, error) {
	contracts := make([]emulator.ContractDescriptor, 0)

	for _, file := range files {
		if path.Ext(file.Path) == ".cdc" && strings.Contains(file.Path, "/contracts/") {
			source, err := p.repository.ReadFile(file.Path)

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

type FlowKitLogger struct {
	logger *zerolog.Logger
}

var _ output.Logger = (*FlowKitLogger)(nil)

func newFlowKitLogger(logger *zerolog.Logger) *FlowKitLogger {
	return &FlowKitLogger{
		logger: logger,
	}
}

func (l *FlowKitLogger) Debug(s string) {
	l.logger.Debug().Msg(s)
}

func (l *FlowKitLogger) Info(s string) {
	l.logger.Info().Msg(s)
}

func (l *FlowKitLogger) Error(s string) {
	l.logger.Error().Msg(s)
}

func (l *FlowKitLogger) StartProgress(s string) {
	l.logger.Info().Msg(s)
}

func (l *FlowKitLogger) StopProgress() {
	// We don't support progress indication, so no need to do anything here
}
