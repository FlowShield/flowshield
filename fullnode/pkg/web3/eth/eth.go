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

type CloudSlit struct {
	Client   *ethclient.Client
	Instance *contract.Slit
	Auth     *bind.TransactOpts
}

var CS *CloudSlit

func InitETH(cfg *confer.Web3) error {
	ctx := context.Background()
	client, err := ethclient.Dial(cfg.ETH.URL)
	if err != nil {
		return err
	}
	chanID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}
	contractAdd := common.HexToAddress(cfg.Contract.Token)
	instance, err := contract.NewSlit(contractAdd, client)
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chanID)
	if err != nil {
		return err
	}
	CS = &CloudSlit{
		Client:   client,
		Instance: instance,
		Auth:     auth,
	}
	return CS.stack(ctx)
}

func (c *CloudSlit) stack(ctx context.Context) error {
	logger.Infof("checking if stacked or not...")
	isDeposit, err := c.Instance.IsDeposit(&bind.CallOpts{
		From: c.Auth.From,
	}, FullNode)
	if err != nil {
		return err
	}
	if isDeposit {
		logger.Infof("you have stacked!")
		return nil
	}
	logger.Infof("you have not stacked! trying to stack...")
	// 尝试质押
	tra, err := c.Instance.Stake(c.Auth, FullNode)
	if err != nil {
		return err
	}
	rec, err := bind.WaitMined(ctx, c.Client, tra)
	if err != nil {
		return err
	}
	if rec.Status > 0 {
		logger.Infof("stack succeed !")
		return nil
	}
	return errors.New("sorry,stacked failed")
}

func Instance() *contract.Slit {
	return CS.Instance
}
