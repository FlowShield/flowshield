package server

import (
	"time"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/service"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util/json"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util"

	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/p2p"
	"github.com/cloudslit/cloudslit/fullnode/pkg/schema"
	"github.com/sirupsen/logrus"
)

func runP2P() error {
	cfg := confer.GlobalConfig()
	// Create a new P2PHost
	p2phost := p2p.NewP2P(cfg.P2P.ServiceDiscoveryID)
	logger.Infof("Completed P2P Setup")
	// Connect to peers with the chosen discovery method
	switch cfg.P2P.ServiceDiscoveryMode {
	case "announce":
		p2phost.AnnounceConnect()
	case "advertise":
		p2phost.AdvertiseConnect()
	default:
		p2phost.AdvertiseConnect()
	}
	logger.Infof("Connected to Service Peers")
	// Join the chat room
	pubsub, err := p2p.JoinPubSub(p2phost, "server_provider", cfg.P2P.ServiceMetadataTopic)
	if err != nil {
		logger.Errorf(nil, "Join PubSub Error: %v", err)
		return err
	}
	logrus.Infof("Successfully joined [%s] P2P channel.", cfg.P2P.ServiceMetadataTopic)
	go startEventHandler(pubsub)
	return nil
}

func startEventHandler(ps *p2p.PubSub) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	info := NewServerInfo(ps.Host)
	for {
		select {
		case msg := <-ps.Inbound:
			//p2p.Generate(msg.Message)
			service.AddNode(nil, p2p.Generate(msg.Message))
		case <-ticker.C:
			// publish
			ps.Outbound <- json.MarshalToString(info)
		}
	}
}

func NewServerInfo(p *p2p.P2P) (server *schema.ServerInfo) {
	server = &schema.ServerInfo{
		PeerId: confer.GlobalConfig().P2P.Account,
		Type:   schema.FullNode,
	}
	trace, err := util.GetCftrace()
	if err != nil {
		logger.Warnf(nil, "Request Cfssl CDN Trace Error:%s", err)
	} else {
		server.MetaData = schema.MetaData{
			Ip:   trace.Ip,
			Loc:  trace.Loc,
			Colo: trace.Colo,
		}
	}
	return
	//return json.MarshalToString(result)
}
