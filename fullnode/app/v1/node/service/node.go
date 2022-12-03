package service

import (
	"github.com/flowshield/flowshield/fullnode/app/v1/node/dao/mysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mapi"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mmysql"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mparam"
	"github.com/flowshield/flowshield/fullnode/pconst"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/flowshield/flowshield/fullnode/pkg/schema"
	"github.com/flowshield/flowshield/fullnode/pkg/web3/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func ListNode(c *gin.Context, param mparam.ListNode) (code int, nodeList mapi.NodeList) {
	count, list, err := mysql.NewNode(c).ListNode(param)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	nodeList.List = list
	nodeList.Paginate.Total = count
	nodeList.Paginate.Current = param.Page
	nodeList.Paginate.PageSize = param.LimitNum
	return
}

func AddNode(c *gin.Context, server *schema.ServerInfo) (code int) {
	if server.Type == "fullnode" {
		return
	}
	// 判断当前传递的节点是否已经质押
	isDeposit, err := eth.Instance().IsDeposit(&bind.CallOpts{
		From: common.HexToAddress(server.PeerId),
	}, eth.Provider)
	if err != nil {
		logger.Errorf(c, "check isDeposit error: %v", err)
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if !isDeposit {
		logger.Warnf(c, "provider %s has not deposited yet: %v", server.PeerId, err)
		return
	}
	node := &mmysql.Node{
		PeerId: server.PeerId,
		Addr:   server.Addr,
		Port:   server.Port,
		IP:     server.MetaData.Ip,
		Loc:    server.MetaData.Loc,
		Colo:   server.MetaData.Colo,
		Price:  server.Price,
		Type:   string(server.Type),
	}
	// 判断服务端哨兵是否存在
	data, err := mysql.NewNode(c).GetNodeByPeerId(node.PeerId)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	if data.ID > 0 {
		node.ID = data.ID
		node.CreatedAt = data.CreatedAt
		err = mysql.NewNode(c).EditNode(node)
		if err != nil {
			return pconst.CODE_COMMON_SERVER_BUSY
		}
		return
	}
	err = mysql.NewNode(c).AddNode(node)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}
