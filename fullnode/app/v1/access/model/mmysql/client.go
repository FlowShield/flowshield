package mmysql

import (
	"database/sql/driver"
	"encoding/json"
)

const (
	WaitingPaid = iota
	Paid
	Success
)

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
}

func (Client) TableName() string {
	return "zta_client"
}

//type Servers ServerAttr

//type Resource []Resource

type ClientTarget struct {
	Host string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
}

//func (c Resource) Value() (driver.Value, error) {
//	b, err := json.Marshal(c)
//	return string(b), err
//}
//
//func (c *Resource) Scan(input interface{}) error {
//	return json.Unmarshal(input.([]byte), c)
//}

func (c ClientTarget) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ClientTarget) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

//type ClientAttrs struct {
//	Type string `json:"type"`
//	Name string `json:"name"`
//	UUID string `json:"uuid"`
//	Port int    `json:"port"`
//	//Relay  []RelayAttrs `json:"relay"`
//	Server ServerAttr   `json:"server"`
//	Target ClientTarget `json:"target"`
//}
