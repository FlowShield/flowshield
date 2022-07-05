package mapi

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
)

type ServerList struct {
	List     []mmysql.Server    `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
