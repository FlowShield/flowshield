package mapi

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
)

type ResourceList struct {
	List     []mmysql.Resource  `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
