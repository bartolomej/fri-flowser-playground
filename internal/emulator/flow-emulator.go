package emulator

import (
	"fri-flowser-playground/internal/emulator/store"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/gateway"
	"github.com/rs/zerolog"
)

type Blockchain struct {
	logger     *zerolog.Logger
	store      *store.InMemory
	blockchain *emulator.Blockchain
	flow       *flowkit.Flowkit
	gateway    *gateway.EmulatorGateway
}

func New(logger *zerolog.Logger) *Blockchain {
	s := store.New()
	return &Blockchain{
		store:  s,
		logger: logger,
		gateway: gateway.NewEmulatorGatewayWithOpts(
			&gateway.EmulatorKey{
				PublicKey: emulator.DefaultServiceKey().AccountKey().PublicKey,
				SigAlgo:   emulator.DefaultServiceKeySigAlgo,
				HashAlgo:  emulator.DefaultServiceKeyHashAlgo,
			},
			gateway.WithEmulatorOptions(
				emulator.WithLogger(*logger),
				emulator.WithStore(s),
				emulator.WithTransactionValidationEnabled(false),
				emulator.WithStorageLimitEnabled(false),
				emulator.WithTransactionFeesEnabled(false),
			),
		),
	}
}

func (b *Blockchain) State() ([]byte, error) {
	return b.store.Json()
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
