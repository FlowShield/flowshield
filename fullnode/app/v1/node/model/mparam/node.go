package mparam

import (
	"github.com/flowshield/flowshield/fullnode/app/base/mdb"
)

type ListNode struct {
	mdb.Paginate
	PeerId string   `json:"peer_id" form:"peer_id"`
	IP     string   `json:"ip" form:"ip"`
	Loc    []string `json:"loc" form:"loc"`
	Colo   string   `json:"colo" form:"colo"`
	Price  int      `json:"price" form:"price"`
	Type   string   `json:"type" form:"type"`
}
