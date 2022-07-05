package service

import (
	"sort"
	"time"

	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/w3s"

	"github.com/cloudslit/cloudslit/fullnode/app/base/mdb"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/dao/api"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pconst"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

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
	// 判断服务端哨兵是否存在
	server, err := mysql.NewServer(c).GetServerByID(param.ServerID)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	if server.ID == 0 {
		return pconst.CODE_COMMON_DATA_NOT_EXIST, nil
	}
	// 查询relay
	total, relayList, err := mysql.NewRelay(c).RelayList(mparam.RelayList{
		Paginate: mdb.Paginate{
			Page:     1,
			LimitNum: 999,
		},
	})
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	data = &mmysql.Client{
		Name:     param.Name,
		ServerID: param.ServerID,
		UUID:     uuid.NewString(),
		Port:     param.Port,
		Expire:   param.Expire,
		Server: mmysql.ServerAttr{
			Name:    server.Name,
			UUID:    server.UUID,
			Host:    server.Host,
			Port:    server.Port,
			OutPort: server.OutPort,
		},
		Target: param.Target,
	}
	attrs := map[string]interface{}{
		"type":   "client",
		"name":   data.Name,
		"uuid":   data.UUID,
		"port":   data.Port,
		"server": data.Server,
		"target": data.Target,
	}
	if total > 0 && len(relayList) > 0 {
		// 默认relay顺序，ID越大越靠前
		sort.Slice(relayList, func(i, j int) bool {
			return relayList[i].ID > relayList[j].ID
		})
		relay := make([]mmysql.RelayAttrs, 0)
		for key, value := range relayList {
			relay = append(relay, mmysql.RelayAttrs{
				Name:    value.Name,
				UUID:    value.UUID,
				Host:    value.Host,
				Port:    value.Port,
				OutPort: value.OutPort,
				Sort:    key,
			})
		}
		data.Relay = relay
		attrs["relay"] = relay
	}
	sentinelSign, err := api.ApplySign(c, attrs, "zero-access", "zero-access", "zero-access", time.Now().AddDate(0, 0, param.Expire))
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	data.CaPem = sentinelSign.CaPEM
	data.CertPem = sentinelSign.CertPEM
	data.KeyPem = sentinelSign.KeyPEM
	// 先存储到w3s
	cid, err := w3s.Put(c.Request.Context(), data)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	data.Cid = cid
	err = mysql.NewClient(c).AddClient(data)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY, nil
	}
	return
}

func EditClient(c *gin.Context, param *mparam.EditClient) (code int) {
	info, err := mysql.NewClient(c).GetClientByID(param.ID)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if info.ID == 0 {
		code = pconst.CODE_COMMON_DATA_NOT_EXIST
		return
	}
	info.Name = param.Name
	info.ServerID = param.ServerID
	info.Port = param.Port
	info.Expire = param.Expire
	info.Target = param.Target

	server, err := mysql.NewServer(c).GetServerByID(param.ServerID)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if server.ID == 0 {
		code = pconst.CODE_COMMON_DATA_NOT_EXIST
		return
	}

	attrs := map[string]interface{}{
		"type":   "client",
		"name":   info.Name,
		"uuid":   info.UUID,
		"port":   info.Port,
		"relay":  info.Relay,
		"server": info.Server,
		"target": info.Target,
	}
	sentinelSign, err := api.ApplySign(c, attrs, "zero-access", "zero-access", "zero-access", time.Now().AddDate(0, 0, 90))
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	info.CaPem = sentinelSign.CaPEM
	info.CertPem = sentinelSign.CertPEM
	info.KeyPem = sentinelSign.KeyPEM
	// 先存储到w3s
	cid, err := w3s.Put(c.Request.Context(), info)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	info.Cid = cid
	err = mysql.NewClient(c).EditClient(info)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func DelClient(c *gin.Context, uuid string) (code int) {
	err := mysql.NewClient(c).DelClient(uuid)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}
