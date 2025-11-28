// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package slimarchive

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ava-labs/libevm/common"
	"github.com/cockroachdb/pebble"
)

var (
	// Key prefixes - we only store traces since blocks/receipts survive pruning
	prefixCallTrace = []byte("ct:")   // ct:{txHash} -> callTracer result JSON
	prefixStateDiff = []byte("sd:")   // sd:{txHash} -> prestateTracer diffMode result JSON
	prefixMeta      = []byte("meta:") // meta:head -> last indexed block

	// Key for marking that slim archive has been enabled (cannot be disabled after)
	keyEnabled = append([]byte(nil), append(prefixMeta, []byte("enabled")...)...)

	ErrNotFound           = errors.New("not found in slim archive")
	ErrAlreadyInitialized = errors.New("slim archive was previously enabled - cannot disable without data loss")
)

// Store wraps PebbleDB for slim archive storage with ZSTD compression.
// Only stores traces - blocks and receipts are already in the chain DB and don't get pruned.
type Store struct {
	db   *pebble.DB
	path string
}

// NewStore opens or creates a new slim archive store at the given path.
// Uses 1MB block size and ZSTD level 1 compression for optimal space efficiency.
func NewStore(path string) (*Store, error) {
	opts := &pebble.Options{
		// Use ZSTD compression for all levels
		Levels: []pebble.LevelOptions{
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L0: 1MB blocks
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L1
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L2
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L3
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L4
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L5
			{Compression: pebble.ZstdCompression, BlockSize: 1 << 20}, // L6
		},
	}

	db, err := pebble.Open(path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open slim archive db: %w", err)
	}

	return &Store{
		db:   db,
		path: path,
	}, nil
}

// Close closes the store.
func (s *Store) Close() error {
	return s.db.Close()
}

// IsEnabled checks if the slim archive has been marked as enabled.
func (s *Store) IsEnabled() bool {
	_, closer, err := s.db.Get(keyEnabled)
	if err != nil {
		return false
	}
	closer.Close()
	return true
}

// MarkEnabled marks the slim archive as enabled.
// Once set, this marker persists and prevents the archive from being disabled.
func (s *Store) MarkEnabled() error {
	return s.db.Set(keyEnabled, []byte{1}, pebble.Sync)
}

// WasEverEnabled checks if slim archive was ever enabled at the given path.
// This is used to prevent enabling after it was previously disabled (data loss).
// Returns (wasEnabled, exists, error).
func WasEverEnabled(path string) (bool, bool, error) {
	opts := &pebble.Options{
		ReadOnly: true,
	}
	db, err := pebble.Open(path, opts)
	if err != nil {
		// If DB doesn't exist, that's fine - it was never enabled
		return false, false, nil
	}
	defer db.Close()

	_, closer, err := db.Get(keyEnabled)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return false, true, nil
		}
		return false, true, err
	}
	closer.Close()
	return true, true, nil
}

// prefixedKey creates a key with the given prefix and tx hash.
func prefixedKey(prefix []byte, hash common.Hash) []byte {
	key := make([]byte, len(prefix)+32)
	copy(key, prefix)
	copy(key[len(prefix):], hash[:])
	return key
}

// PutCallTrace stores callTracer result by tx hash.
func (s *Store) PutCallTrace(txHash common.Hash, data []byte) error {
	return s.db.Set(prefixedKey(prefixCallTrace, txHash), data, pebble.Sync)
}

// GetCallTrace retrieves callTracer result by tx hash.
func (s *Store) GetCallTrace(txHash common.Hash) ([]byte, error) {
	data, closer, err := s.db.Get(prefixedKey(prefixCallTrace, txHash))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer closer.Close()

	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

// HasCallTrace checks if a call trace exists for the given transaction hash.
func (s *Store) HasCallTrace(txHash common.Hash) bool {
	_, closer, err := s.db.Get(prefixedKey(prefixCallTrace, txHash))
	if err != nil {
		return false
	}
	closer.Close()
	return true
}

// PutStateDiff stores prestateTracer diffMode result by tx hash.
func (s *Store) PutStateDiff(txHash common.Hash, data []byte) error {
	return s.db.Set(prefixedKey(prefixStateDiff, txHash), data, pebble.Sync)
}

// GetStateDiff retrieves prestateTracer diffMode result by tx hash.
func (s *Store) GetStateDiff(txHash common.Hash) ([]byte, error) {
	data, closer, err := s.db.Get(prefixedKey(prefixStateDiff, txHash))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer closer.Close()

	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

// HasStateDiff checks if a state diff exists for the given transaction hash.
func (s *Store) HasStateDiff(txHash common.Hash) bool {
	_, closer, err := s.db.Get(prefixedKey(prefixStateDiff, txHash))
	if err != nil {
		return false
	}
	closer.Close()
	return true
}

// HasTrace checks if any trace (call or state diff) exists for the given transaction hash.
func (s *Store) HasTrace(txHash common.Hash) bool {
	return s.HasCallTrace(txHash) || s.HasStateDiff(txHash)
}

// GetTrace retrieves the appropriate trace based on tracer type.
// For "callTracer" or empty, returns call trace. For "prestateTracer", returns state diff.
func (s *Store) GetTrace(txHash common.Hash, tracer string) ([]byte, error) {
	if tracer == "prestateTracer" {
		return s.GetStateDiff(txHash)
	}
	return s.GetCallTrace(txHash)
}

// PutHead stores the last indexed block number.
func (s *Store) PutHead(num uint64) error {
	numBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numBytes, num)
	return s.db.Set(append(prefixMeta, []byte("head")...), numBytes, pebble.Sync)
}

// GetHead retrieves the last indexed block number.
func (s *Store) GetHead() (uint64, error) {
	data, closer, err := s.db.Get(append(prefixMeta, []byte("head")...))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	defer closer.Close()

	return binary.BigEndian.Uint64(data), nil
}

// WriteBatch allows writing multiple traces atomically.
type WriteBatch struct {
	batch *pebble.Batch
	store *Store
}

// NewBatch creates a new write batch.
func (s *Store) NewBatch() *WriteBatch {
	return &WriteBatch{
		batch: s.db.NewBatch(),
		store: s,
	}
}

// PutCallTrace adds a call trace to the batch.
func (b *WriteBatch) PutCallTrace(txHash common.Hash, data []byte) error {
	return b.batch.Set(prefixedKey(prefixCallTrace, txHash), data, nil)
}

// PutStateDiff adds a state diff to the batch.
func (b *WriteBatch) PutStateDiff(txHash common.Hash, data []byte) error {
	return b.batch.Set(prefixedKey(prefixStateDiff, txHash), data, nil)
}

// PutHead adds head update to the batch.
func (b *WriteBatch) PutHead(num uint64) error {
	numBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numBytes, num)
	return b.batch.Set(append(prefixMeta, []byte("head")...), numBytes, nil)
}

// Commit writes the batch to the database.
func (b *WriteBatch) Commit() error {
	return b.batch.Commit(pebble.Sync)
}

// Close discards the batch.
func (b *WriteBatch) Close() error {
	return b.batch.Close()
}
