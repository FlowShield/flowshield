package eth

import (
	"context"
	"math/big"

	"github.com/cloudslit/cloudslit/provider/internal/config"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

func Init(cfg *config.Web3) (err error) {
	client, err = ethclient.Dial(cfg.EthAddress())
	if err != nil {
		return err
	}
	return nil
}

func GetEthBalance(ctx context.Context, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func GetContractBalance(address string) (*big.Int, error) {
	return nil, nil
}
