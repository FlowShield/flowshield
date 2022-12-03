package mapi

import (
	"github.com/flowshield/flowshield/fullnode/app/base/mapi"
	"github.com/flowshield/flowshield/fullnode/app/v1/access/model/mmysql"
)

type ResourceList struct {
	List     []mmysql.Resource  `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
