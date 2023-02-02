// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

// AllowList is a precompile that allows the admin to set the permissions of
// addresses. AllowList itself is not a precompile, but rather a helper
// function that is used by other precompiles to determine if an address has
// permissions.
// The permissions are:
// 1. No role assigned - this is equivalent to common.Hash{} and deletes the key from the DB when set
// 2. Enabled are allowed to take certain actions
// 3. Admin - allowed to modify both the admin and enabled list as well as take certain actions

const (
	SetAdminFuncKey      = "setAdmin"
	SetEnabledFuncKey    = "setEnabled"
	SetNoneFuncKey       = "setNone"
	ReadAllowListFuncKey = "readAllowList"

	ModifyAllowListGasCost = precompile.WriteGasCostPerSlot
	ReadAllowListGasCost   = precompile.ReadGasCostPerSlot
)

var (
	NoRole      Role = Role(common.BigToHash(big.NewInt(0))) // No role assigned - this is equivalent to common.Hash{} and deletes the key from the DB when set
	EnabledRole Role = Role(common.BigToHash(big.NewInt(1))) // Deployers are allowed to create new contracts
	AdminRole   Role = Role(common.BigToHash(big.NewInt(2))) // Admin - allowed to modify both the admin and deployer list as well as deploy contracts

	// AllowList function signatures
	setAdminSignature      = precompile.CalculateFunctionSelector("setAdmin(address)")
	setEnabledSignature    = precompile.CalculateFunctionSelector("setEnabled(address)")
	setNoneSignature       = precompile.CalculateFunctionSelector("setNone(address)")
	readAllowListSignature = precompile.CalculateFunctionSelector("readAllowList(address)")
	// Error returned when an invalid write is attempted
	ErrCannotModifyAllowList = errors.New("non-admin cannot modify allow list")

	allowListInputLen = common.HashLength
)

// GetAllowListStatus returns the allow list role of [address] for the precompile
// at [precompileAddr]
func GetAllowListStatus(state precompile.StateDB, precompileAddr common.Address, address common.Address) Role {
	// Generate the state key for [address]
	addressKey := address.Hash()
	return Role(state.GetState(precompileAddr, addressKey))
}

// SetAllowListRole sets the permissions of [address] to [role] for the precompile
// at [precompileAddr].
// assumes [role] has already been verified as valid.
func SetAllowListRole(stateDB precompile.StateDB, precompileAddr, address common.Address, role Role) {
	// Generate the state key for [address]
	addressKey := address.Hash()
	// Assign [role] to the address
	// This stores the [role] in the contract storage with address [precompileAddr]
	// and [addressKey] hash. It means that any reusage of the [addressKey] for different value
	// conflicts with the same slot [role] is stored.
	// Precompile implementations must use a different key than [addressKey]
	stateDB.SetState(precompileAddr, addressKey, common.Hash(role))
}

// PackModifyAllowList packs [address] and [role] into the appropriate arguments for modifying the allow list.
// Note: [role] is not packed in the input value returned, but is instead used as a selector for the function
// selector that should be encoded in the input.
func PackModifyAllowList(address common.Address, role Role) ([]byte, error) {
	// function selector (4 bytes) + hash for address
	input := make([]byte, 0, precompile.SelectorLen+common.HashLength)

	switch role {
	case AdminRole:
		input = append(input, setAdminSignature...)
	case EnabledRole:
		input = append(input, setEnabledSignature...)
	case NoRole:
		input = append(input, setNoneSignature...)
	default:
		return nil, fmt.Errorf("cannot pack modify list input with invalid role: %s", role)
	}

	input = append(input, address.Hash().Bytes()...)
	return input, nil
}

// PackReadAllowList packs [address] into the input data to the read allow list function
func PackReadAllowList(address common.Address) []byte {
	input := make([]byte, 0, precompile.SelectorLen+common.HashLength)
	input = append(input, readAllowListSignature...)
	input = append(input, address.Hash().Bytes()...)
	return input
}

// createAllowListRoleSetter returns an execution function for setting the allow list status of the input address argument to [role].
// This execution function is speciifc to [precompileAddr].
func createAllowListRoleSetter(precompileAddr common.Address, role Role) precompile.RunStatefulPrecompileFunc {
	return func(evm precompile.PrecompileAccessibleState, callerAddr, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = precompile.DeductGas(suppliedGas, ModifyAllowListGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != allowListInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for modifying allow list: %d", len(input))
		}

		modifyAddress := common.BytesToAddress(input)

		if readOnly {
			return nil, remainingGas, vmerrs.ErrWriteProtection
		}

		stateDB := evm.GetStateDB()

		// Verify that the caller is in the allow list and therefore has the right to call this function.
		callerStatus := GetAllowListStatus(stateDB, precompileAddr, callerAddr)
		if !callerStatus.IsAdmin() {
			return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotModifyAllowList, callerAddr)
		}

		SetAllowListRole(stateDB, precompileAddr, modifyAddress, role)
		// Return an empty output and the remaining gas
		return []byte{}, remainingGas, nil
	}
}

// createReadAllowList returns an execution function that reads the allow list for the given [precompileAddr].
// The execution function parses the input into a single address and returns the 32 byte hash that specifies the
// designated role of that address
func createReadAllowList(precompileAddr common.Address) precompile.RunStatefulPrecompileFunc {
	return func(evm precompile.PrecompileAccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = precompile.DeductGas(suppliedGas, ReadAllowListGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != allowListInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for read allow list: %d", len(input))
		}

		readAddress := common.BytesToAddress(input)
		role := GetAllowListStatus(evm.GetStateDB(), precompileAddr, readAddress)
		roleBytes := common.Hash(role).Bytes()
		return roleBytes, remainingGas, nil
	}
}

// createAllowListPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr]
func CreateAllowListPrecompile(precompileAddr common.Address) precompile.StatefulPrecompiledContract {
	// Construct the contract with no fallback function.
	allowListFuncs := CreateAllowListFunctions(precompileAddr)
	contract, err := precompile.NewStatefulPrecompileContract(nil, allowListFuncs)
	// TODO Change this to be returned as an error after refactoring this precompile
	// to use the new precompile template.
	if err != nil {
		panic(err)
	}
	return contract
}

func CreateAllowListFunctions(precompileAddr common.Address) []*precompile.StatefulPrecompileFunction {
	setAdmin := precompile.NewStatefulPrecompileFunction(setAdminSignature, createAllowListRoleSetter(precompileAddr, AdminRole))
	setEnabled := precompile.NewStatefulPrecompileFunction(setEnabledSignature, createAllowListRoleSetter(precompileAddr, EnabledRole))
	setNone := precompile.NewStatefulPrecompileFunction(setNoneSignature, createAllowListRoleSetter(precompileAddr, NoRole))
	read := precompile.NewStatefulPrecompileFunction(readAllowListSignature, createReadAllowList(precompileAddr))

	return []*precompile.StatefulPrecompileFunction{setAdmin, setEnabled, setNone, read}
}
