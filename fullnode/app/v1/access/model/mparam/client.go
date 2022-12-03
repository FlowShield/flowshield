package mparam

import (
	"time"

	"github.com/flowshield/flowshield/fullnode/app/base/mdb"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mmysql"
)

type ClientList struct {
	mdb.Paginate
	Name        string `json:"name" form:"name"`
	PeerID      string `json:"peer_id"`
	ResourceCID int    `json:"resource_cid" form:"resource_cid"`
	Working     bool   `json:"working" form:"working"`
}

type AddClient struct {
	PeerID string `json:"peer_id" form:"peer_id" binding:"required"`
	Name   string `json:"name" form:"name" binding:"required"`
	//Port        int    `json:"port" form:"port" binding:"required"`         // 443
	Duration    uint   `json:"duration" form:"duration" binding:"required"` // 使用时间：小时
	ResourceCID string `json:"resource_cid" binding:"required"`
}

type EditClient struct {
	ID       uint64              `json:"id" form:"id" binding:"required"`
	ServerID uint64              `json:"server_id" form:"server_id" binding:"required"`
	Name     string              `json:"name" form:"name" binding:"required"`
	Port     int                 `json:"port" form:"port" binding:"required"`     // 443
	Expire   int                 `json:"expire" form:"expire" binding:"required"` // 过期时间：天
	Target   mmysql.ClientTarget `json:"target" binding:"required"`
}

type NotifyClient struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
}

type CheckStatus struct {
	Status   uint          `json:"status"`
	Duration time.Duration `json:"duration"`
}
