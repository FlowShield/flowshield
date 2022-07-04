package test

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/cloudslit/cloudslit/fullnode/pkg/contract"
)

var client *ethclient.Client
var err error

func init() {
	client, err = ethclient.Dial("https://ropsten.infura.io/v3/45630f96f9d841679dc200a7c97763d2")
	if err != nil {
		panic(err)
	}
}

func TestETH1(t *testing.T) {
	account := common.HexToAddress("0x828233e3908fB45d40baC6B2F19F8A239ab7ae7d")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue, err)
}

func TestETH2(t *testing.T) {
	contractAddress := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bytecode)) // 60806...10029
}

func TestETH3(t *testing.T) {
	// 智能合约地址
	address := common.HexToAddress("0x4E9bfAB50AE5aA47838921450BBc1b12a81798ba")
	instance, err := contract.NewToken(address, client)
	if err != nil {
		log.Fatal(err)
	}

	bal, err := instance.BalanceOf(nil, common.HexToAddress("0x1623c4E373f80fa7B3d5E46c2F71bc50708bA5A9"))
	fmt.Println(bal, err)

	privateKey, err := crypto.HexToECDSA("xxxxxxxxxx")
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(3))
	opt := &bind.TransactOpts{
		From:     auth.From,
		Nonce:    auth.Nonce,
		Signer:   auth.Signer,
		Value:    big.NewInt(10000),
		GasPrice: big.NewInt(10000),
		GasLimit: 2381623,
	}
	ty, err := instance.Transfer(opt, common.HexToAddress("0x828233e3908fB45d40baC6B2F19F8A239ab7ae7d"), big.NewInt(100000))
	fmt.Println(ty, err)
}
