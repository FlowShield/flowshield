package verify

import (
	"context"
	"sync"
	"time"

	"github.com/flowshield/flowshield/verifier/pkg/recover"
)

var VerObj = NewVerifier(&Options{Often: 30})

type Verifier struct {
	opt      *Options
	ctx      context.Context
	Provider []*Provider
	ticker   *time.Ticker
	Record   map[string]*Record
	Health   map[string]*Record
	sync.RWMutex
}

func NewVerifier(opt *Options) (verifier *Verifier) {
	verifier = &Verifier{
		opt:    opt,
		ticker: time.NewTicker(time.Second * opt.Often),
		//ticker: time.NewTicker(time.Second),
		Record: make(map[string]*Record),
		Health: make(map[string]*Record),
	}
	return
}

func (v *Verifier) Run(ctx context.Context) error {
	recover.Recovery(ctx, func() {
		for {
			select {
			case <-v.ticker.C:
				go v.SyncProviderAndOrder()
				go v.Statistics()
			}
		}
	})
	return nil
}

func (v *Verifier) SyncProviderAndOrder() {
	ordersMysql, err := orders()
	if err != nil {
		return
	}
	providers, err := providers(ordersMysql)
	if err != nil {
		return
	}
	v.Lock()
	defer v.Unlock()
	v.Provider = providers
	if v.Provider == nil {
		return
	}
	for _, value := range v.Provider {
		for _, va := range value.Order {
			va.CheckHealthy(value.IP)
			if v.Record[value.PeerId] == nil {
				v.Record[value.PeerId] = &Record{}
			}
			if v.Health[va.OrderID] == nil {
				v.Health[va.OrderID] = &Record{}
			}
			if !va.Healthy.Health {
				record := v.Record[value.PeerId]
				record.AddFail(1)
				health := v.Health[va.OrderID]
				health.AddFail(1)
			} else {
				record := v.Record[value.PeerId]
				record.AddSuccess(1)
				health := v.Health[va.OrderID]
				health.AddFail(1)
			}
		}
	}
}

func (v *Verifier) Statistics() {
	v.RLock()
	defer v.RUnlock()
	if v.Record == nil {
		return
	}
	//for key, value := range v.Record {
	// TODO
	//log.Printf("provider %s can not be connected for %d times!\n", key, value.Fail)
	//}
}

func (v *Verifier) ProviderHealth(orderID []string) (record map[string]*Record) {
	if len(orderID) == 0 {
		return nil
	}
	record = make(map[string]*Record)
	v.RLock()
	defer v.RUnlock()
	for _, value := range orderID {
		record[value] = v.Health[value]
	}
	return
}
