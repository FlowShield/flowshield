package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/cloudslit/cloudslit/fullnode/pkg/util"

	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/eth"

	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/w3s"

	"github.com/cloudslit/cloudslit/fullnode/pkg/p2p"

	"github.com/cloudslit/cloudslit/fullnode/pkg/schema"

	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/dao/api"
	mysqlNode "github.com/cloudslit/cloudslit/fullnode/app/v1/node/dao/mysql"

	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pconst"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type CA struct {
	CaPem   string `json:"ca_pem"`
	CertPem string `json:"cert_pem"`
	KeyPem  string `json:"key_pem"`
}

type ServerAttr struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Message struct {
	Type string      `json:"type"` // node, order
	Data interface{} `json:"data"`
}

func ClientList(c *gin.Context, param mparam.ClientList) (code int, ClientList mapi.ClientList) {
	count, list, err := mysql.NewClient(c).ClientList(param)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	ClientList.List = list
	ClientList.Paginate.Total = count
	ClientList.Paginate.PageSize = param.LimitNum
	ClientList.Paginate.Current = param.Page
	return
}

func AddClient(c *gin.Context, param *mparam.AddClient) (code int, data *mmysql.Client) {
	// 判断Node是否存在
	node, err := mysqlNode.NewNode(c).GetNodeByPeerId(param.PeerID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	if node.ID == 0 {
		return pconst.CODE_COMMON_DATA_NOT_EXIST, nil
	}
	data = &mmysql.Client{
		Name:   param.Name,
		PeerID: param.PeerID,
		UUID:   uuid.NewString(),
		//Port:     param.Port,
		ResourceCid: param.ResourceCID,
		Duration:    param.Duration,
		Price:       param.Duration * node.Price,
	}
	// 查询Resource是否存在
	resource, err := mysql.NewResource(c).GetResourceByCID(param.ResourceCID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	if resource.ID == 0 {
		return pconst.CODE_COMMON_DATA_NOT_EXIST, nil
	}

	serverSign, err := api.ApplySign(c, map[string]interface{}{"type": "provider"}, "cloud-slit", "cloud-slit", node.Addr, time.Duration(param.Duration)*time.Hour)
	if err != nil {
		logger.Errorf(c, "AddClient ApplySign err : %v", err)
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	// 先存储到w3s
	cid, err := w3s.Put(c.Request.Context(), &CA{
		CaPem:   util.Base64Encode(serverSign.CaPEM),
		CertPem: util.Base64Encode(serverSign.CertPEM),
		KeyPem:  util.Base64Encode(serverSign.KeyPEM),
	}, data.UUID, []byte(node.PeerId[len(node.PeerId)-8:]))
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	data.ServerCid = cid
	err = mysql.NewClient(c).AddClient(data)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	return
}

//func EditClient(c *gin.Context, param *mparam.EditClient) (code int) {
//	info, err := mysql.NewClient(c).GetClientByID(param.ID)
//	if err != nil {
//		code = pconst.CODE_COMMON_SERVER_BUSY
//		return
//	}
//	if info.ID == 0 {
//		code = pconst.CODE_COMMON_DATA_NOT_EXIST
//		return
//	}
//	info.Name = param.Name
//	info.ServerID = param.ServerID
//	info.Port = param.Port
//	info.Expire = param.Expire
//	info.Target = param.Target
//
//	server, err := mysql.NewServer(c).GetServerByID(param.ServerID)
//	if err != nil {
//		code = pconst.CODE_COMMON_SERVER_BUSY
//		return
//	}
//	if server.ID == 0 {
//		code = pconst.CODE_COMMON_DATA_NOT_EXIST
//		return
//	}
//
//	//attrs := map[string]interface{}{
//	//	"type":   "client",
//	//	"name":   info.Name,
//	//	"uuid":   info.UUID,
//	//	"port":   info.Port,
//	//	"relay":  info.Relay,
//	//	"server": info.Server,
//	//	"target": info.Target,
//	//}
//	//sentinelSign, err := api.ApplySign(c, attrs, "zero-access", "zero-access", "zero-access", time.Now().AddDate(0, 0, 90))
//	//if err != nil {
//	//	code = pconst.CODE_COMMON_SERVER_BUSY
//	//	return
//	//}
//	//info.CaPem = sentinelSign.CaPEM
//	//info.CertPem = sentinelSign.CertPEM
//	//info.KeyPem = sentinelSign.KeyPEM
//	// 先存储到w3s
//	cid, err := w3s.Put(c.Request.Context(), info)
//	if err != nil {
//		return pconst.CODE_COMMON_SERVER_BUSY
//	}
//	info.Cid = cid
//	err = mysql.NewClient(c).EditClient(info)
//	if err != nil {
//		code = pconst.CODE_COMMON_SERVER_BUSY
//		return
//	}
//	return
//}

//func DelClient(c *gin.Context, uuid string) (code int) {
//	err := mysql.NewClient(c).DelClient(uuid)
//	if err != nil {
//		code = pconst.CODE_COMMON_SERVER_BUSY
//	}
//	return
//}

func AcceptClientOrder(c *gin.Context, client *schema.ClientP2P) {
	if client == nil {
		logger.Errorf(c, "AcceptClientOrder client os nil")
		return
	}
	if client.Port <= 0 {
		logger.Errorf(c, "AcceptClientOrder client port <= 0")
		return
	}
	// 根据uuid查询client
	info, err := mysql.NewClient(c).GetClientByUUID(client.UUID)
	if err != nil {
		return
	}
	if info == nil || info.ID == 0 {
		logger.Errorf(c, "AcceptClientOrder client not exist")
		return
	}
	// 判断订单状态是否已支付
	if info.Status != mmysql.Paid {
		//logger.Warnf(c, "order status is not right: %v", info.Status)
		return
	}
	// 申请客户端ca，保存至ipfs，修改订单信息
	// 查询Resource是否存在
	resource, err := mysql.NewResource(c).GetResourceByCID(info.ResourceCid)
	if err != nil {
		return
	}
	if resource == nil || resource.ID == 0 {
		logger.Errorf(c, "AcceptClientOrder resource not exist")
		return
	}
	port, err := strconv.Atoi(resource.Port)
	if err != nil {
		logger.Errorf(c, "resource.Port err : %v", err)
		return
	}
	// 查询node信息
	node, err := mysqlNode.NewNode(c).GetNodeByPeerId(info.PeerID)
	attrs := map[string]interface{}{
		"type": "client",
		"server": ServerAttr{
			Host: node.Addr,
			Port: client.Port,
		},
		"target": mmysql.ClientTarget{
			Host: resource.Host,
			Port: port,
		},
	}
	clientSign, err := api.ApplySign(c, attrs, "cloud-slit", "cloud-slit", "cloud-slit", time.Duration(info.Duration)*time.Hour)
	if err != nil {
		logger.Errorf(c, "AcceptClientOrder ApplySign err : %v", err)
		return
	}
	// 先存储到w3s
	cid, err := w3s.Put(context.Background(), &CA{
		CaPem:   util.Base64Encode(clientSign.CaPEM),
		CertPem: util.Base64Encode(clientSign.CertPEM),
		KeyPem:  util.Base64Encode(clientSign.KeyPEM),
	}, info.UUID, []byte(node.PeerId[len(node.PeerId)-8:]))
	info.Port = client.Port
	info.ClientCid = cid
	info.Status = mmysql.Success
	err = mysql.NewClient(c).EditClient(info)
	if err != nil {
		logger.Errorf(c, "AcceptClientOrder EditClient err : %v", err)
		return
	}
	return
}

func NotifyClient(c *gin.Context, param *mparam.NotifyClient) (code int) {
	client, err := mysql.NewClient(c).GetClientByUUID(param.UUID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if client == nil || client.ID == 0 {
		return pconst.CODE_COMMON_DATA_NOT_EXIST
	}
	// 如果已经成功，则忽略
	if client.Status == mmysql.Success {
		logger.Infof("client status is success, ignore")
		return
	}
	// 查询合约判断该笔订单是否支付
	check, err := eth.Instance().CheckOrder(&bind.CallOpts{
		From: eth.CS.Auth.From,
	}, param.UUID)
	if err != nil {
		logger.Errorf(c, "notifyClient CheckOrder err : %v", err)
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if !check {
		logger.Warnf(c, "notifyClient Order: %s has not been paid", param.UUID)
		return pconst.CODE_DATA_WRONG_STATE
	}
	client.Status = mmysql.Paid
	err = mysql.NewClient(c).EditClient(client)
	if err != nil {
		return pconst.CODE_COMMON_DATA_NOT_EXIST
	}
	// 通知相应的矿工
	message := &Message{
		Type: "order",
		Data: &schema.ClientP2P{
			ServerCID: client.ServerCid,
			Wallet:    client.PeerID,
			UUID:      client.UUID,
		},
	}
	bytesData, _ := json.Marshal(message)
	p2p.GetPubSub().Outbound <- string(bytesData)
	return
}
