package mapi

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mapi"
	"github.com/cloudslit/cloudslit/fullnode/app/v1/access/model/mmysql"
)

type ClientList struct {
	List     []mmysql.Client    `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
