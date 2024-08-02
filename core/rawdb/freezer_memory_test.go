// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2024 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/rawdb/ancienttest"
	"github.com/ethereum/go-ethereum/ethdb"
)

func TestMemoryFreezer(t *testing.T) {
	ancienttest.TestAncientSuite(t, func(kinds []string) ethdb.AncientStore {
		tables := make(map[string]bool)
		for _, kind := range kinds {
			tables[kind] = true
		}
		return NewMemoryFreezer(false, tables)
	})
	ancienttest.TestResettableAncientSuite(t, func(kinds []string) ethdb.ResettableAncientStore {
		tables := make(map[string]bool)
		for _, kind := range kinds {
			tables[kind] = true
		}
		return NewMemoryFreezer(false, tables)
	})
}
