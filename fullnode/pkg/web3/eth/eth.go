package eth

import (
	"context"
	"errors"

	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
	"github.com/cloudslit/cloudslit/fullnode/pkg/contract"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	FullNode = 1

	Provider = 2
)

var (
	client   *ethclient.Client
	instance *contract.Slit
	err      error
)

func InitETH(cfg *confer.Web3) error {
	ctx := context.Background()
	client, err = ethclient.Dial(cfg.EthAddress())
	if err != nil {
		return err
	}
	chanID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}
	contractAdd := common.HexToAddress(cfg.Contract.Token)
	instance, err = contract.NewSlit(contractAdd, client)
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chanID)
	if err != nil {
		return err
	}
	logger.Infof("checking if deposited or not...")
	isDeposit, err := instance.IsDeposit(&bind.CallOpts{
		From: auth.From,
	}, FullNode)
	if err != nil {
		return err
	}
	if isDeposit {
		logger.Infof("you have deposited!")
		return nil
	}
	logger.Infof("you have not deposited! trying to deposit...")
	// 尝试质押
	tra, err := instance.Deposit(auth, FullNode)
	if err != nil {
		return err
	}
	rec, err := bind.WaitMined(ctx, client, tra)
	if err != nil {
		return err
	}
	if rec.Status > 0 {
		logger.Infof("deposited succeed !")
		return nil
	}
	return errors.New("sorry,deposit failed")
}

func Contract() *contract.Slit {
	return instance
}
