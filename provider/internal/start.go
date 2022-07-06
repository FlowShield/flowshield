package internal

import (
	"context"
	"github.com/cloudslit/cloudslit/provider/internal/bll"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/initer"
	"github.com/cloudslit/cloudslit/provider/internal/schema"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/p2p"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Up struct {
	server *schema.ServerInfo
	s      sync.Mutex
}

func NewUp() *Up {
	return &Up{}
}

// SetServerPrice 设置服务启动price
func (a *Up) SetServerPrice(price int) {
	a.s.Lock()
	a.server.Price = price
	a.s.Unlock()
}

// GetServerPrice 获取服务启动price
func (a *Up) GetServerPrice() int {
	return a.server.Price
}

func InitProviderServer(ctx context.Context) {
	a := NewUp()
	err := initer.InitSelfCert()
	if err != nil {
		logger.Fatalf("Init Certificate Error: %v", err)
	}
	if config.C.P2p.Enable {
		// Create a new P2PHost
		p2phost := p2p.NewP2P(config.C.P2p.ServiceDiscoveryID)
		logrus.Infoln("Completed P2P Setup")

		// Connect to peers with the chosen discovery method
		switch config.C.P2p.ServiceDiscoverMode {
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
		a.server = a.NewServerInfo(p2phost)
		go a.starteventhandler(chatapp)
	}
	bll.NewServer().Listen(ctx)
}

// eventhandle
func (a *Up) starteventhandler(ps *p2p.PubSub) {
	refreshticker := time.NewTicker(10 * time.Second)
	defer refreshticker.Stop()
	for {
		select {
		case msg := <-ps.Inbound:
			// Print the recieved messages to the message box
			logger.Infof("Recieved Msg:%s", msg)

		case <-refreshticker.C:
			// publish
			logger.Infof("I'm:%s ", ps.Host.Host.ID().String())
			ps.Outbound <- a.server.String()
		}
	}
}

func (a *Up) NewServerInfo(p *p2p.P2P) *schema.ServerInfo {
	result := schema.ServerInfo{
		PeerId: config.C.Common.PeerId,
		Addr:   config.C.Common.LocalAddr,
		Port:   config.C.Common.LocalPort,
		Price:  config.C.Common.Price,
		Type:   schema.ServerTypeProvider,
	}
	trace, err := a.GetCftrace()
	if err != nil {
		logger.Warnf("Request Cfssl CDN Trace Error:%s", err)
	} else {
		result.MetaData = schema.MetaData{
			Ip:   trace.Ip,
			Loc:  trace.Loc,
			Colo: trace.Colo,
		}
	}
	return &result
}

type CfTrace struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

func (a *Up) GetCftrace() (*CfTrace, error) {
	resp, err := http.Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := new(CfTrace)
	sb := strings.Split(string(b), "\n")
	for _, item := range sb {
		is := strings.Split(item, "=")
		if is[0] == "ip" {
			result.Ip = is[1]
		}
		if is[0] == "loc" {
			result.Loc = is[1]
		}
		if is[0] == "colo" {
			result.Colo = is[1]
		}
	}
	return result, err
}
