// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SlitMetaData contains all meta data concerning the Slit contract.
var SlitMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"orderId\",\"type\":\"string\"}],\"name\":\"checkOrder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"orderId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_orderPrice\",\"type\":\"uint256\"}],\"name\":\"clientOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fullnodeDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_type\",\"type\":\"uint8\"}],\"name\":\"isDeposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"privoderDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_type\",\"type\":\"uint8\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_type\",\"type\":\"uint8\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SlitABI is the input ABI used to generate the binding from.
// Deprecated: Use SlitMetaData.ABI instead.
var SlitABI = SlitMetaData.ABI

// Slit is an auto generated Go binding around an Ethereum contract.
type Slit struct {
	SlitCaller     // Read-only binding to the contract
	SlitTransactor // Write-only binding to the contract
	SlitFilterer   // Log filterer for contract events
}

// SlitCaller is an auto generated read-only Go binding around an Ethereum contract.
type SlitCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SlitTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SlitTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SlitFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SlitFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SlitSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SlitSession struct {
	Contract     *Slit             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SlitCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SlitCallerSession struct {
	Contract *SlitCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SlitTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SlitTransactorSession struct {
	Contract     *SlitTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SlitRaw is an auto generated low-level Go binding around an Ethereum contract.
type SlitRaw struct {
	Contract *Slit // Generic contract binding to access the raw methods on
}

// SlitCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SlitCallerRaw struct {
	Contract *SlitCaller // Generic read-only contract binding to access the raw methods on
}

// SlitTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SlitTransactorRaw struct {
	Contract *SlitTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSlit creates a new instance of Slit, bound to a specific deployed contract.
func NewSlit(address common.Address, backend bind.ContractBackend) (*Slit, error) {
	contract, err := bindSlit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Slit{SlitCaller: SlitCaller{contract: contract}, SlitTransactor: SlitTransactor{contract: contract}, SlitFilterer: SlitFilterer{contract: contract}}, nil
}

// NewSlitCaller creates a new read-only instance of Slit, bound to a specific deployed contract.
func NewSlitCaller(address common.Address, caller bind.ContractCaller) (*SlitCaller, error) {
	contract, err := bindSlit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SlitCaller{contract: contract}, nil
}

// NewSlitTransactor creates a new write-only instance of Slit, bound to a specific deployed contract.
func NewSlitTransactor(address common.Address, transactor bind.ContractTransactor) (*SlitTransactor, error) {
	contract, err := bindSlit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SlitTransactor{contract: contract}, nil
}

// NewSlitFilterer creates a new log filterer instance of Slit, bound to a specific deployed contract.
func NewSlitFilterer(address common.Address, filterer bind.ContractFilterer) (*SlitFilterer, error) {
	contract, err := bindSlit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SlitFilterer{contract: contract}, nil
}

// bindSlit binds a generic wrapper to an already deployed contract.
func bindSlit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SlitABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Slit *SlitRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Slit.Contract.SlitCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Slit *SlitRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Slit.Contract.SlitTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Slit *SlitRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Slit.Contract.SlitTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Slit *SlitCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Slit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Slit *SlitTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Slit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Slit *SlitTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Slit.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Slit *SlitCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Slit *SlitSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Slit.Contract.Allowance(&_Slit.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Slit *SlitCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Slit.Contract.Allowance(&_Slit.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Slit *SlitCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Slit *SlitSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Slit.Contract.BalanceOf(&_Slit.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Slit *SlitCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Slit.Contract.BalanceOf(&_Slit.CallOpts, account)
}

// CheckOrder is a free data retrieval call binding the contract method 0xdd6887ed.
//
// Solidity: function checkOrder(string orderId) view returns(bool)
func (_Slit *SlitCaller) CheckOrder(opts *bind.CallOpts, orderId string) (bool, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "checkOrder", orderId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckOrder is a free data retrieval call binding the contract method 0xdd6887ed.
//
// Solidity: function checkOrder(string orderId) view returns(bool)
func (_Slit *SlitSession) CheckOrder(orderId string) (bool, error) {
	return _Slit.Contract.CheckOrder(&_Slit.CallOpts, orderId)
}

// CheckOrder is a free data retrieval call binding the contract method 0xdd6887ed.
//
// Solidity: function checkOrder(string orderId) view returns(bool)
func (_Slit *SlitCallerSession) CheckOrder(orderId string) (bool, error) {
	return _Slit.Contract.CheckOrder(&_Slit.CallOpts, orderId)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Slit *SlitCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Slit *SlitSession) Decimals() (uint8, error) {
	return _Slit.Contract.Decimals(&_Slit.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Slit *SlitCallerSession) Decimals() (uint8, error) {
	return _Slit.Contract.Decimals(&_Slit.CallOpts)
}

// FullnodeDepositAmount is a free data retrieval call binding the contract method 0xd35b1ac2.
//
// Solidity: function fullnodeDepositAmount() view returns(uint256)
func (_Slit *SlitCaller) FullnodeDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "fullnodeDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FullnodeDepositAmount is a free data retrieval call binding the contract method 0xd35b1ac2.
//
// Solidity: function fullnodeDepositAmount() view returns(uint256)
func (_Slit *SlitSession) FullnodeDepositAmount() (*big.Int, error) {
	return _Slit.Contract.FullnodeDepositAmount(&_Slit.CallOpts)
}

// FullnodeDepositAmount is a free data retrieval call binding the contract method 0xd35b1ac2.
//
// Solidity: function fullnodeDepositAmount() view returns(uint256)
func (_Slit *SlitCallerSession) FullnodeDepositAmount() (*big.Int, error) {
	return _Slit.Contract.FullnodeDepositAmount(&_Slit.CallOpts)
}

// IsDeposit is a free data retrieval call binding the contract method 0xf276b8aa.
//
// Solidity: function isDeposit(uint8 _type) view returns(bool)
func (_Slit *SlitCaller) IsDeposit(opts *bind.CallOpts, _type uint8) (bool, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "isDeposit", _type)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDeposit is a free data retrieval call binding the contract method 0xf276b8aa.
//
// Solidity: function isDeposit(uint8 _type) view returns(bool)
func (_Slit *SlitSession) IsDeposit(_type uint8) (bool, error) {
	return _Slit.Contract.IsDeposit(&_Slit.CallOpts, _type)
}

// IsDeposit is a free data retrieval call binding the contract method 0xf276b8aa.
//
// Solidity: function isDeposit(uint8 _type) view returns(bool)
func (_Slit *SlitCallerSession) IsDeposit(_type uint8) (bool, error) {
	return _Slit.Contract.IsDeposit(&_Slit.CallOpts, _type)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Slit *SlitCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Slit *SlitSession) Name() (string, error) {
	return _Slit.Contract.Name(&_Slit.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Slit *SlitCallerSession) Name() (string, error) {
	return _Slit.Contract.Name(&_Slit.CallOpts)
}

// PrivoderDepositAmount is a free data retrieval call binding the contract method 0xcd74c0a7.
//
// Solidity: function privoderDepositAmount() view returns(uint256)
func (_Slit *SlitCaller) PrivoderDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "privoderDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PrivoderDepositAmount is a free data retrieval call binding the contract method 0xcd74c0a7.
//
// Solidity: function privoderDepositAmount() view returns(uint256)
func (_Slit *SlitSession) PrivoderDepositAmount() (*big.Int, error) {
	return _Slit.Contract.PrivoderDepositAmount(&_Slit.CallOpts)
}

// PrivoderDepositAmount is a free data retrieval call binding the contract method 0xcd74c0a7.
//
// Solidity: function privoderDepositAmount() view returns(uint256)
func (_Slit *SlitCallerSession) PrivoderDepositAmount() (*big.Int, error) {
	return _Slit.Contract.PrivoderDepositAmount(&_Slit.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Slit *SlitCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Slit *SlitSession) Symbol() (string, error) {
	return _Slit.Contract.Symbol(&_Slit.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Slit *SlitCallerSession) Symbol() (string, error) {
	return _Slit.Contract.Symbol(&_Slit.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Slit *SlitCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Slit.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Slit *SlitSession) TotalSupply() (*big.Int, error) {
	return _Slit.Contract.TotalSupply(&_Slit.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Slit *SlitCallerSession) TotalSupply() (*big.Int, error) {
	return _Slit.Contract.TotalSupply(&_Slit.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Slit *SlitTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Slit *SlitSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.Approve(&_Slit.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Slit *SlitTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.Approve(&_Slit.TransactOpts, spender, value)
}

// ClientOrder is a paid mutator transaction binding the contract method 0x1b005437.
//
// Solidity: function clientOrder(string orderId, uint256 _orderPrice) returns()
func (_Slit *SlitTransactor) ClientOrder(opts *bind.TransactOpts, orderId string, _orderPrice *big.Int) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "clientOrder", orderId, _orderPrice)
}

// ClientOrder is a paid mutator transaction binding the contract method 0x1b005437.
//
// Solidity: function clientOrder(string orderId, uint256 _orderPrice) returns()
func (_Slit *SlitSession) ClientOrder(orderId string, _orderPrice *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.ClientOrder(&_Slit.TransactOpts, orderId, _orderPrice)
}

// ClientOrder is a paid mutator transaction binding the contract method 0x1b005437.
//
// Solidity: function clientOrder(string orderId, uint256 _orderPrice) returns()
func (_Slit *SlitTransactorSession) ClientOrder(orderId string, _orderPrice *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.ClientOrder(&_Slit.TransactOpts, orderId, _orderPrice)
}

// Stake is a paid mutator transaction binding the contract method 0x604f2177.
//
// Solidity: function stake(uint8 _type) returns()
func (_Slit *SlitTransactor) Stake(opts *bind.TransactOpts, _type uint8) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "stake", _type)
}

// Stake is a paid mutator transaction binding the contract method 0x604f2177.
//
// Solidity: function stake(uint8 _type) returns()
func (_Slit *SlitSession) Stake(_type uint8) (*types.Transaction, error) {
	return _Slit.Contract.Stake(&_Slit.TransactOpts, _type)
}

// Stake is a paid mutator transaction binding the contract method 0x604f2177.
//
// Solidity: function stake(uint8 _type) returns()
func (_Slit *SlitTransactorSession) Stake(_type uint8) (*types.Transaction, error) {
	return _Slit.Contract.Stake(&_Slit.TransactOpts, _type)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Slit *SlitTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Slit *SlitSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.Transfer(&_Slit.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Slit *SlitTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.Transfer(&_Slit.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Slit *SlitTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Slit *SlitSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.TransferFrom(&_Slit.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Slit *SlitTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Slit.Contract.TransferFrom(&_Slit.TransactOpts, from, to, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc6ab5d90.
//
// Solidity: function withdraw(uint8 _type) returns()
func (_Slit *SlitTransactor) Withdraw(opts *bind.TransactOpts, _type uint8) (*types.Transaction, error) {
	return _Slit.contract.Transact(opts, "withdraw", _type)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc6ab5d90.
//
// Solidity: function withdraw(uint8 _type) returns()
func (_Slit *SlitSession) Withdraw(_type uint8) (*types.Transaction, error) {
	return _Slit.Contract.Withdraw(&_Slit.TransactOpts, _type)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc6ab5d90.
//
// Solidity: function withdraw(uint8 _type) returns()
func (_Slit *SlitTransactorSession) Withdraw(_type uint8) (*types.Transaction, error) {
	return _Slit.Contract.Withdraw(&_Slit.TransactOpts, _type)
}

// SlitApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Slit contract.
type SlitApprovalIterator struct {
	Event *SlitApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SlitApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SlitApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SlitApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SlitApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SlitApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SlitApproval represents a Approval event raised by the Slit contract.
type SlitApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Slit *SlitFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SlitApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Slit.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SlitApprovalIterator{contract: _Slit.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Slit *SlitFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SlitApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Slit.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SlitApproval)
				if err := _Slit.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Slit *SlitFilterer) ParseApproval(log types.Log) (*SlitApproval, error) {
	event := new(SlitApproval)
	if err := _Slit.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SlitTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Slit contract.
type SlitTransferIterator struct {
	Event *SlitTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SlitTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SlitTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SlitTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SlitTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SlitTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SlitTransfer represents a Transfer event raised by the Slit contract.
type SlitTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Slit *SlitFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SlitTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Slit.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SlitTransferIterator{contract: _Slit.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Slit *SlitFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SlitTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Slit.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SlitTransfer)
				if err := _Slit.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Slit *SlitFilterer) ParseTransfer(log types.Log) (*SlitTransfer, error) {
	event := new(SlitTransfer)
	if err := _Slit.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
