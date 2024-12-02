// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"github.com/ava-labs/coreth/precompile/allowlist"
	"github.com/ava-labs/coreth/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileconfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// TxAllowList specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// TxAllowList with the given [admins], [enableds] and [managers] as members of the allowlist.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade: precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables TxAllowList.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (c *Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*TxAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	return c.AllowListConfig.Verify(chainConfig, c.Upgrade)
}
