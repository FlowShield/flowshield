package mapi

import (
	"github.com/flowshield/flowshield/fullnode/app/base/mapi"
	"github.com/flowshield/flowshield/fullnode/app/v1/node/model/mmysql"
)

type NodeList struct {
	List     []mmysql.Node      `json:"list"`
	Paginate mapi.AdminPaginate `json:"paginate"`
}
