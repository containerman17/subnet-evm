// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package slimarchive

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/hexutil"
	"github.com/ava-labs/subnet-evm/rpc"
)

// ArchiveAPI provides RPC methods for querying slim archive trace data.
// Only serves traces - blocks and receipts come from the normal chain DB.
type ArchiveAPI struct {
	store *Store
}

// NewArchiveAPI creates a new archive API instance.
func NewArchiveAPI(store *Store) *ArchiveAPI {
	return &ArchiveAPI{store: store}
}

// TraceTransaction returns the stored trace for a transaction.
// Supports callTracer (default) and prestateTracer (with diffMode).
func (api *ArchiveAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (json.RawMessage, error) {
	tracer := ""
	if config != nil && config.Tracer != nil {
		tracer = *config.Tracer
	}

	data, err := api.store.GetTrace(hash, tracer)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, errors.New("trace not found in slim archive (tx may be pending or node still syncing)")
		}
		return nil, err
	}
	return json.RawMessage(data), nil
}

// HasTrace checks if a trace exists for the given transaction.
func (api *ArchiveAPI) HasTrace(ctx context.Context, hash common.Hash) (bool, error) {
	return api.store.HasTrace(hash), nil
}

// GetHead returns the last indexed block number.
func (api *ArchiveAPI) GetHead(ctx context.Context) (hexutil.Uint64, error) {
	head, err := api.store.GetHead()
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return hexutil.Uint64(head), nil
}

// TraceConfig mirrors the tracer config for API compatibility.
type TraceConfig struct {
	Tracer       *string          `json:"tracer,omitempty"`
	TracerConfig *json.RawMessage `json:"tracerConfig,omitempty"`
	Timeout      *string          `json:"timeout,omitempty"`
}

// APIs returns the RPC API descriptors for the slim archive service.
func APIs(store *Store) []rpc.API {
	api := NewArchiveAPI(store)
	return []rpc.API{
		{
			Namespace: "slimarchive",
			Service:   api,
			Name:      "slim-archive",
		},
	}
}
