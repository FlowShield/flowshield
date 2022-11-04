package server

import (
	"context"

	"github.com/cloudslit/cloudslit/provider/internal/bll"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/p2p"
	"github.com/cloudslit/cloudslit/provider/pkg/web3/eth"
	"github.com/cloudslit/cloudslit/provider/pkg/web3/w3s"
	"github.com/sirupsen/logrus"
)

func InitNode(ctx context.Context) error {
	if err := w3s.Init(&config.C.Web3); err != nil {
		logger.Errorf("w3s init error : %v", err)
		return err
	}
	if err := eth.Init(&config.C.Web3); err != nil {
		logger.Errorf("eth init error : %v", err)
		return err
	}
	if err := runETH(); err != nil {
		logger.Errorf("runETH error : %v", err)
		return err
	}
	if config.C.P2p.Enable {
		// Create a new P2PHost
		logrus.Infoln("use service discovery id:", config.C.P2p.ServiceDiscoveryID)
		p2phost := p2p.NewP2P(config.C.P2p.ServiceDiscoveryID)
		logrus.Infoln("Completed P2P Setup")
		logrus.Infoln("Please wait for about 30 seconds ...")
		// Connect to peers with the chosen discovery method
		switch config.C.P2p.ServiceDiscoveryMode {
		case "announce":
			p2phost.AnnounceConnect()
		case "advertise":
			p2phost.AdvertiseConnect()
		default:
			p2phost.AdvertiseConnect()
		}
		logrus.Infoln("Connected to Service Peers")

		// Join the chat room
		chatapp, err := p2p.JoinPubSub(p2phost, "server_provider", config.C.P2p.ServiceMetadataTopic)
		if err != nil {
			logger.Fatalf("Join PubSub Error: %v", err)
		}
		logrus.Infof("Successfully joined [%s] P2P channel.", config.C.P2p.ServiceMetadataTopic)
		psBll := bll.NewPubsub()
		go psBll.InitByDB(ctx, chatapp)
		go psBll.StartPubsubHandler(ctx, chatapp, p2phost)
	}
	return nil
}
