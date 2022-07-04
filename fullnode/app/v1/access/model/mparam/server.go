package mparam

import "github.com/cloudslit/cloudslit/fullnode/app/base/mdb"

type ServerList struct {
	mdb.Paginate
	Name       string `json:"name" form:"name"`
	ResourceID int    `json:"resource_id" form:"resource_id"`
}

type AddServer struct {
	Name       string `json:"name" form:"name" binding:"required"`
	ResourceID string `json:"resource_id" form:"resource_id"`
	Host       string `json:"host" form:"host" binding:"required"`         // api.github.com
	Port       int    `json:"port" form:"port" binding:"required"`         // 443
	OutPort    int    `json:"out_port" form:"out_port" binding:"required"` // 443
}

type EditServer struct {
	ID         uint64 `json:"id" form:"id" binding:"required"`
	Name       string `json:"name" form:"name" binding:"required"`
	ResourceID string `json:"resource_id" form:"resource_id"`
	Host       string `json:"host" form:"host" binding:"required"`         // api.github.com
	Port       int    `json:"port" form:"port" binding:"required"`         // 443
	OutPort    int    `json:"out_port" form:"out_port" binding:"required"` // 443
}
