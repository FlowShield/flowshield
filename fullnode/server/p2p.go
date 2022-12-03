package server

import (
	"encoding/json"

	service2 "github.com/flowshield/flowshield/fullnode/app/v1/access/service"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/service"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/p2p"
	"github.com/flowshield/flowshield/fullnode/pkg/schema"
	"github.com/tidwall/gjson"
)

func runP2P(cfg *confer.P2P) error {
	// Create a new P2PHost
	if err := p2p.InitP2P(cfg); err != nil {
		return err
	}
	go startEventHandler(p2p.GetPubSub())
	return nil
}

func startEventHandler(ps *p2p.PubSub) {
	//ticker := time.NewTicker(time.Second * 10)
	//defer ticker.Stop()
	//info := NewServerInfo(ps.Host)
	for {
		select {
		case msg := <-ps.Inbound:
			//p2p.Generate(msg.Message)
			HandleMessage(msg.Message)
		//case <-ticker.C:
		//	// publish
		//	ps.Outbound <- json.MarshalToString(info)
		case logData := <-ps.Logs:
			logger.Errorf(nil, "p2p receive error: %s", logData.String())
		}
	}
}

func HandleMessage(message string) {
	// 判断消息类型，是属于节点通信，还是订单通信
	messageType := gjson.Get(message, "type")
	switch messageType.String() {
	case "node":
		// 节点通信
		service.AddNode(nil, generateNode(gjson.Get(message, "data").String()))
	case "order":
		service2.AcceptClientOrder(nil, generateClient(gjson.Get(message, "data").String()))
	default:

	}
}

func generateNode(node string) (server *schema.ServerInfo) {
	_ = json.Unmarshal([]byte(node), &server)
	return
}

func generateClient(client string) (info *schema.ClientP2P) {
	_ = json.Unmarshal([]byte(client), &info)
	return
}

//func NewServerInfo(p *p2p.P2P) (server *schema.ServerInfo) {
//	server = &schema.ServerInfo{
//		PeerId: confer.GlobalConfig().P2P.Account,
//		Type:   schema.FullNode,
//	}
//	trace, err := util.GetCftrace()
//	if err != nil {
//		logger.Warnf(nil, "Request Cfssl CDN Trace Error:%s", err)
//	} else {
//		server.MetaData = schema.MetaData{
//			Ip:   trace.Ip,
//			Loc:  trace.Loc,
//			Colo: trace.Colo,
//		}
//	}
//	return
//}
