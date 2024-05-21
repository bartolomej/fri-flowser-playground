package emulator

import (
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-emulator/storage/memstore"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/gateway"
	"github.com/rs/zerolog"
)

type Blockchain struct {
	logger     *zerolog.Logger
	blockchain *emulator.Blockchain
	flow       *flowkit.Flowkit
	gateway    *gateway.EmulatorGateway
}

func New(logger *zerolog.Logger) *Blockchain {
	return &Blockchain{
		logger: logger,
		gateway: gateway.NewEmulatorGatewayWithOpts(
			&gateway.EmulatorKey{
				PublicKey: emulator.DefaultServiceKey().AccountKey().PublicKey,
				SigAlgo:   emulator.DefaultServiceKeySigAlgo,
				HashAlgo:  emulator.DefaultServiceKeyHashAlgo,
			},
			gateway.WithEmulatorOptions(
				emulator.WithLogger(*logger),
				emulator.WithStore(memstore.New()),
				emulator.WithTransactionValidationEnabled(false),
				emulator.WithStorageLimitEnabled(false),
				emulator.WithTransactionFeesEnabled(false),
			),
		),
	}
}

func (b *Blockchain) Gateway() gateway.Gateway {
	return b.gateway
}

func (b *Blockchain) Start() error {
	blockchain, err := emulator.New(
		emulator.WithLogger(*b.logger),
	)

	if err != nil {
		return err
	}

	b.blockchain = blockchain

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
