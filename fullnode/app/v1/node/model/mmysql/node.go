package mmysql

import (
	"database/sql/driver"
	"encoding/json"
)

type Node struct {
	ID        uint   `gorm:"primarykey"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
	PeerId    string `json:"peer_id"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	IP        string `json:"ip"`
	Loc       string `json:"loc"`
	Colo      string `json:"colo"`
	Price     uint   `json:"price"`
	Type      string `json:"type"`
}

func (Node) TableName() string {
	return "zta_node"
}

type MetaData struct {
	Ip   string `json:"ip"`
	Loc  string `json:"loc"`
	Colo string `json:"colo"`
}

func (c MetaData) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *MetaData) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
