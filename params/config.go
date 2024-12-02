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
// Copyright 2016 The go-ethereum Authors
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

package params

import (
	"math/big"

	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/coreth/utils"
	gethparams "github.com/ethereum/go-ethereum/params"
)

// Guarantees extras initialisation before a call to [ChainConfig.Rules].
var _ = libevmInit()

var (
	// SubnetEVMDefaultConfig is the default configuration
	// without any network upgrades.
	SubnetEVMDefaultChainConfig = WithExtra(
		&ChainConfig{
			ChainID: DefaultChainID,

			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		},
		&ChainConfigExtra{
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades:    getDefaultNetworkUpgrades(upgrade.GetConfig(constants.MainnetID)), // This can be changed to correct network (local, test) via VM.
			GenesisPrecompiles: Precompiles{},
		},
	)

	TestChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			ShanghaiTime:        utils.TimeToNewUint64(upgrade.GetConfig(constants.UnitTestID).DurangoTime),
			CancunTime:          utils.TimeToNewUint64(upgrade.GetConfig(constants.UnitTestID).EtnaTime),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades:    getDefaultNetworkUpgrades(upgrade.GetConfig(constants.UnitTestID)), // This can be changed to correct network (local, test) via VM.
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestPreSubnetEVMChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestSubnetEVMChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestDurangoChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			ShanghaiTime:        utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.UnscheduledActivationTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)

	TestEtnaChainConfig = WithExtra(
		&ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0),
			ShanghaiTime:        utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
			CancunTime:          utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
		},
		&ChainConfigExtra{
			AvalancheContext:   AvalancheContext{utils.TestSnowContext()},
			FeeConfig:          DefaultFeeConfig,
			AllowFeeRecipients: false,
			NetworkUpgrades: NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
				EtnaTimestamp:      utils.TimeToNewUint64(upgrade.InitiallyActiveTime),
			},
			GenesisPrecompiles: Precompiles{},
			UpgradeConfig:      UpgradeConfig{},
		},
	)
	TestRules = TestChainConfig.Rules(new(big.Int), IsMergeTODO, 0)
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig = gethparams.ChainConfig

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules = gethparams.Rules
