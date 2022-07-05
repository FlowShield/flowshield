package mmysql

import (
	"database/sql/driver"
	"encoding/json"
)

type Server struct {
	ID         uint   `gorm:"primarykey"`
	ResourceID string `json:"resource_id"`
	UserUUID   string `json:"user_uuid" gorm:"user_uuid"`
	Name       string `json:"name"`
	UUID       string `json:"uuid" gorm:"column:uuid"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	OutPort    int    `json:"out_port"`
	//CaPem      string `json:"ca_pem"`
	//CertPem    string `json:"cert_pem"`
	//KeyPem     string `json:"key_pem"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

func (Server) TableName() string {
	return "zta_server"
}

type ServerAttr struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	OutPort int    `json:"out_port"`
}

func (c ServerAttr) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ServerAttr) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
