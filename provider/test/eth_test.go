package test

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/cloudslit/cloudslit/provider/pkg/contract"
)

var (
	client *ethclient.Client
	err    error
)

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
	// 根据client获取chanid
	// chanID, err := client.ChainID(context.Background())
	// 智能合约地址
	address := common.HexToAddress("0xAc0A5A821d7b818f7495062e2a2FD38cEe207397")
	instance, err := contract.NewSlit(address, client)
	if err != nil {
		log.Fatal(err)
	}
	bal, err := instance.BalanceOf(nil, common.HexToAddress("0x828233e3908fB45d40baC6B2F19F8A239ab7ae7d"))
	fmt.Println(bal, err)

	//privateKey, err := crypto.HexToECDSA("8829b5d74cdfa86ae17b11d2df83f627a888fab3b86a139c6d442ef7d0e9dd76")
	//auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chanID)
	//trc, err := instance.Stake(auth)
	//fmt.Println(trc, err, "********************")
	////trc, err := instance.Transfer(auth, common.HexToAddress("0x1623c4E373f80fa7B3d5E46c2F71bc50708bA5A9"), big.NewInt(7000))
	//fmt.Println(trc.Hash().String(), err)
	//isStack, err := instance.IsStake(&bind.CallOpts{
	//	From: auth.From,
	//})
	//fmt.Println(isStack, err, "++++++++++")
}
