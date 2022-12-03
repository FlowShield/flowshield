package p2p

import (
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
)

var pubSubObj *PubSub

func InitP2P(cfg *confer.P2P) (err error) {
	logger.Infof("Starting P2P...")
	p2phost := NewP2P(cfg.ServiceDiscoveryID)
	logger.Infof("Completed P2P Setup")
	// Connect to peers with the chosen discovery method
	switch cfg.ServiceDiscoveryMode {
	case "announce":
		p2phost.AnnounceConnect()
	case "advertise":
		p2phost.AdvertiseConnect()
	default:
		p2phost.AdvertiseConnect()
	}
	logger.Infof("Connected to Service Peers")
	// Join the chat room
	pubSubObj, err = JoinPubSub(p2phost, "server_provider", cfg.ServiceMetadataTopic)
	if err != nil {
		logger.Errorf(nil, "Join PubSub Error: %v", err)
		return err
	}
	logger.Infof("Successfully joined [%s] P2P channel.", cfg.ServiceMetadataTopic)
	return
}

func GetPubSub() *PubSub {
	return pubSubObj
}
