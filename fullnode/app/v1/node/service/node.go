package service

import (
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	"github.com/cloudslit/cloudslit/fullnode/pkg/schema"

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
