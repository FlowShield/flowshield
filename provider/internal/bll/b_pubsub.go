package bll

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/dao/provider/model"
	"github.com/cloudslit/cloudslit/provider/internal/dao/provider/service"
	"github.com/cloudslit/cloudslit/provider/internal/schema"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/p2p"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
)

// Pubsub
type Pubsub struct {
	Orders map[string]*schema.NodeOrder
	mu     sync.RWMutex
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		Orders: make(map[string]*schema.NodeOrder),
		mu:     sync.RWMutex{},
	}
}

// 检查数据库是否存在
func (a *Pubsub) InitByDB(ctx context.Context, ps *p2p.PubSub) {
	// 查询数据库
	list, err := service.ListProvider(&model.Provider{})
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("读取持久化数据失败：%v", err)
		return
	}
	for _, item := range list {
		if _, ok := a.getOrder(item.Uuid); ok {
			continue
		}
		if item.Wallet != config.C.Web3.Account {
			logger.Warnf("The persistent order information is inconsistent with the node information and will not be started")
			continue
		}
		order := item.ToNodeOrder()
		// 检测端口是否异常
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(order.Port))
		if err != nil {
			logger.WithErrorStack(ctx, err).Errorf("The order is restarted, and the port is abnormal. Please deal with it in time, Err:%s", err)
			continue
		}
		_ = ln.Close()
		err = a.handleOrder(ctx, ps, order)
		if err != nil {
			logger.WithErrorStack(ctx, err).Errorf("Init Handle Order Err:%s", err)
			continue
		}
	}
}

// eventhandle
func (a *Pubsub) StartPubsubHandler(ctx context.Context, ps *p2p.PubSub, p *p2p.P2P) {
	//go func() {
	//	msg := `{"type":"order","data":{"server_cid":"bafybeifnc734brtng4tn2wxfp7pjtjz7qicd2dqvsxs36shugxoogzu4zu","wallet":"0x0e5518bfef2b0e0c6600742c662b797445020F99","uuid":"cde5260e-47ac-4a07-88b4-9a7ffc357a0b","port":0}}`
	//	err := a.ReceiveHandle(ctx, ps, msg)
	//	if err != nil {
	//		logger.Errorf("Receive Msg Handle Err:%s", err)
	//	}
	//}()
	server := a.NewServerInfo()
	refreshticker := time.NewTicker(10 * time.Second)
	defer refreshticker.Stop()
	for {
		select {
		case msg := <-ps.Inbound:
			// Print the recieved messages to the message box
			err := a.ReceiveHandle(ctx, ps, msg.Message)
			if err != nil {
				logger.WithErrorStack(ctx, err).Errorf("Receive Msg Handle Err:%s", err)
			}

		case <-refreshticker.C:
			// Timing publish
			a.nodeHeartBeat(ps, server)

		case log := <-ps.Logs:
			// Add the log to the message box
			logger.Warnf("PubSub Log:%s", log)
		}
	}
}

func (a *Pubsub) ReceiveHandle(ctx context.Context, ps *p2p.PubSub, msg string) error {
	logger.Debugf("Received Msg:%s", msg)
	var pss schema.PsMessage
	err := json.Unmarshal([]byte(msg), &pss)
	if err != nil {
		return errors.WithStack(err)
	}
	if pss.Type == "order" {
		err := a.orderReceive(ctx, ps, &pss)
		if err != nil {
			return err
		}
	}
	return nil
}

// orderReceive 接收订单信息
func (a *Pubsub) orderReceive(ctx context.Context, ps *p2p.PubSub, pss *schema.PsMessage) error {
	order, err := pss.ToNodeOrder()
	if err != nil {
		return err
	}
	return a.handleOrder(ctx, ps, order)
}

