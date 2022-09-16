package verify

import (
	"time"
)

type Provider struct {
	PeerId string   `json:"peer_id"`
	Addr   string   `json:"addr"`
	Port   int      `json:"port"`
	IP     string   `json:"ip"`
	Loc    string   `json:"loc"`
	Colo   string   `json:"colo"`
	Price  uint     `json:"price"`
	Order  []*Order `json:"order"`
}

func providers(orders []*OrderMysql) (providers []*Provider, err error) {
	if orders == nil || len(orders) == 0 {
		return
	}
	var providerMap = make(map[string][]*OrderMysql)
	for _, value := range orders {
		providerMap[value.PeerID] = append(providerMap[value.PeerID], value)
	}
	for key, value := range providerMap {
		orderSli := make([]*Order, 0)
		for _, v := range value {
			orderObj := &Order{
				PeerID:    v.PeerID,
				Port:      v.Port,
				OrderID:   v.UUID,
				StartTime: time.Unix(v.UpdatedAt, 0),
				EndTime:   time.Unix(v.UpdatedAt, 0).Add(time.Duration(v.Duration) * 60 * 60),
				Healthy:   nil,
			}
			orderSli = append(orderSli, orderObj)
		}
		provider := &Provider{
			PeerId: key,
			Addr:   value[0].NodeIP,
			//Port:   value[0].Port,
			IP: value[0].NodeIP,
			//Loc:    "",
			//Colo:   "",
			//Price:  0,
			Order: orderSli,
		}
		providers = append(providers, provider)
	}
	return
}
