// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package rewardmanager

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileconfig.Config = &Config{}

type InitialRewardConfig struct {
	AllowFeeRecipients bool           `json:"allowFeeRecipients"`
	RewardAddress      common.Address `json:"rewardAddress,omitempty"`
}

func (u *InitialRewardConfig) ToBytesWithPacker(p *wrappers.Packer) error {
	p.PackBool(u.AllowFeeRecipients)
	if p.Err != nil {
		return p.Err
	}
	p.PackBytes(u.RewardAddress[:])
	return p.Err
}

func (u *InitialRewardConfig) FromBytesWithPacker(p *wrappers.Packer) error {
	u.AllowFeeRecipients = p.UnpackBool()
	if p.Err != nil {
		return p.Err
	}
	u.RewardAddress = common.BytesToAddress(p.UnpackBytes())
	if p.Err != nil {
		return p.Err
	}
	return nil
}

func (i *InitialRewardConfig) Equal(other *InitialRewardConfig) bool {
	if other == nil {
		return false
	}

	return i.AllowFeeRecipients == other.AllowFeeRecipients && i.RewardAddress == other.RewardAddress
}

func (i *InitialRewardConfig) Verify() error {
	switch {
	case i.AllowFeeRecipients && i.RewardAddress != (common.Address{}):
		return ErrCannotEnableBothRewards
	default:
		return nil
	}
}

func (i *InitialRewardConfig) Configure(state contract.StateDB) error {
	// enable allow fee recipients
	if i.AllowFeeRecipients {
		EnableAllowFeeRecipients(state)
	} else if i.RewardAddress == (common.Address{}) {
		// if reward address is empty and allow fee recipients is false
		// then disable rewards
		DisableFeeRewards(state)
	} else {
		// set reward address
		return StoreRewardAddress(state, i.RewardAddress)
	}
	return nil
}

// Config implements the StatefulPrecompileConfig interface while adding in the
// RewardManager specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins], [enableds] and [managers] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address, initialConfig *InitialRewardConfig) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade:             precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (*Config) Key() string { return ConfigKey }

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	if c.InitialRewardConfig != nil {
		if err := c.InitialRewardConfig.Verify(); err != nil {
			return err
		}
	}
	return c.AllowListConfig.Verify(chainConfig, c.Upgrade)
}

// Equal returns true if [cfg] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}

	if c.InitialRewardConfig != nil {
		if other.InitialRewardConfig == nil {
			return false
		}
		if !c.InitialRewardConfig.Equal(other.InitialRewardConfig) {
			return false
		}
	}

	return c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

func (c *Config) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 32 * 1024,
	}

	if err := c.AllowListConfig.ToBytesWithPacker(&p); err != nil {
		return nil, err
	}

	if err := c.Upgrade.ToBytesWithPacker(&p); err != nil {
		return nil, err
	}

	p.PackBool(c.InitialRewardConfig == nil)
	if p.Err != nil {
		return nil, p.Err
	}

	if c.InitialRewardConfig != nil {
		if err := c.InitialRewardConfig.ToBytesWithPacker(&p); err != nil {
			return nil, err
		}
	}

	return p.Bytes, nil
}

func (c *Config) UnmarshalBinary(bytes []byte) error {
	p := wrappers.Packer{
		Bytes: bytes,
	}
	if err := c.AllowListConfig.FromBytesWithPacker(&p); err != nil {
		return err
	}
	if err := c.Upgrade.FromBytesWithPacker(&p); err != nil {
		return err
	}

	isNil := p.UnpackBool()
	if !isNil {
		c.InitialRewardConfig = &InitialRewardConfig{}
		if err := c.InitialRewardConfig.FromBytesWithPacker(&p); err != nil {
			return err
		}
	}

	return nil
}
