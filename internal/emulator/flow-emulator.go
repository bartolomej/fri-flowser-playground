package emulator

import (
	"fmt"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/rs/zerolog"
	"os"
)

type Blockchain struct {
	logger     *zerolog.Logger
	blockchain *emulator.Blockchain
}

func (b *Blockchain) Start() error {
	logger := initLogger()
	blockchain, err := emulator.New(
		emulator.WithLogger(*logger),
	)

	if err != nil {
		return err
	}

	b.blockchain = blockchain
	b.logger = logger

	return nil
}

type ContractDescriptor struct {
	Source []byte
}

func (b *Blockchain) Deploy(descriptors []ContractDescriptor) error {
	contracts := make([]emulator.ContractDescription, 0)

	serviceKey := b.blockchain.ServiceKey()
	serviceAddress := serviceKey.Address

	for _, descriptor := range descriptors {
		contracts = append(contracts, emulator.ContractDescription{
			Name:        "Example",
			Description: "Deploying example contract",
			Address:     serviceAddress,
			Source:      descriptor.Source,
		})
	}

	return emulator.DeployContracts(b.blockchain, contracts)
}

func initLogger() *zerolog.Logger {

	level := zerolog.InfoLevel
	zerolog.MessageFieldName = "msg"

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &logger
}

func NewTextWriter() zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("%-44s", i)
	}

	return writer
}
