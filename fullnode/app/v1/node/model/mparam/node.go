package mparam

import (
	"github.com/cloudslit/cloudslit/fullnode/app/base/mdb"
)

type ListNode struct {
	mdb.Paginate
	PeerId   string   `json:"peer_id" form:"peer_id"`
	IP       string   `json:"ip" form:"ip"`
	Loc      []string `json:"loc" form:"loc"`
	Colo     string   `json:"colo" form:"colo"`
	GasPrice int      `json:"gas_price" form:"gas_price"`
	Type     string   `json:"type" form:"type"`
}
