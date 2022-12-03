package server

import (
	"github.com/flowshield/flowshield/fullnode/agent"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/mysql"
	"github.com/flowshield/flowshield/fullnode/pkg/redis"
	"github.com/flowshield/flowshield/fullnode/pkg/web3/eth"
	"github.com/flowshield/flowshield/fullnode/pkg/web3/w3s"
	"github.com/urfave/cli"
)

func InitService(c *cli.Context) (err error) {
	if err = confer.Init(c.String("c")); err != nil {
		return
	}
	cfg := confer.GlobalConfig()
	logger.Init(&logger.Config{
		Level:       logger.LogLevel(),
		Filename:    logger.LogFile(),
		SendToFile:  logger.SendLogToFile(),
		Development: confer.ConfigEnvIsDev(),
	})
	if err = redis.Init(&cfg.Redis); err != nil {
		logger.Errorf(nil, "redis init error : %v", err)
		return
	}
	if err = mysql.Init(&cfg.Mysql); err != nil {
		logger.Errorf(nil, "mysql init error : %v", err)
		return
	}
	if err = w3s.Init(&cfg.Web3); err != nil {
		logger.Errorf(nil, "w3s init error : %v", err)
		return
	}
	if confer.GlobalConfig().Web3.Register == "true" {
		if err = eth.InitETH(&cfg.Web3); err != nil {
			logger.Errorf(nil, "eth init error : %v", err)
			return
		}
	}
	// 判断是否开启P2P
	if confer.GlobalConfig().P2P.Enable {
		if err = runP2P(&cfg.P2P); err != nil {
			logger.Errorf(nil, "p2p init error : %v", err)
			return
		}
	}
	if confer.GlobalConfig().Web3.Register == "true" && confer.GlobalConfig().P2P.Enable {
		// 启动订单同步检查agent
		go agent.SyncClientOrderStatus()
		go agent.SyncUserBindStatus()
	}
	return
}
