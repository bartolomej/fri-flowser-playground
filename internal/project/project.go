package project

import (
	"context"
	"encoding/json"
	"fmt"
	"fri-flowser-playground/internal/emulator"
	"fri-flowser-playground/internal/git"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/accounts"
	"github.com/onflow/flowkit/arguments"
	"github.com/onflow/flowkit/output"
	"github.com/onflow/flowkit/transactions"
	"github.com/rs/zerolog"
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

func (p *Project) Files() ([]git.RepositoryFile, error) {
	return p.repository.Files()
}

func (p *Project) Open(projectUrl string) error {
	fmt.Printf("Cloning project: %s\n", projectUrl)

	err := p.repository.Clone(projectUrl)

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

	err = p.setupAccounts()

	if err != nil {
		return err
	}

	contracts, err := p.kit.DeployProject(context.Background(), flowkit.UpdateExistingContract(true))

	if err != nil {
		return err
	}

	p.logger.Info().Msg(fmt.Sprintf("Deployed %d contracts\n", len(contracts)))

	return nil
}

func (p *Project) ExecuteScript(code []byte, location string, argsJson string) (res []byte, err error) {
	var args []cadence.Value
	if argsJson != "" {
		args, err = arguments.ParseJSON(argsJson)
	}
	if err != nil {
		return nil, err
	}

	result, err := p.kit.ExecuteScript(
		context.Background(),
		flowkit.Script{Code: code, Args: args, Location: location},
		flowkit.LatestScriptQuery,
	)

	if err != nil {
		return nil, err
	}

	return []byte(result.String()), err
}

func (p *Project) ExecuteTransaction(code []byte, location string, argsJson string) ([]byte, error) {
	state, err := p.kit.State()
	if err != nil {
		return nil, err
	}
	serviceAccount, err := state.EmulatorServiceAccount()

	var args []cadence.Value
	if argsJson != "" {
		args, err = arguments.ParseJSON(argsJson)
	}
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(1000)

	_, result, err := p.kit.SendTransaction(
		context.Background(),
		transactions.AccountRoles{
			Proposer:    *serviceAccount,
			Authorizers: []accounts.Account{},
			Payer:       *serviceAccount,
		},
		flowkit.Script{Code: code, Args: args, Location: location},
		gasLimit,
	)

	if err != nil {
		return nil, err
	}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return jsonResult, err
}

// setupAccounts creates account on the network and updates the state
// Uses the same approach as in: https://github.com/onflow/flow-cli/blob/f1bcd08d61bf1f20a41b1005158662d094004c65/internal/super/project.go#L207
func (p *Project) setupAccounts() error {
	state, err := p.kit.State()
	if err != nil {
		return err
	}

	for _, confAccount := range state.Config().Accounts {
		// Must be run each loop for now, as we are modifying elements of the accounts array
		// with the Accounts().Remove() call, which is a temporary solution for a big in flow-kit.
		serviceAccount, err := state.EmulatorServiceAccount()
		if err != nil {
			return err
		}
		privateKey, err := serviceAccount.Key.PrivateKey()
		if err != nil {
			return err
		}
		pubKey := (*privateKey).PublicKey()

		p.logger.Info().Msg(fmt.Sprintf("Creating account %s %s", confAccount.Name, confAccount.Address))

		existingAccount, _ := p.kit.Gateway().GetAccount(context.Background(), confAccount.Address)

		// Only create non-existing accounts
		if existingAccount != nil {
			continue
		}

		created, _, err := p.kit.CreateAccount(
			context.Background(),
			serviceAccount,
			[]accounts.PublicKey{{
				Public:   pubKey,
				Weight:   flow.AccountKeyWeightThreshold,
				SigAlgo:  crypto.ECDSA_P256,
				HashAlgo: crypto.SHA3_256,
			}},
		)

		if err != nil {
			panic(err)
		}

		// There is a bug that prevents `AddOrUpdate` from updating an existing record, so we must remove it first.
		// See: https://github.com/onflow/flowkit/blame/2f09f4a76225d658c31147edc419695efb241e25/accounts/account.go#L188
		_ = state.Accounts().Remove(confAccount.Name)

		state.Accounts().AddOrUpdate(&accounts.Account{
			Name:    confAccount.Name,
			Address: created.Address,
			Key:     accounts.NewHexKeyFromPrivateKey(0, crypto.SHA3_256, *privateKey),
		})

		p.logger.Info().Msg(fmt.Sprintf("Created account %s", created.Address))
	}

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
