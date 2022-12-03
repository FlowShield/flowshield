package mapi

import (
	"github.com/flowshield/flowshield/fullnode/app/base/mapi"
)

type ClientList struct {
	List     []Client           `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}

type Client struct {
	ID          uint   `gorm:"primarykey"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
	UserUUID    string `json:"user_uuid" gorm:"user_uuid"`
	Name        string `json:"name"`
	PeerID      string `json:"peer_id"`
	UUID        string `json:"uuid" gorm:"column:uuid"`
	Port        int    `json:"port"`
	Duration    uint   `json:"duration"` // 使用时间：小时
	Price       uint   `json:"price"`    // 金额
	ResourceCid string `json:"resource_cid"`
	ServerCid   string `json:"server_cid"`
	ClientCid   string `json:"client_cid"`
	Status      uint   `json:"status"` // 0:待支付，1:已支付,待回调，2:已完成
	NodeIP      string `json:"node_ip"`
}
