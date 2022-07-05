package mapi

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/node/model/mmysql"
)

type NodeList struct {
	List     []mmysql.Node      `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
