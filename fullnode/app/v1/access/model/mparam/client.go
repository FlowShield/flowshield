package mparam

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mdb"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
)

type ClientList struct {
	mdb.Paginate
	Name     string `json:"name" form:"name"`
	ServerID int    `json:"server_id" form:"server_id"`
}

type AddClient struct {
	ServerID uint64              `json:"server_id" form:"server_id" binding:"required"`
	Name     string              `json:"name" form:"name" binding:"required"`
	Port     int                 `json:"port" form:"port" binding:"required"`     // 443
	Expire   int                 `json:"expire" form:"expire" binding:"required"` // 过期时间：天
	Target   mmysql.ClientTarget `json:"target" binding:"required"`
}

type EditClient struct {
	ID       uint64              `json:"id" form:"id" binding:"required"`
	ServerID uint64              `json:"server_id" form:"server_id" binding:"required"`
	Name     string              `json:"name" form:"name" binding:"required"`
	Port     int                 `json:"port" form:"port" binding:"required"`     // 443
	Expire   int                 `json:"expire" form:"expire" binding:"required"` // 过期时间：天
	Target   mmysql.ClientTarget `json:"target" binding:"required"`
}
