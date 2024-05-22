/*
 * Flow Emulator
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/onflow/flow-go/fvm/storage/snapshot"
	flowgo "github.com/onflow/flow-go/model/flow"

	"github.com/onflow/flow-emulator/storage"
	"github.com/onflow/flow-emulator/types"
)

// InMemory implements the InMemory interface with an in-memory store.
// This is a copy of the implementation found in flow-emulator,
// with added functionality for serializing it to JSON.
// See: https://github.com/onflow/flow-emulator/blob/988cd96886514245ca3fbefaf6b4e621147bd72c/storage/memstore/memstore.go
type InMemory struct {
	mu sync.RWMutex
	// block ID to block height
	blockIDToHeight map[flowgo.Identifier]uint64
	// blocks by height
	blocks map[uint64]flowgo.Block
	// collections by ID
	collections map[flowgo.Identifier]flowgo.LightCollection
	// transactions by ID
	transactions map[flowgo.Identifier]flowgo.TransactionBody
	// Transaction results by ID
	transactionResults map[flowgo.Identifier]types.StorableTransactionResult
	// Ledger states by block height
	ledger map[uint64]snapshot.SnapshotTree
	// events by block height
	eventsByBlockHeight map[uint64][]flowgo.Event
	// highest block height
	blockHeight uint64
}

type InMemoryJson struct {
	BlockIDToHeight     map[flowgo.Identifier]uint64                          `json:"blockIDToHeight"`
	Blocks              map[uint64]flowgo.Block                               `json:"blocks"`
	Collections         map[flowgo.Identifier]flowgo.LightCollection          `json:"collections"`
	Transactions        map[flowgo.Identifier]flowgo.TransactionBody          `json:"transactions"`
	TransactionResults  map[flowgo.Identifier]types.StorableTransactionResult `json:"transactionResults"`
	EventsByBlockHeight map[uint64][]flowgo.Event                             `json:"eventsByBlockHeight"`
	BlockHeight         uint64                                                `json:"blockHeight"`
}

// New returns a new in-memory InMemory implementation.
func New() *InMemory {
	return &InMemory{
		mu:                  sync.RWMutex{},
		blockIDToHeight:     make(map[flowgo.Identifier]uint64),
		blocks:              make(map[uint64]flowgo.Block),
		collections:         make(map[flowgo.Identifier]flowgo.LightCollection),
		transactions:        make(map[flowgo.Identifier]flowgo.TransactionBody),
		transactionResults:  make(map[flowgo.Identifier]types.StorableTransactionResult),
		ledger:              make(map[uint64]snapshot.SnapshotTree),
		eventsByBlockHeight: make(map[uint64][]flowgo.Event),
	}
}

var _ storage.Store = &InMemory{}

func (s *InMemory) Json() ([]byte, error) {
	inMemoryJson := InMemoryJson{
		BlockIDToHeight:     s.blockIDToHeight,
		Blocks:              s.blocks,
		Collections:         s.collections,
		Transactions:        s.transactions,
		TransactionResults:  s.transactionResults,
		EventsByBlockHeight: s.eventsByBlockHeight,
	}
	return json.Marshal(inMemoryJson)
}

func (s *InMemory) Start() error {
	return nil
}

func (s *InMemory) Stop() {
}

func (s *InMemory) LatestBlockHeight(ctx context.Context) (uint64, error) {
	b, err := s.LatestBlock(ctx)
	if err != nil {
		return 0, err
	}

	return b.Header.Height, nil
}

func (s *InMemory) LatestBlock(ctx context.Context) (flowgo.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	latestBlock, ok := s.blocks[s.blockHeight]
	if !ok {
		return flowgo.Block{}, storage.ErrNotFound
	}
	return latestBlock, nil
}

func (s *InMemory) StoreBlock(ctx context.Context, block *flowgo.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.storeBlock(block)
}

func (s *InMemory) storeBlock(block *flowgo.Block) error {
	s.blocks[block.Header.Height] = *block
	s.blockIDToHeight[block.ID()] = block.Header.Height

	if block.Header.Height > s.blockHeight {
		s.blockHeight = block.Header.Height
	}

	return nil
}

func (s *InMemory) BlockByID(ctx context.Context, blockID flowgo.Identifier) (*flowgo.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockHeight, ok := s.blockIDToHeight[blockID]
	if !ok {
		return nil, storage.ErrNotFound
	}

	block, ok := s.blocks[blockHeight]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return &block, nil

}

func (s *InMemory) BlockByHeight(ctx context.Context, height uint64) (*flowgo.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	block, ok := s.blocks[height]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return &block, nil
}

func (s *InMemory) CommitBlock(
	ctx context.Context,
	block flowgo.Block,
	collections []*flowgo.LightCollection,
	transactions map[flowgo.Identifier]*flowgo.TransactionBody,
	transactionResults map[flowgo.Identifier]*types.StorableTransactionResult,
	executionSnapshot *snapshot.ExecutionSnapshot,
	events []flowgo.Event,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(transactions) != len(transactionResults) {
		return fmt.Errorf(
			"transactions count (%d) does not match result count (%d)",
			len(transactions),
			len(transactionResults),
		)
	}

	err := s.storeBlock(&block)
	if err != nil {
		return err
	}

	for _, col := range collections {
		err := s.insertCollection(*col)
		if err != nil {
			return err
		}
	}

	for _, tx := range transactions {
		err := s.insertTransaction(tx.ID(), *tx)
		if err != nil {
			return err
		}
	}

	for txID, result := range transactionResults {
		err := s.insertTransactionResult(txID, *result)
		if err != nil {
			return err
		}
	}

	err = s.insertExecutionSnapshot(
		block.Header.Height,
		executionSnapshot)
	if err != nil {
		return err
	}

	err = s.insertEvents(block.Header.Height, events)
	if err != nil {
		return err
	}

	return nil
}

func (s *InMemory) CollectionByID(
	ctx context.Context,
	collectionID flowgo.Identifier,
) (flowgo.LightCollection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.collections[collectionID]
	if !ok {
		return flowgo.LightCollection{}, storage.ErrNotFound
	}
	return tx, nil
}

func (s *InMemory) TransactionByID(
	ctx context.Context,
	transactionID flowgo.Identifier,
) (flowgo.TransactionBody, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.transactions[transactionID]
	if !ok {
		return flowgo.TransactionBody{}, storage.ErrNotFound
	}
	return tx, nil

}

func (s *InMemory) TransactionResultByID(
	ctx context.Context,
	transactionID flowgo.Identifier,
) (types.StorableTransactionResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.transactionResults[transactionID]
	if !ok {
		return types.StorableTransactionResult{}, storage.ErrNotFound
	}
	return result, nil

}

func (s *InMemory) LedgerByHeight(
	ctx context.Context,
	blockHeight uint64,
) (snapshot.StorageSnapshot, error) {
	return s.ledger[blockHeight], nil
}

func (s *InMemory) EventsByHeight(
	ctx context.Context,
	blockHeight uint64,
	eventType string,
) ([]flowgo.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	allEvents := s.eventsByBlockHeight[blockHeight]

	events := make([]flowgo.Event, 0)

	for _, event := range allEvents {
		if eventType == "" {
			events = append(events, event)
		} else {
			if string(event.Type) == eventType {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

func (s *InMemory) insertCollection(col flowgo.LightCollection) error {
	s.collections[col.ID()] = col
	return nil
}

func (s *InMemory) insertTransaction(txID flowgo.Identifier, tx flowgo.TransactionBody) error {
	s.transactions[txID] = tx
	return nil
}

func (s *InMemory) insertTransactionResult(txID flowgo.Identifier, result types.StorableTransactionResult) error {
	s.transactionResults[txID] = result
	return nil
}

func (s *InMemory) insertExecutionSnapshot(
	blockHeight uint64,
	executionSnapshot *snapshot.ExecutionSnapshot,
) error {
	oldLedger := s.ledger[blockHeight-1]

	s.ledger[blockHeight] = oldLedger.Append(executionSnapshot)

	return nil
}

func (s *InMemory) insertEvents(blockHeight uint64, events []flowgo.Event) error {
	if s.eventsByBlockHeight[blockHeight] == nil {
		s.eventsByBlockHeight[blockHeight] = events
	} else {
		s.eventsByBlockHeight[blockHeight] = append(s.eventsByBlockHeight[blockHeight], events...)
	}

	return nil
}
