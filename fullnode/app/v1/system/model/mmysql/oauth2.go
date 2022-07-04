package mmysql

import (
	"database/sql/driver"
	"encoding/json"
)

type Oauth2 struct {
	ID           uint   `gorm:"primarykey"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
	Company      string `json:"company" binding:"required,oneof= github facebook google"`
	ClientId     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectUrl  string `json:"redirect_url" binding:"required"`
	Scopes       Scopes `json:"scopes" binding:"required"`
	AuthUrl      string `json:"auth_url" binding:"required"`
	TokenUrl     string `json:"token_url" binding:"required"`
}

func (Oauth2) TableName() string {
	return "zta_oauth2"
}

type Scopes []string

func (c Scopes) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Scopes) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
