// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package slimarchive

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/state"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/libevm/eth/tracers"
	"github.com/ava-labs/libevm/log"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params"
)

// Indexer implements core.BlockHook to capture and store trace data.
// Only stores traces - blocks and receipts are already in chain DB and don't get pruned.
type Indexer struct {
	store       *Store
	chainConfig *params.ChainConfig
	chain       ChainContext
}

// ChainContext provides chain data needed for tracing.
type ChainContext interface {
	core.ChainContext
	GetHeader(common.Hash, uint64) *types.Header
}

// NewIndexer creates a new indexer with the given store.
func NewIndexer(store *Store, chainConfig *params.ChainConfig, chain ChainContext) *Indexer {
	return &Indexer{
		store:       store,
		chainConfig: chainConfig,
		chain:       chain,
	}
}

// OnBlockProcessed implements core.BlockHook.
// It captures trace data for each transaction after block execution.
func (idx *Indexer) OnBlockProcessed(block *types.Block, receipts types.Receipts, statedb *state.StateDB) error {
	blockNum := block.NumberU64()
	txs := block.Transactions()

	if len(txs) == 0 {
		// No transactions to trace, just update head
		if err := idx.store.PutHead(blockNum); err != nil {
			return fmt.Errorf("failed to update head: %w", err)
		}
		return nil
	}

	batch := idx.store.NewBatch()
	defer batch.Close()

	// Get parent for block context
	parent := idx.chain.GetHeader(block.ParentHash(), blockNum-1)
	if parent == nil && blockNum > 0 {
		return fmt.Errorf("parent header not found for block %d", blockNum)
	}

	signer := types.MakeSigner(idx.chainConfig, block.Number(), block.Time())
	blockContext := core.NewEVMBlockContext(block.Header(), idx.chain, nil)

	// Create a copy of statedb for tracing - we'll re-execute txs with tracer
	traceStateDB := statedb.Copy()

	for i, tx := range txs {
		txHash := tx.Hash()

		// Generate and store call trace using callTracer
		callTrace, err := idx.runTracer(block, tx, i, signer, blockContext, traceStateDB.Copy(), "callTracer", nil)
		if err != nil {
			log.Error("Failed to trace transaction with callTracer", "tx", txHash, "err", err)
			callTrace, _ = json.Marshal(map[string]string{"error": err.Error()})
		}
		if err := batch.PutCallTrace(txHash, callTrace); err != nil {
			log.Error("Failed to store call trace", "tx", txHash, "err", err)
		}

		// Generate and store state diff using prestateTracer with diffMode
		diffConfig := json.RawMessage(`{"diffMode": true}`)
		stateDiff, err := idx.runTracer(block, tx, i, signer, blockContext, traceStateDB.Copy(), "prestateTracer", diffConfig)
		if err != nil {
			log.Error("Failed to trace transaction with prestateTracer", "tx", txHash, "err", err)
			stateDiff, _ = json.Marshal(map[string]string{"error": err.Error()})
		}
		if err := batch.PutStateDiff(txHash, stateDiff); err != nil {
			log.Error("Failed to store state diff", "tx", txHash, "err", err)
		}
	}

	// Update head
	if err := batch.PutHead(blockNum); err != nil {
		return fmt.Errorf("failed to update head: %w", err)
	}

	if err := batch.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch for block %d: %w", blockNum, err)
	}

	log.Debug("Indexed block traces", "number", blockNum, "hash", block.Hash(), "txs", len(txs))
	return nil
}

// runTracer executes a transaction with the specified tracer and returns the result.
func (idx *Indexer) runTracer(block *types.Block, tx *types.Transaction, txIndex int, signer types.Signer, blockContext vm.BlockContext, statedb *state.StateDB, tracerName string, tracerConfig json.RawMessage) ([]byte, error) {
	if statedb == nil {
		return nil, fmt.Errorf("state not available for tracing")
	}

	// Create tracer
	tracer, err := tracers.DefaultDirectory.New(tracerName, nil, tracerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer %s: %w", tracerName, err)
	}

	// Prepare message
	msg, err := core.TransactionToMessage(tx, signer, block.BaseFee())
	if err != nil {
		return nil, fmt.Errorf("failed to convert tx to message: %w", err)
	}

	// Create EVM with tracer
	txContext := core.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, statedb, idx.chainConfig, vm.Config{
		Tracer:    tracer,
		NoBaseFee: true,
	})

	// Execute
	statedb.SetTxContext(tx.Hash(), txIndex)
	_, err = core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.GasLimit))
	if err != nil {
		// Transaction execution failed, but we still want the trace
		log.Debug("Transaction execution failed during tracing", "tx", tx.Hash(), "tracerName", tracerName, "err", err)
	}

	// Get trace result
	result, err := tracer.GetResult()
	if err != nil {
		return nil, fmt.Errorf("failed to get trace result: %w", err)
	}

	return json.Marshal(result)
}

// Store returns the underlying store.
func (idx *Indexer) Store() *Store {
	return idx.store
}

// Close closes the indexer and its store.
func (idx *Indexer) Close() error {
	return idx.store.Close()
}
