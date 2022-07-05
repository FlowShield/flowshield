package mmysql

import (
	"database/sql/driver"
	"encoding/json"
)

type Client struct {
	ID        uint         `gorm:"primarykey"`
	CreatedAt int64        `gorm:"autoCreateTime"`
	UpdatedAt int64        `gorm:"autoUpdateTime"`
	UserUUID  string       `json:"user_uuid" gorm:"user_uuid"`
	Name      string       `json:"name"`
	ServerID  uint64       `json:"server_id"`
	UUID      string       `json:"uuid" gorm:"column:uuid"`
	Port      int          `json:"port"`
	Expire    int          `json:"expire"` // 过期时间：天
	Relay     Relays       `json:"relay"`
	Server    ServerAttr   `json:"server"`
	Target    ClientTarget `json:"target"`
	//CaPem     string       `json:"ca_pem"`
	//CertPem   string       `json:"cert_pem"`
	//KeyPem    string       `json:"key_pem"`
	Cid string `json:"cid"`
}

func (Client) TableName() string {
	return "zta_client"
}

type Relays []RelayAttrs

//type Servers ServerAttr

//type Resource []Resource

type ClientTarget struct {
	Host string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
}

func (c Relays) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Relays) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
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

type ClientAttrs struct {
	Type   string       `json:"type"`
	Name   string       `json:"name"`
	UUID   string       `json:"uuid"`
	Port   int          `json:"port"`
	Relay  []RelayAttrs `json:"relay"`
	Server ServerAttr   `json:"server"`
	Target ClientTarget `json:"target"`
}
