package mparam

import "github.com/flowshield/flowshield/fullnode/app/base/mdb"

type ResourceList struct {
	mdb.Paginate
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}

type AddResource struct {
	Name string `json:"name" form:"name" binding:"required"`
	Type string `json:"type" form:"name" binding:"oneof=cidr dns"`
	Host string `json:"host" form:"host" binding:"required"` // api.github.com,192.168.1.1/16
	Port string `json:"port" form:"port" binding:"required"` // 80-443;3306;6379
}

type EditResource struct {
	ID   uint64 `json:"id" form:"id" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
	Type string `json:"type" form:"name" binding:"oneof=cidr dns"`
	Host string `json:"host" form:"host" binding:"required"` // api.github.com,192.168.1.1/16
	Port string `json:"port" form:"port" binding:"required"` // 80-443;3306;6379
}
