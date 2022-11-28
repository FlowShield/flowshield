package service

import (
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/dao/mysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mparam"
	"github.com/cloudslit/cloudslit/fullnode/pconst"
	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
	"github.com/cloudslit/cloudslit/fullnode/pkg/util"
	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/w3s"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func ResourceList(c *gin.Context, param mparam.ResourceList) (code int, ResourceList mapi.ResourceList) {
	count, list, err := mysql.NewResource(c).ResourceList(param)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	ResourceList.List = list
	ResourceList.Paginate.Total = count
	ResourceList.Paginate.PageSize = param.LimitNum
	ResourceList.Paginate.Current = param.Page
	return
}

func AddResource(c *gin.Context, param *mparam.AddResource) (code int) {
	data := &mmysql.Resource{
		Name: param.Name,
		UUID: uuid.NewString(),
		Type: param.Type,
		Host: param.Host,
		Port: param.Port,
	}
	if user := util.User(c); user != nil {
		data.UserUUID = user.UUID
	}
	// 先存储到w3s
	account := confer.GlobalConfig().P2P.Account
	cid, err := w3s.Put(c.Request.Context(), data, data.UUID, []byte(account[len(account)-8:]))
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	data.Cid = cid
	err = mysql.NewResource(c).AddResource(data)
	if err != nil {
		return pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}

func EditResource(c *gin.Context, param *mparam.EditResource) (code int) {
	info, err := mysql.NewResource(c).GetResourceByID(param.ID)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if info.ID == 0 {
		code = pconst.CODE_COMMON_DATA_NOT_EXIST
		return
	}
	info.Name = param.Name
	info.Type = param.Type
	info.Host = param.Host
	info.Port = param.Port
	err = mysql.NewResource(c).EditResource(info)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	return
}

func DelResource(c *gin.Context, uuid string) (code int) {
	// check if any servers under this resource
	resource, err := mysql.NewResource(c).GetResourceByUUID(uuid)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
		return
	}
	if resource == nil || resource.ID == 0 {
		code = pconst.CODE_COMMON_DATA_NOT_EXIST
		return
	}
	err = mysql.NewResource(c).DelResource(uuid)
	if err != nil {
		code = pconst.CODE_COMMON_SERVER_BUSY
	}
	return
}