// 处理订单
func (a *Pubsub) handleOrder(ctx context.Context, ps *p2p.PubSub, order *schema.NodeOrder) error {
	// 包含Port 不处理心跳包
	if order.IsHeart {
		return nil
	}
	// 处理订单
	logger.WithContext(ctx).Infof("handle Order: %s", order.String())
	if _, ok := a.getOrder(order.Uuid); ok {
		return fmt.Errorf("this order has been launched:%s", order.Uuid)
	}
	if order.Wallet != config.C.Web3.Account {
		return fmt.Errorf("wallet Abnormal，expect:%s, get:%s", config.C.Web3.Account, order.Wallet)
	}
	// 解析 w3s 配置
	pc, err := a.ParseW3sData(ctx, order)
	if err != nil {
		return err
	}
	// 预制端口
	port, err := verifyPort(order.Port)
	if err != nil {
		return errors.WithStack(err)
	}
	order.Port = port
	// 入库
	err = a.addProvider(model.OrderToProvider(order, pc))
	if err != nil {
		return errors.WithStack(err)
	}
	// 监听端口
	p := NewProvider()
	err = p.Listen(ctx, port, pc)
	if err != nil {
		return errors.WithStack(err)
	}
	go p.Handle(ctx)
	go a.providerHeartBeat(ctx, ps, order, p, pc)
	a.setOrder(order)
	logger.WithContext(ctx).Infof("Handle Order Success: %s", order.String())
	return nil
}

func (a *Pubsub) addProvider(item *model.Provider) error {
	err := service.AddProvider(item)
	return err
}

// Provider 心跳
func (a *Pubsub) providerHeartBeat(ctx context.Context, ps *p2p.PubSub, order *schema.NodeOrder, p *Provider, pc *schema.ProviderConfig) {
	order.IsHeart = true
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 检测证书到期，关闭服务
			if time.Now().After(pc.CertConfig.NotAfter) {
				p.Close()
				a.closeOrder(order)
				logger.Warnf("Provider Service expired, Closing soon: %s", order.String())
				return
			}
			// 服务心跳包
			pub := &schema.PsMessage{
				Type: schema.PsMsgTypeOrder,
				Data: order,
			}
			str := pub.String()
			ps.InsertOutbound(str)
			logger.Debugf("Provider Heart Beat - running at [::]:%d", order.Port)
		}
	}
}

func (a *Pubsub) getOrder(uuid string) (*schema.NodeOrder, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.Orders[uuid]
	return v, ok
}

func (a *Pubsub) setOrder(order *schema.NodeOrder) {
	a.mu.Lock()
	a.Orders[order.Uuid] = order
	defer a.mu.Unlock()
}

func (a *Pubsub) closeOrder(order *schema.NodeOrder) {
	a.mu.Lock()
	delete(a.Orders, order.Uuid)
	defer a.mu.Unlock()
}

// nodeHeartBeat 节点发布自身信息
func (a *Pubsub) nodeHeartBeat(ps *p2p.PubSub, server *schema.NodeInfo) {
	logger.Debugf("Node Heart Beat - PeerId:%s", ps.Host.Host.ID().String())
	pub := &schema.PsMessage{
		Type: schema.PsMsgTypeNode,
		Data: server,
	}
	str := pub.String()
	ps.InsertOutbound(str)
}

func (a *Pubsub) ParseW3sData(ctx context.Context, order *schema.NodeOrder) (*schema.ProviderConfig, error) {
	cfg := config.C.Web3
	// 解析配置
	tryCount := 0
retry:
	key := []byte(config.C.Web3.Account[len(cfg.Account)-8:])
	pc, err := order.ToProviderConfig(ctx, key)
	if err != nil {
		tryCount++
		logger.Warnf("get w3s data err:%v", err)
		if tryCount > cfg.W3S.RetryCount {
			return nil, err
		}
		goto retry
	}
	return pc, nil
}

func (a *Pubsub) NewServerInfo() *schema.NodeInfo {
	result := schema.NodeInfo{
		PeerId: config.C.Web3.Account,
		Addr:   config.C.App.LocalAddr,
		Port:   config.C.App.LocalPort,
		Price:  config.C.Web3.Price,
		Type:   schema.NodeTypeProvider,
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

func (a *Pubsub) GetCftrace() (*CfTrace, error) {
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
