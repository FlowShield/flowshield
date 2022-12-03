package model

import "github.com/flowshield/flowshield/provider/internal/schema"

type Providers []*Provider

type Provider struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiredAt int64  `gorm:"autoExpiredTime" json:"expired_at"`
	Uuid      string `gorm:"uuid" json:"uuid"`
	Wallet    string `gorm:"wallet" json:"wallet"`
	ServerCid string `gorm:"server_cid" json:"server_cid"`
	Port      int    `gorm:"port" json:"port"`
}

func (Provider) TableName() string {
	return "pr_provider"
}

func (a *Provider) ToNodeOrder() *schema.NodeOrder {
	result := &schema.NodeOrder{
		Uuid:      a.Uuid,
		Wallet:    a.Wallet,
		ServerCid: a.ServerCid,
		Port:      a.Port,
	}
	return result
}

func OrderToProvider(order *schema.NodeOrder, pc *schema.ProviderConfig) *Provider {
	result := &Provider{
		Uuid:      order.Uuid,
		Wallet:    order.Wallet,
		ServerCid: order.ServerCid,
		Port:      order.Port,
		ExpiredAt: pc.CertConfig.NotAfter.Unix(),
	}
	return result
}
